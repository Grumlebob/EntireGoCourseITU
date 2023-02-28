package main

import (
    "context"
    "log"
    "math/rand"
    "time"
    "google.golang.org/grpc"
    pb "example.com/go-distributed-hashtable/DistributedHashtable"
)

var (
	address = "localhost"
    mclient pb.HashtableClient
    rclient pb.HashtableClient
)

func put(key int32, value int32) bool {
    errOnMain := false
    errOnRep := false
    var mstream pb.Hashtable_PutClient
    var rstream pb.Hashtable_PutClient
    var err error
    // put on main server
    if mclient != nil {
        mstream, err = mclient.Put(context.Background())
        if err != nil {
            log.Printf("open main stream error %v\n", err)
            errOnMain = true
        } else {
            // wait for lock on main server
            _, err := mstream.Recv()
            if err != nil {
                log.Printf("did not get lock on main server %v\n", err)
                errOnMain = true
            }
        } 
    } else {
        errOnMain = true
    }

    // put on replication server
    if rclient != nil {
        rstream, err = rclient.Put(context.Background())
        if err != nil {
            log.Printf("open replication stream error %v\n", err)
            errOnRep = true
        } else {
            // wait for lock on replication server
            _, err := rstream.Recv()
            if err != nil {
                log.Printf("did not get lock on replication server %v:\n", err)
                errOnRep = true
            }
        }
    } else {
        errOnRep = true
    }

    if errOnRep && errOnMain {
        log.Fatalf("both servers failed\n")
    }
    // create put request
    putReq := pb.PutMsg{Key: key, Value: value}

    // send put request to replication server
    if !errOnRep {
        if err := rstream.Send(&putReq); err != nil {
            log.Printf("can not send %v\n", err)
            errOnRep = true
        } else {
            // recieve reply from replication server
            rrsp, err := rstream.Recv()
            if err != nil {
                log.Printf("can not get reply from replication server %v\n", err)
                errOnRep = true
            } else {
                errOnRep = !rrsp.GetSucces()
            }

        }
    }

    // send put request to main server 
    if !errOnMain {
        if err := mstream.Send(&putReq); err != nil {
            log.Printf("can not send %v\n", err)
            errOnMain = true
        } else {
            // recieve reply from main server 
            mrsp, err := mstream.Recv()
            if err != nil {
                log.Printf("can not get reply from main server %v\n", err)
                errOnMain = true
            } else {
                errOnMain = !mrsp.GetSucces()
            }
        }
    }
    succes := !(errOnMain && errOnRep)
    log.Printf("put(%v, %v) -> %v\n", key, value, succes)
    return succes
}
func get(key int32) int32 {
    // create get request
    getReq := pb.GetMsg{Key: key}
    
    var value int32
    mainReached := false
    if mclient != nil {
        rsp, err := mclient.Get(context.Background(), &getReq)
        if err == nil {
            value = rsp.GetValue()
            mainReached = true
        } else {
            log.Printf("failed to get from main server %v", err)
        }
    }
    if !mainReached && rclient != nil {
        rsp, err := rclient.Get(context.Background(), &getReq)
        if err == nil {
            value = rsp.GetValue()
        } else {
            log.Printf("failed to get from replication server %v", err)
        }
    }
    log.Printf("get(%v) -> %v\n", key, value)
    return value
}
func main() {
    // set seed for generating random requests
    rand.Seed(time.Now().Unix())
    
    errOnMain := false
    errOnRep := false
    
    // dial main server
    mconn, err := grpc.Dial(":50001", grpc.WithInsecure())
    if err != nil {
		log.Printf("can not connect with main server %v\n", err)
        errOnMain = true
	} 

    // dial replication server
    rconn, err := grpc.Dial(":50002", grpc.WithInsecure())
    if err != nil {
		log.Printf("can not connect with replication server %v\n", err)
        errOnRep = true
	}
    
    if errOnRep && errOnMain {
        log.Fatalf("failed to connect with both servers\n")
    }
    // create clients
    if !errOnMain {
        mclient = pb.NewHashtableClient(mconn)
    }
    if !errOnRep {
        rclient = pb.NewHashtableClient(rconn)
    }
	
    // do 10 random put calls
    for i := 0; i < 10; i++{
        put(int32(rand.Intn(10)),int32(rand.Intn(10)))
        time.Sleep(time.Second)
    }
    // do 10 random get calls
    for i := 0; i < 10; i++{
        get(int32(rand.Intn(10)))
        time.Sleep(time.Second)
    }

}
