package main

import (
	"context"
	pb "example.com/Replication/repService"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
)

var (
	address = "localhost"
	ports   = []int{50001, 50002, 50003, 50004, 50005, 50006, 50007, 50008, 50009}
)

var (
	id = uuid.New()
)

func main() {
	ctx, conn, conerr := Connect()
	var client pb.ReplicationClient
	if conerr != nil {
		log.Fatalf("could not connect to a server: %v", conerr)
	}

	defer conn.Close()
	for {
		client = pb.NewReplicationClient(conn)
		x := int32(rand.Int() % 10)
		y := int32(rand.Int() % 10)
		var putted, err = client.Put(ctx, &pb.KeyVal{Key: x, Val: y})
		if err != nil {
			ctx2, conn2, errc := Connect()
			if errc != nil {
				log.Fatalf("could not get result: %v", err)
			}
			conn = conn2
			ctx = ctx2
			continue
		}
		log.Printf("succesfully putted %v, %v : %v", x, y, putted.Res)
		x = int32(rand.Int() % 10)
		var got, errr = client.Get(ctx, &pb.Int{X: x})
		if errr != nil {
			ctx3, conn3, err3 := Connect()
			if err3 != nil {
				log.Fatalf("something went wrong: %v", errr)
			}
			conn = conn3
			ctx = ctx3
			continue
		}
		log.Printf("succesfully got %v : %v", x, got.X)

		time.Sleep(time.Millisecond * 100)
	}
}

func Connect() (context.Context, *grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var connErr error
	var timeoutCtx context.Context

	for i := range ports {
		timeoutCtx, _ = context.WithTimeout(context.Background(), 1*time.Second)
		conn, connErr = grpc.DialContext(timeoutCtx, fmt.Sprintf("%s:%d", address, ports[i]), grpc.WithInsecure(), grpc.WithBlock())
		if connErr != nil {
			log.Printf("Connection failed %d: %v", ports[i], connErr)
			continue
		}
		break
	}
	return timeoutCtx, conn, connErr
}
