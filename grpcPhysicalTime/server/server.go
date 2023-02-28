package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Grumlebob/grpcPhysicalTime/protos"

	"google.golang.org/grpc"
)

type Server struct {
	protos.ChatServiceServer
}

func (s *Server) GetTime(ctx context.Context, clientMessage *protos.ClientRequest) (*protos.ServerResponse, error) {
	var serverTime = time.Now().String()
	var response = &protos.ServerResponse{
		Text: serverTime,
	}

	fmt.Println("Client message: ", clientMessage)

	return response, nil
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
