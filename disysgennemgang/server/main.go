package main

import (
  "context"
  "log"
  "net"
  "fmt"
  pb "example.com/go-distributed-hashtable/DistributedHashtable"
  "google.golang.org/grpc"
)

var (
  lock = make(chan bool,1)
  ht = make(map[int32]int32)
)

type HashtableServer struct {
  pb.UnimplementedHashtableServer
} 

func (s *HashtableServer) Put(stream pb.Hashtable_PutServer) error{
  // lock
  <-lock
  
  // tell client they have lock with empty response message
  if err := stream.Send(&pb.PutRsp{}); err != nil {  
    lock<-true
    log.Printf("failed to send lock accept: %v\n", err)
    return nil
  } 

  // recieve put request
  req, err := stream.Recv() 
  if err != nil {
    lock<-true
    log.Printf("failed to recieve put request: %v\n", err)
    return nil
  }

  // put
  log.Printf("Put: %v, %v\n", req.GetKey(), req.GetValue())
  ht[req.GetKey()] = req.GetValue()

  // return succes
  if err := stream.Send(&pb.PutRsp{Succes: true}); err != nil {  
    lock<-true
    log.Printf("failed to send succes respoce: %v\n", err)
    return nil

  } 

  // unlock
  lock<-true
  return nil
}

func (s *HashtableServer) Get(ctx context.Context, in *pb.GetMsg) (*pb.GetRsp, error){
  // lock
  <-lock
  
  log.Printf("Get: %v\n", in.GetKey())
  
  val := ht[in.GetKey()]
  
  // unlock
  lock<-true
  return &pb.GetRsp{Value: val}, nil
}

func main() {
  port := ""
  //unlock
  lock<-true 
  
  fmt.Printf("listen on port[50001, 50002]: \n\t>>>")
  fmt.Scanln(&port)
  lis, err := net.Listen("tcp", ":" + port)
  if err != nil {
    log.Fatalf("failed to listen: %v\n", err)
  }

  s := grpc.NewServer()
  pb.RegisterHashtableServer(s, &HashtableServer{})
  log.Printf("server listening at %v\n", port)
  if err:= s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v\n", err)
  }

}
