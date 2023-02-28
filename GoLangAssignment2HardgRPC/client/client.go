package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/Grumlebob/GoLangAssignment2HardgRPC/protos"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := protos.NewChatServiceClient(conn)

	tcpSimulation(c)
}

func tcpSimulation(c protos.ChatServiceClient) {
	message := protos.Message{Text: "Client sent first handshake, with Syn flag True and Seq 0", Ack: 0, Seq: 0}
	fmt.Println(message.Text)

	//First handshake
	firstHandshake, err := c.GetHeader(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling GetHeader(Message): %s", err)
	}
	//Wait for second handshake
	for firstHandshake.Ack != 1 {

	}
	fmt.Printf("Client recieved second handshake from server with Syn flag True and Ack: %d, and Seq: %d \n", firstHandshake.Ack, firstHandshake.Seq)

	//Third handshake
	message = protos.Message{Text: "Client sent third hardshake, with Ack: 1 and Seq: ", Ack: firstHandshake.Ack, Seq: firstHandshake.Seq + 1}
	fmt.Println(message.Text, message.Seq)
	thirdHandshake, err := c.GetHeader(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling GetHeader(Message): %s", err)
	}
	fmt.Printf("Client Recieved from server, Ack: %d \n", thirdHandshake.Ack)

	//Data exhange logic here
	for i := 0; i < 10; i++ {
		message.Seq++
		fmt.Printf("Client Sent to server fictional data with Seq: %d \n", message.Seq)
		dataSimulation, err := c.GetHeader(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling GetHeader(Message): %s", err)
		}
		fmt.Printf("Client recieved from server Ack: %d \n", dataSimulation.Ack)
	}
}
