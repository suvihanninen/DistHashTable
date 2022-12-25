package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	hashTable "github.com/suvihanninen/DistHashTable.git/grpc"
	"google.golang.org/grpc"
)

type HashTableServer struct {
	hashTable.UnimplementedHashTableServer        // You need this line if you have a server
	port                                   string // Not required but useful if your server needs to know what port it's listening to
	ctx                                    context.Context
	ht                                     map[int32]int32
	lock                                   chan bool
}

func main() {
	port := os.Args[1]
	address := ":" + port
	list, err := net.Listen("tcp", address)

	if err != nil {
		log.Printf("FEServer %s: Server on port %s: Failed to listen on port %s: %v", port, port, address, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	grpcServer := grpc.NewServer()

	//log to file instead of console
	f := setLog()
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := &HashTableServer{
		port: os.Args[1],
		ctx:  ctx,
		ht:   make(map[int32]int32),
		lock: make(chan bool, 1),
	}

	//unlock
	server.lock <- true
	hashTable.RegisterHashTableServer(grpcServer, server)

	go func() {
		log.Printf("Server %s: We are trying to listen calls from client: %s", server.port, port)
		println("Server is listening")
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to serve %v", err)
		}

		log.Printf("FEServer %s: We have started to listen calls from client: %s", server.port, port)
	}()

	for {
	}
}

func (server *HashTableServer) Put(ctx context.Context, PutRequest *hashTable.PutRequest) (*hashTable.PutResponse, error) {
	<-server.lock
	time.Sleep(10 * time.Second)
	success := true

	server.ht[PutRequest.GetKey()] = PutRequest.GetValue()
	for key, val := range server.ht {
		log.Printf("Key: %v and Value: %v", key, val)
	}
	server.lock <- true
	return &hashTable.PutResponse{Response: success}, nil

}

func (server *HashTableServer) Get(ctx context.Context, GetRequest *hashTable.GetRequest) (*hashTable.GetResponse, error) {
	//lock
	<-server.lock
	//time.Sleep(5 * time.Second)
	for key, val := range server.ht {

		if key == GetRequest.GetKey() {
			log.Printf("We found a value!")
			server.lock <- true
			return &hashTable.GetResponse{Value: val}, nil
		}

	}
	server.lock <- true
	return &hashTable.GetResponse{Value: 0}, nil
}

func setLog() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
