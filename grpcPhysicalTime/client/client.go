package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Grumlebob/grpcPhysicalTime/protos"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := protos.NewChatServiceClient(conn)

	timeSync(c)
}

func timeSync(c protos.ChatServiceClient) {
	t := time.Now()
	time := timestamppb.New(t)

	clientRequest := protos.ClientRequest{
		Text:      "Client sent first handshake, with Syn flag True and Seq 0",
		Timestamp: time,
	}
	fmt.Println(clientRequest.Timestamp)

	//First handshake
	firstHandshake, err := c.GetTime(context.Background(), &clientRequest)
	if err != nil {
		log.Fatalf("Error when calling GetTime : %s", err)
	}
	fmt.Println(firstHandshake)

}
