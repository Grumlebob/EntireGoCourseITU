package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	guid "github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	node "jacob.com/go-grpc/grpc"
)

var userId = "0"

var (
	servers map[int32]node.ReplicaServiceClient
	ctx     context.Context
)

var (
	lock = make(chan bool, 1)
)

func main() {
	//Print date and time in log.
	//Create logfile with ownport as name. Print to both stdout and logfile
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetFlags(log.Ltime)
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	servers = make(map[int32]node.ReplicaServiceClient)
	context, cancel := context.WithCancel(context.Background())
	ctx = context
	defer cancel()
	lock <- true

	//Create a connection to each replica server, 4000, 4001, 4002
	for i := 0; i < 3; i++ {
		port := int32(4000) + int32(i)
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		server := node.NewReplicaServiceClient(conn)
		servers[port] = server
	}

	//Assigns a unique ID to a new client
	getClientId()

	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		go sendIncrementToAllServers(1)
	}

	fmt.Println("Enter 'leave()' to leave the chatroom. \nEnter your message here:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		FindAllNumbers := regexp.MustCompile(`\d+`).FindAllString(text, 1)
		if len(FindAllNumbers) != 0 {
			numeric, _ := strconv.ParseInt(FindAllNumbers[0], 10, 32)
			go sendIncrementToAllServers(numeric)
		}
	}
	// block main
	block := make(chan bool)
	<-block
}

func sendIncrementToAllServers(increment int64) (amount int64, err error) {
	<-lock
	//Map for checking consistency of servers
	var mapOfResponse = make(map[int64]int32)

	for serverId, server := range servers {
		//Send something to server
		response, err := server.Increment(ctx, &node.RequestIncrement{IncrementAmount: increment})
		if err != nil {
			//Delete server if it has crashed
			log.Printf("Server %v down - Deleting from map", serverId)
			delete(servers, serverId)
		} else {
			//If server alive, add response to map. Increment frequency counter
			if mapOfResponse[response.CurrentAmount] == 0 {
				mapOfResponse[response.CurrentAmount] = 1
			} else {
				mapOfResponse[response.CurrentAmount]++
			}
		}
	}
	//If all servers sent same answer, then we are consistent and length of map is 1
	length := len(mapOfResponse)
	//log.Print("Length of map: ", length)
	if length > 1 {
		//If length over 2 we are not getting consistent answers and wants to recheck.
		//Resend increment with 0
		for {
			log.Printf("Received %v different responses. Rechecking for consistency \n", length)
			lock <- true
			amount, err := sendIncrementToAllServers(0)
			if err == nil {
				log.Println("We are consistent with amount: ", amount)
				return amount, nil
			}
		}
	} else {
		//Print the 1 response
		for key, _ := range mapOfResponse {
			log.Printf("Consistent answer: %v \n", key)
			lock <- true
			return key, nil
		}
		//If map contains no values (should not happen), return -1
		return -1, nil
	}
}

func getClientId() {
	if userId != "0" {
		return
	}
	//Uses library "github.com/rs/xid" to generate a unique ID for each client
	userId = guid.New().String()
	log.Println("Hello! - You are the client with ID: ", userId)
}

func sendSomethingToAllServers(something string) (result string, err error) {

	//Checks consistency of servers
	<-lock
	var mapOfResponse = make(map[string]int32)

	for serverId, server := range servers {
		//Send something to server
		response, err := server.SendSomeMessage(ctx, &node.SomeMessage{Message: something})
		if err != nil {
			//delete server here
			log.Printf("Server %v down - Deleting from map", serverId)
			delete(servers, serverId)
		} else {
			//Add to map, and increment frequency counter
			if mapOfResponse[response.Message] == 0 {
				mapOfResponse[response.Message] = 1
			} else {
				mapOfResponse[response.Message]++
			}
		}
	}

	length := len(mapOfResponse)
	//log.Print("Length of map: ", length)
	if length > 1 {
		//If length over 2, then we have a problem, and we need to check for consistency
		//Resend increment with 0
		for {
			log.Printf("Received %v different responses. Rechecking for consistency \n", length)
			lock <- true
			result, err := sendSomethingToAllServers("0")
			if err == nil {
				log.Println("We are consistent with amount: ", result)
				return result, nil
			}
		}
	} else {
		//Print the 1 response
		for key, _ := range mapOfResponse {
			log.Printf("Consistent answer: %v \n", key)
			lock <- true
			return key, nil
		}
		//If map contains no values (should not happen), return -1
		return "-1", nil
	}
}
