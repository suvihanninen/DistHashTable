package main

import (

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	hashTable "github.com/suvihanninen/DistHashTable.git/grpc"
	"google.golang.org/grpc"
)

func main() {

	//Make connection to 5001
	port1 := ":5001"
	connection1, err := grpc.Dial(port1, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}

	//Make connection to 5002
	port2 := ":5002"
	connection2, err := grpc.Dial(port2, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	//log to file instead of console
	f := setLogClient()
	defer f.Close()

	server1 := hashTable.NewHashTableClient(connection1) //creates a new client
	server2 := hashTable.NewHashTableClient(connection2) //creates a new client

	defer connection1.Close()
	defer connection2.Close()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		println("Enter 'put [key] [value]' to insert a value to HashTable, or 'get [key]' to retrieve a value.")

		for {
			scanner.Scan()
			text := scanner.Text()

			if strings.Contains(text, "put") {
				input := strings.Fields(text)

				key, err := strconv.ParseInt(input[1], 10, 32)
				if err != nil {
					log.Fatalf("Problem with key: %v", err)
				}

				value, err := strconv.ParseInt(input[2], 10, 32)
				if err != nil {
					log.Fatalf("Problem with value: %v", err)
				}

				putRequest := &hashTable.PutRequest{
					Key:   int32(key),
					Value: int32(value),
				}

				response1, err := server1.Put(context.Background(), putRequest)
				if err != nil {
					log.Printf("Response from server %s: Put failed: ", port1, err)
				}

				response2, err := server2.Put(context.Background(), putRequest)
				if err != nil {
					log.Printf("Response from server %s: Put failed: ", port2, err)
				}

				if response1.GetResponse() == true {
					log.Printf("Response from server %s: Put response: %s", port1, response1.GetResponse())
					println("Response from server "+port1+": Put response: ", response1.GetResponse())
				} else {
					//log.Printf("Put failed on server %s . The response: ", port1, response1.GetResponse())
					println("Response from server "+port1+": Put response: ", response1.GetResponse())
				}

				if response2.GetResponse() == true {
					log.Printf("Response from server %s: Put response: %s", port2, response2.GetResponse())
					println("Response from server "+port2+": Put response: ", response2.GetResponse())
				} else {
					//log.Printf("Put failed on server %s . The response: ", port2, response2.GetResponse())
					println("Response from server "+port2+": Put response: ", response2.GetResponse())
				}

			} else if strings.Contains(text, "get") {
				input := strings.Fields(text)
				success1 := true
				success2 := true
				key, err := strconv.ParseInt(input[1], 10, 32)
				if err != nil {
					log.Fatalf("Problem with key: %v", err)
				}

				getRequest := &hashTable.GetRequest{
					Key: int32(key),
				}

				response1, err := server1.Get(context.Background(), getRequest)
				if err != nil {
					log.Printf("Response from server %s: Get failed: ", port1, err)
					success1 = false
				}
				if success1 {
					log.Printf("Response from server %s: Get response: ", port1, response1.GetValue())
					println("Response from server "+port1+": Get response: ", response1.GetValue())
				} else {
					//server 5001 has crached
				}

				response2, err := server2.Get(context.Background(), getRequest)
				if err != nil {
					log.Printf("Response from server %s: Get failed: ", port2, err)
					success2 = false
				}

				if success2 {
					log.Printf("Response from server %s: Get response: ", port2, response2.GetValue())
					println("Response from server "+port2+": Get response: ", response2.GetValue())
				} else {
					//server 5002 has crached
				}

			} else {
				println("Sorry didn't catch that, try again ")
			}
		}
	}()

	for {

	}

}

// sets the logger to use a log.txt file instead of the console
func setLogClient() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
