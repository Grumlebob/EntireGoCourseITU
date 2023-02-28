package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	replicationRole               string
	idOfPrimaryReplicationManager int32
	clients                       map[int32]node.NodeClient
	requestsHandled               map[int32]bool
	distributedDictionary         map[string]string
	ctx                           context.Context
	lastHeartbeat                 *timestamppb.Timestamp
}

const (
	// Node role for primary-back replication
	PRIMARY = "PRIMARY"
	BACKUP  = "BACKUP"
)

var (
	lock = make(chan bool, 1)
)

var uniqueIdentifier = int32(0)

func main() {
	//When lock has content, it is unlocked. When it is empty, it is locked. Initialize with a value, so it is unlocked.
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
		clients:                     make(map[int32]node.NodeClient),
		requestsHandled:             make(map[int32]bool),
		ctx:                         ctx,
		replicationRole:             BACKUP,
		distributedDictionary:       make(map[string]string),
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

	//Peer to peer architecture - 2 Peers connected on port 5000, 5001
	for i := 0; i < 2; i++ {
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

	log.Printf("I am peer with Id: %v", p.id)
	//Assign primary / backup role
	p.LeaderElection()

	//Leader sends heartbeats to backups
	if p.replicationRole == PRIMARY {
		go p.leaderSendHeartbeatToBackups()
	}

	//Backups check their heartbeats and does leader election on heartbeat timeout.
	if p.replicationRole == BACKUP {
		//Initial wait, so leader has time to send first heartbeat.
		time.Sleep(1000 * time.Millisecond)
		go p.backupCheckHeartbeat()
	}

	//Instructions for user
	log.Println("Welcome to the distributed dictionary")
	log.Println("Write 'add <key> <value>' to put a key-value pair in the dictionary")
	log.Println("Write 'read <key>' to get the value of a key")
	log.Println("Write 'printrequests' to print map of requests already handled")
	log.Println("Write 'printdictionary' to print the distributed dictionary")

	//Automatic system demo
	if p.replicationRole == PRIMARY {
		log.Println("--- Adding demo data ---")
		for i := 0; i < 5; i++ {
			word := fmt.Sprintf("word%v", i)
			definition := fmt.Sprintf("definition%v", i)
			//Add word to dictionary and read it back
			p.Add(p.ctx, &node.DictionaryAdd{Word: word, Definition: definition, UniqueIdentifierForRequest: p.getUniqueIdentifier()})
			resultOne, _ := p.Read(p.ctx, &node.ReadWord{Word: word})
			//Correct Add - Read consistency correctness check
			log.Printf("Read word: %v. Had: %v. Expected: %v", word, resultOne.Definition, definition)
			if resultOne.Definition != definition {
				log.Printf("Should not happen! - Read was not the last add!")
			}
			//Changing definition of word, to check that read returns the latest definition
			p.Add(p.ctx, &node.DictionaryAdd{Word: word, Definition: definition + "New", UniqueIdentifierForRequest: p.getUniqueIdentifier()})
			resultTwo, _ := p.Read(p.ctx, &node.ReadWord{Word: word})
			if resultTwo.Definition != definition {
				log.Printf("Read word: %v. Old: %v. New: %v. Expected: %v", word, resultOne.Definition, resultTwo.Definition, definition+"New")
				if resultTwo.Definition == definition+"New" {
					log.Printf("Success - Latest value was read")
				}
			} else {
				log.Printf("Should not happen! - Read was not the last add!")
			}
		}
	}
	//print demo data
	log.Printf("--- Printing maps on all nodes. Maps should be equal ---")
	log.Printf("Size of distributed dictionary: %v. Should be: %v", len(p.distributedDictionary), 5)
	log.Printf("Size of unique requests handled map: %v. Should be: %v", len(p.requestsHandled), 10)
	p.printDistributedDictionary()
	p.printRequestsHandledMap()
	log.Printf("--- Demo done ---")

	//Read user input - Keeps main thread alive
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		textSplitOnWhitespace := strings.Fields(text)
		//ADD command
		if len(textSplitOnWhitespace) == 3 && textSplitOnWhitespace[0] == "add" {
			//Phase 1 of passive replication, Request: The front end issues the request, containing a unique identifier, to the primary.
			uniqueId := p.getUniqueIdentifier()
			//Re-send. Happens in case a primary crashes, and a new one is promoted, that message might have been lost. And to ensure liveliness, we want to resend it to the newly elected primary. Thus At-least-once semantics.
			for {
				//Always forward to primary
				if p.replicationRole == PRIMARY {
					response, _ := p.Add(p.ctx, &node.DictionaryAdd{Word: textSplitOnWhitespace[1], Definition: textSplitOnWhitespace[2], UniqueIdentifierForRequest: uniqueId})
					log.Printf("Was added: %v. Word: %v. Definition: %v", response.Result, textSplitOnWhitespace[1], textSplitOnWhitespace[2])
				} else if p.replicationRole == BACKUP {
					//Check primary alive
					p.sendHeartbeatPingToPeerId(p.idOfPrimaryReplicationManager)
					_, foundPrimary := p.clients[p.idOfPrimaryReplicationManager]
					if foundPrimary {
						response, _ := p.clients[p.idOfPrimaryReplicationManager].Add(p.ctx, &node.DictionaryAdd{Word: textSplitOnWhitespace[1], Definition: textSplitOnWhitespace[2], UniqueIdentifierForRequest: uniqueId})
						log.Printf("Was added: %v. Word: %v. Definition: %v", response.Result, textSplitOnWhitespace[1], textSplitOnWhitespace[2])
					}
				}
				//If response is stored break, else re-send
				_, foundUniqueRequest := p.requestsHandled[uniqueId]
				if foundUniqueRequest {
					break
				} else {
					log.Println("Request not found, re-sending after 1 second")
					time.Sleep(time.Second * 1)
				}

			}
		}
		//READ command
		//No need to re-send, since it is a read-only operation and thus at-most-once semantics
		if len(textSplitOnWhitespace) == 2 && textSplitOnWhitespace[0] == "read" {
			if p.replicationRole == PRIMARY {
				outcome, _ := p.Read(p.ctx, &node.ReadWord{Word: textSplitOnWhitespace[1]})
				log.Printf("Read word: %v. Has definition: %v", textSplitOnWhitespace[1], outcome.Definition)
			} else if p.replicationRole == BACKUP {
				p.sendHeartbeatPingToPeerId(p.idOfPrimaryReplicationManager)
				_, foundPrimary := p.clients[p.idOfPrimaryReplicationManager]
				if foundPrimary {
					outcome, _ := p.clients[p.idOfPrimaryReplicationManager].Read(p.ctx, &node.ReadWord{Word: textSplitOnWhitespace[1]})
					log.Printf("Read word: %v. Has definition: %v", textSplitOnWhitespace[1], outcome.Definition)

				}
			}
		}
		//force CRASH command
		if strings.Contains(text, "crash") && len(textSplitOnWhitespace) == 1 {
			log.Printf("Crashing node id %v ", p.id)
			os.Exit(1)
		}
		//Print requests already handled
		if strings.Contains(text, "printrequests") && len(textSplitOnWhitespace) == 1 {
			p.printRequestsHandledMap()
		}
		//Print Distributed Dictionary
		if strings.Contains(text, "printdictionary") && len(textSplitOnWhitespace) == 1 {
			p.printDistributedDictionary()
		}
	}
}

