package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Grumlebob/GoLangAssignment2HardgRPC/protos"

	"google.golang.org/grpc"
)

type Server struct {
	protos.ChatServiceServer
}

func (s *Server) GetHeader(ctx context.Context, message *protos.Message) (*protos.Message, error) {
	var msgServer = message
	//If recieving first handshake
	if (message.Seq == 0) && (message.Ack == 0) {
		fmt.Println("Server recieved first handshake, with Syn flag True and Seq 0")
		msgServer = &protos.Message{Text: "Second handshake  sent from Server, with Syn flag True and Ack 1", Seq: message.Seq, Ack: message.Seq + 1}
		fmt.Printf("Server sending second handshake with Ack: %d \n", msgServer.Ack)
	} else {
		fmt.Printf("Server acknowledged it got %d and processed it. \n", message.Seq)
		msgServer = &protos.Message{Seq: message.Seq, Ack: message.Seq + 1}
		fmt.Printf("Server expect next to be Seq %d \n", msgServer.Ack)
	}
	return msgServer, nil
}

func main() {
	// Create listener tcp on port 9080
	listener, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	protos.RegisterChatServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
