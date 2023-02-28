package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	pb "jacob.com/go-grpc/grpc"

	"google.golang.org/grpc"
)

type Server struct {
	pb.ReplicaServiceServer
}

var (
	amount int64
	mu     sync.Mutex
	lock   = make(chan bool, 1)
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 4000

	log.SetFlags(log.Ltime)

	// Create listener tcp on port 4000, 4001, 4002
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	log.Printf("Hello, This is server %v", ownPort)
	grpcServer := grpc.NewServer()
	pb.RegisterReplicaServiceServer(grpcServer, &Server{})

	amount = 0
	lock <- true
	// Kill 1 server after 20 seconds
	if ownPort == 4000 {
		go func() {
			time.Sleep(20 * time.Second)
			os.Exit(1)
		}()
	}

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (s *Server) SendSomeMessage(ctx context.Context, in *pb.SomeMessage) (*pb.SomeMessage, error) {

	log.Printf("Received: %v", in.Message)

	response := &pb.SomeMessage{
		Message: in.Message,
	}

	return response, nil
}

func (s *Server) Increment(ctx context.Context, in *pb.RequestIncrement) (*pb.ResponseIncrement, error) {
	<-lock

	log.Printf("Received: %v \n", in.IncrementAmount)

	//return error, amount is negative
	if in.IncrementAmount < 0 {
		return &pb.ResponseIncrement{CurrentAmount: amount}, errors.New("amount is negative")
	}

	amount += in.IncrementAmount
	log.Printf("Total: %v \n", amount)
	response := &pb.ResponseIncrement{
		CurrentAmount: amount,
	}
	lock <- true
	return response, nil
}
