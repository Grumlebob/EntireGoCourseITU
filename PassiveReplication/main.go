package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	node "jacob.com/go-grpc/grpc"
)

type peer struct {
	node.UnimplementedNodeServer
	id                            int32
	agreementsNeededFromBackups   int32
	auctionState                  string
	replicationRole               string
	idOfPrimaryReplicationManager int32
	highestBidOnCurrentAuction    int32
	winnerId                      int32
	currentItem                   string
	clients                       map[int32]node.NodeClient
	requestsHandled               map[int32]string
	ctx                           context.Context
	lastHeartbeat                 *timestamppb.Timestamp
}

//Primary-backup replication implementation.
//• 1. Request: The front end issues the request, containing a unique identifier, to the primary.
//• 2. Coordination: The primary takes each request atomically, in the order in which
//it receives it. It checks the unique identifier, in case it has already executed the
//request, and if so it simply resends the response.
//• 3. Execution: The primary executes the request and stores the response.
//• 4. Agreement: If the request is an update, then the primary sends the updated
//state, the response and the unique identifier to all the backups. The backups send an acknowledgement.
//• 5. Response: The primary responds to the front end, which hands the response back to the client.

const (
	// Node role for primary-back replication
	PRIMARY = "PRIMARY"
	BACKUP  = "BACKUP"
)

const (
	// Auction state
	OPEN   = "Open"
	CLOSED = "Closed"
)

var (
	lock = make(chan bool, 1)
)

var uniqueIdentifier = int32(0)

func main() {
	lock <- true
	//Get port from command line
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	//Create logfile with ownport as name. Print to both stdout and logfile
	logfile := fmt.Sprintf("%d%s", ownPort, "Log.txt")
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetFlags(log.Ltime)
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := &peer{
		id:                          ownPort,
		agreementsNeededFromBackups: 0,
		winnerId:                    0,
		clients:                     make(map[int32]node.NodeClient),
		requestsHandled:             make(map[int32]string),
		ctx:                         ctx,
		auctionState:                CLOSED,
		replicationRole:             BACKUP,
		currentItem:                 "",
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	node.RegisterNodeServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()
	//3 Peers connected on port 5000, 5001, 5002
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)
		if port == ownPort {
			continue
		}
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := node.NewNodeClient(conn)
		p.clients[port] = c
	}

	//Assign primary / backup role
	log.Printf("I am peer %v", p.id)
	p.LeaderElection()

	//Leader sends heartbeats to backups
	if p.replicationRole == PRIMARY {
		p.leaderSendHeartbeatToBackups()
	}
	//Backups check their heartbeats and does leader election, if haven't received a heartbeat
	if p.replicationRole == BACKUP {
		//Initial wait, so leader has time to send first heartbeat.
		time.Sleep(1000 * time.Millisecond)
		go p.backupCheckHeartbeat()
	}

	//Open and close auctions after set timers
	if p.replicationRole == PRIMARY {
		p.hostAuction()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		//RESULT COMMAND
		if strings.Contains(text, "result") {
			if p.replicationRole == PRIMARY {
				outcome, _ := p.Result(p.ctx, &emptypb.Empty{})
				log.Printf("%s \nThe highest bid: %v", outcome.AuctionStatus, outcome.HighestBid)
			} else if p.replicationRole == BACKUP {
				p.sendHeartbeatPingToPeerId(p.idOfPrimaryReplicationManager)
				_, foundPrimary := p.clients[p.idOfPrimaryReplicationManager]
				if foundPrimary {
					outcome, _ := p.clients[p.idOfPrimaryReplicationManager].Result(p.ctx, &emptypb.Empty{})
					log.Printf("%s \nThe highest bid: %v", outcome.AuctionStatus, outcome.HighestBid)
				}
			}
		}
		//BID COMMAND
		FindAllNumbers := regexp.MustCompile(`\d+`).FindAllString(text, 1)
		if len(FindAllNumbers) > 0 {
			numeric, _ := strconv.ParseInt(FindAllNumbers[0], 10, 32)
			if strings.Contains(text, "bid") && numeric > 0 {
				uniqueId := p.getUniqueIdentifier()
				//Re-send. Happens in case a primary crashes, and a new one is promoted, that message might have been lost.
				for {
					if p.replicationRole == PRIMARY {
						ack, _ := p.Bid(p.ctx, &node.Bid{ClientId: p.id, UniqueBidId: uniqueId, Amount: int32(numeric)})
						log.Println(ack.Ack)
					} else if p.replicationRole == BACKUP {
						p.sendHeartbeatPingToPeerId(p.idOfPrimaryReplicationManager)
						_, foundPrimary := p.clients[p.idOfPrimaryReplicationManager]
						if foundPrimary {
							ack, _ := p.clients[p.idOfPrimaryReplicationManager].Bid(p.ctx, &node.Bid{ClientId: p.id, UniqueBidId: uniqueId, Amount: int32(numeric)})
							log.Println(ack.Ack)
						}
					}
					//If response is stored break, else re-send
					_, foundUniqueRequest := p.requestsHandled[uniqueId]
					if foundUniqueRequest {
						break
					}
				}
			}
		}
		//FORCE CRASH COMMAND
		if strings.Contains(text, "crash") {
			log.Printf("Crashing node id %v ", p.id)
			os.Exit(1)
		}
	}
}