func (p *peer) LeaderElection() {
	//Bully principles - Highest Id wins, and announces to rest. All other are backups.
	//In this implementation no election/alive message is sent during election, instead if the election is not completed, the heartbeat timeout will trigger (which is kinda the same as if a timeout on answer/coordinate ran out), starting a new election and killing off any dead peers. This only works because all peers knows about all peers in whole system.
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

// Bully winner - Highest ID declared himself as leader to all other.
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
					return
				}
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

func (p *peer) getAgreementFromAllPeersAndReplicateLeaderData(word string, definition string, ack bool, identifier int32) (agreementReached bool) {
	//N-1 Agreements needed from backups
	p.agreementsNeededFromBackups = int32(len(p.clients))
	for id, client := range p.clients {
		_, err := client.HandleAgreementAndReplicationFromLeader(p.ctx, &node.Replicate{AddWord: word, AddDefinition: definition, ResponseForRequest: ack, UniqueIdentifierForRequest: identifier})
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
	p.requestsHandled[replicate.UniqueIdentifierForRequest] = replicate.ResponseForRequest
	p.distributedDictionary[replicate.AddWord] = replicate.AddDefinition
	return &emptypb.Empty{}, nil
}

func (p *peer) getUniqueIdentifier() (uniqueId int32) {
	uniqueIdentifier++
	asString := fmt.Sprintf("%v%v", uniqueIdentifier, p.id)
	//If we just convert to int32, we will get 1+5000 = 5001, which is not unique.
	//Instead we concatenate as a string, to get 15000, making it unique until we reach overflow
	realId, _ := strconv.ParseInt(asString, 10, 32)
	return int32(realId)
}

func (p *peer) Add(ctx context.Context, in *node.DictionaryAdd) (*node.AddResult, error) {
	//Uses lock for exclusion, so there is no race-condition between clients reading/adding concurrently.
	<-lock

	addResult := &node.AddResult{}

	//If already processed, send same response.
	previousResponse, found := p.requestsHandled[in.UniqueIdentifierForRequest]
	if found {
		addResult.Result = previousResponse
		return addResult, nil
	}

	//Passive replication phase 3: Primary executes request, and stores response
	p.distributedDictionary[in.Word] = in.Definition
	addResult.Result = true
	p.requestsHandled[in.UniqueIdentifierForRequest] = addResult.Result

	//Passive replication phase 4: Primary sends update to backups, and they send back an acknowledgement.
	agreement := p.getAgreementFromAllPeersAndReplicateLeaderData(in.Word, in.Definition, addResult.Result, in.UniqueIdentifierForRequest)
	if !agreement {
		addResult.Result = false
		p.requestsHandled[in.UniqueIdentifierForRequest] = addResult.Result
	}

	lock <- true
	//Passive replication phase 5: Primary responds to client
	return addResult, nil
}

func (p *peer) Read(ctx context.Context, in *node.ReadWord) (*node.ReadResult, error) {
	//Locking to prevent race conditions, such as another client adding a word while we are reading it.
	//No need for updating replicas, as operation is read-only.
	<-lock
	result := &node.ReadResult{Definition: p.distributedDictionary[in.Word]}
	lock <- true
	return result, nil
}

// A helper Method to read key values of a distributed dictionary
func (p *peer) printDistributedDictionary() {
	log.Printf("--- Printing Distributed Dictionary ---\n")
	for key, value := range p.distributedDictionary {
		log.Printf("Word: %v, Definition: %v\n", key, value)
	}
}

// A helper Method to read key values of the handled requests
func (p *peer) printRequestsHandledMap() {
	log.Printf("--- Printing requests already handled ---\n")
	for key, value := range p.requestsHandled {
		log.Printf("Unique Request Id: %v, Response: %v\n", key, value)
	}
}