func (p *peer) LeaderElection() {
	//Bully principles - Highest Id wins, and announces to rest. All other are backups.
	//In this implementation no election/alive message is sent during election, instead if the election is not completed, the heartbeat timeout will trigger (which is kinda the same as if a timeout on answer/coordinate ran out), starting a new election and killing off any dead peers.
	//When heartbeat fails, only the highest peer will change role to leader, and rest will simply be re-assigned as backups.
	shouldBeRole := BACKUP
	highestClientSeen := int32(0)
	for i := range p.clients {
		if i < p.id && p.id > highestClientSeen {
			shouldBeRole = PRIMARY
			highestClientSeen = p.id
		} else if i > p.id {
			highestClientSeen = i
			shouldBeRole = BACKUP
		}
	}
	p.replicationRole = shouldBeRole
	p.idOfPrimaryReplicationManager = highestClientSeen
	if len(p.clients) == 0 {
		p.replicationRole = PRIMARY
		p.idOfPrimaryReplicationManager = p.id
	}
	if p.replicationRole == PRIMARY {
		for i := range p.clients {
			if i != p.id {
				_, err := p.clients[i].Coordinator(p.ctx, &node.CoordinatorPort{LeaderPort: p.id})
				if err != nil {
					log.Printf("Error sending coordinator message to node %v", i)
					delete(p.clients, i)
				}
			}
		}
	}
	log.Printf("Replication role: %v \n", p.replicationRole)
}

// Bully - Highest ID declared himself as leader to all other.
func (p *peer) Coordinator(ctx context.Context, leaderPort *node.CoordinatorPort) (*emptypb.Empty, error) {
	p.idOfPrimaryReplicationManager = leaderPort.LeaderPort
	p.replicationRole = BACKUP
	log.Printf("Leader is now: %v \n", p.idOfPrimaryReplicationManager)
	return &emptypb.Empty{}, nil
}

// Helper method to send a heartbeat ping to a peer, while checking if peer is alive.
func (p *peer) sendHeartbeatPingToPeerId(peerId int32) {
	//Check map of clients, if it exists.
	_, found := p.clients[peerId]
	if !found {
		log.Printf("Client node %v is dead", peerId)
		//Remove dead node from list of clients
		delete(p.clients, peerId)
		return
	}
	//Ping the node with current timestamp.
	t := time.Now()
	time := timestamppb.New(t)
	heartbeatTimestamp := node.HeartbeatTimestamp{
		Timestamp: time,
	}
	_, err := p.clients[peerId].Heartbeat(p.ctx, &heartbeatTimestamp)
	if err != nil {
		//Remove dead node from list of clients
		log.Printf("Client node %v is dead", peerId)
		delete(p.clients, peerId)
	}
}

// Heartbeat
func (p *peer) Heartbeat(ctx context.Context, heartbeat *node.HeartbeatTimestamp) (*emptypb.Empty, error) {
	p.lastHeartbeat = heartbeat.Timestamp
	return &emptypb.Empty{}, nil
}

// Backups continuously checks heartbeat from leader.
func (p *peer) backupCheckHeartbeat() {
	go func() {
		for {
			//Every 500 ms, they check their heartbeat time compared to current time.
			//If over 3 seconds, leader must have crashed and a new election is held.
			time.Sleep(500 * time.Millisecond)
			lastBeat := p.lastHeartbeat
			t := time.Now()
			currentTime := timestamppb.New(t)
			if currentTime.AsTime().Sub(lastBeat.AsTime()).Seconds() > 5 {
				fmt.Println("Leader is dead. Starting election")
				delete(p.clients, p.idOfPrimaryReplicationManager)
				p.LeaderElection()
				//If p is promoted, it should instead send heartbeats and stop listening for heartbeats themselves.
				if p.replicationRole == PRIMARY {
					go p.leaderSendHeartbeatToBackups()
					go p.hostAuction()
					return
				}
			} else {
				//Heartbeat is still alive.
				//fmt.Println("hb: ", currentTime.AsTime().Sub(lastBeat.AsTime()).Seconds())
			}
		}
	}()
}

// Leader continuously sends heartbeat to backups.
func (p *peer) leaderSendHeartbeatToBackups() {
	go func() {
		for {
			for id := range p.clients {
				p.sendHeartbeatPingToPeerId(id)
			}
			time.Sleep(300 * time.Millisecond)
		}
	}()
}

func (p *peer) getAgreementFromAllPeersAndReplicateLeaderData(ack string, identifier int32) (agreementReached bool) {
	//N-1 Agreements needed
	p.agreementsNeededFromBackups = int32(len(p.clients))
	for id, client := range p.clients {
		_, err := client.HandleAgreementAndReplicationFromLeader(p.ctx, &node.Replicate{AuctionStatus: p.auctionState, HighestBidOnCurrentAuction: p.highestBidOnCurrentAuction, ResponseForRequest: ack, UniqueIdentifierForRequest: identifier, CurrentItem: p.currentItem, WinnerId: p.winnerId})
		if err != nil {
			log.Printf("Client node %v is dead", id)
			//Remove dead node from list of clients
			delete(p.clients, id)
		}
		p.agreementsNeededFromBackups--
	}
	return p.agreementsNeededFromBackups == 0
}

func (p *peer) HandleAgreementAndReplicationFromLeader(ctx context.Context, replicate *node.Replicate) (*emptypb.Empty, error) {
	p.auctionState = replicate.AuctionStatus
	p.highestBidOnCurrentAuction = replicate.HighestBidOnCurrentAuction
	p.currentItem = replicate.CurrentItem
	p.winnerId = replicate.WinnerId
	p.requestsHandled[replicate.UniqueIdentifierForRequest] = replicate.ResponseForRequest
	return &emptypb.Empty{}, nil
}

// Helper method to broadcast to all peers
func (p *peer) SendMessageToAllPeers(message string) {
	for id, client := range p.clients {
		_, err := client.BroadcastMessage(p.ctx, &node.MessageString{Text: message})
		if err != nil {
			log.Printf("Client node %v is dead", id)
			//Remove dead node from list of clients
			delete(p.clients, id)
		}
	}
}

// Send a string to a peer
func (p *peer) BroadcastMessage(ctx context.Context, message *node.MessageString) (*emptypb.Empty, error) {
	log.Println(message.Text)
	return &emptypb.Empty{}, nil
}

func (p *peer) getUniqueIdentifier() (uniqueId int32) {
	uniqueIdentifier++
	asString := fmt.Sprintf("%v%v", uniqueIdentifier, p.id)
	//We don't want to 1+5000, we want 15000 to make them unique...well, unique enough, until overflow.
	realId, _ := strconv.ParseInt(asString, 10, 32)
	return int32(realId)
}

func (p *peer) Bid(ctx context.Context, bid *node.Bid) (*node.Acknowledgement, error) {
	<-lock
	agreement := false
	acknowledgement := &node.Acknowledgement{}

	//If already processed, send same Ack.
	ack, found := p.requestsHandled[bid.UniqueBidId]
	if found {
		acknowledgement.Ack = ack
	} else if p.auctionState == CLOSED && !found {
		acknowledgement.Ack = "Fail, auction is closed"
	} else if (p.highestBidOnCurrentAuction < bid.Amount) && (!found) && (p.auctionState == OPEN) {
		p.highestBidOnCurrentAuction = bid.Amount
		p.winnerId = bid.ClientId
		acknowledgement.Ack = "OK"
	} else if p.highestBidOnCurrentAuction >= bid.Amount && !found && p.auctionState == OPEN {
		acknowledgement.Ack = "Fail, your bid was too low"
	}

	p.requestsHandled[bid.UniqueBidId] = acknowledgement.Ack

	agreement = p.getAgreementFromAllPeersAndReplicateLeaderData(acknowledgement.Ack, bid.UniqueBidId)
	if !agreement {
		//Can theoretically not enter this if block, due to assumptions stated in the assignment
		acknowledgement.Ack = "Fail, couldn't reach agreement in phase 4"
		p.requestsHandled[bid.UniqueBidId] = acknowledgement.Ack
	}
	lock <- true
	return acknowledgement, nil
}

func (p *peer) Result(ctx context.Context, empty *emptypb.Empty) (*node.Outcome, error) {
	if p.auctionState == OPEN {
		return &node.Outcome{AuctionStatus: "Auction is OPEN", HighestBid: p.highestBidOnCurrentAuction}, nil
	} else {
		return &node.Outcome{AuctionStatus: "Auction is CLOSED", HighestBid: p.highestBidOnCurrentAuction}, nil
	}
}

func (p *peer) hostAuction() {
	go func() {
		items := []string{"Laptop", "Phone", "Tablet", "TV", "Headphones", "Watch"}
		for _, item := range items {
			p.openAuction(item)
			time.Sleep(45 * time.Second)
			p.closeAuction()
			time.Sleep(10 * time.Second)
		}
	}()
}

func (p *peer) openAuction(item string) {
	if p.auctionState == OPEN {
		return
	}
	p.auctionState = OPEN
	p.currentItem = item
	p.agreementsNeededFromBackups = int32(len(p.clients))
	p.highestBidOnCurrentAuction = 0
	p.winnerId = 0
	announceAuction := "Auction for " + item + " is open for 45 seconds! \nEnter Bid <amount> to bid on the item. \nEnter Result to see highest bid"
	log.Println(announceAuction)
	p.SendMessageToAllPeers(announceAuction)
}

func (p *peer) closeAuction() {
	p.auctionState = CLOSED
	p.agreementsNeededFromBackups = int32(len(p.clients))
	announceWinner := fmt.Sprintf("Auction for %s is Over! \nHighest bid was %v\nAuction was won by %v\nNext auction starts in 10 seconds", p.currentItem, p.highestBidOnCurrentAuction, p.winnerId)
	log.Println(announceWinner)
	p.SendMessageToAllPeers(announceWinner)
}
