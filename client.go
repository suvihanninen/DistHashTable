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

				if response1.GetResponse() == true {
					log.Printf("Response from server %s: Put response: %s", port1, response1.GetResponse())
					println("Response from server "+port1+": Put response: ", response1.GetResponse())
				} else {
					//log.Printf("Put failed on server %s . The response: ", port1, response1.GetResponse())
					println("Response from server "+port1+": Put response: ", response1.GetResponse())
				}

				response2, err := server2.Put(context.Background(), putRequest)
				if err != nil {
					log.Printf("Response from server %s: Put failed: ", port2, err)
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

				key, err := strconv.ParseInt(input[1], 10, 32)
				if err != nil {
					log.Fatalf("Problem with key: %v", err)
				}

				result := get(key, server1, server2, port1, port2)
				keyString := strconv.FormatInt(key, 10)
				log.Printf("Response to GET request when key is %v: %v", input, result)
				println("Response to GET request when key is "+keyString+": ", result)

			} else {
				println("Sorry didn't catch that, try again ")
			}
		}
	}()

	for {

	}

}

func get(key int64, server1 hashTable.HashTableClient, server2 hashTable.HashTableClient, port1 string, port2 string) int32 {
	var result int32
	success1 := true
	success2 := true
	result = 0

	getRequest := &hashTable.GetRequest{
		Key: int32(key),
	}

	response1, err := server1.Get(context.Background(), getRequest)
	if err != nil {
		log.Printf("Response from server %s: Get failed: ", port1, err)
		success1 = false
	}

	response2, err := server2.Get(context.Background(), getRequest)
	if err != nil {
		log.Printf("Response from server %s: Get failed: ", port2, err)
		success2 = false
	}

	if success2 && success1 {
		println("response1: ", response1.GetValue())
		println("response2: ", response2.GetValue())
		println("Success 1 and 2 are true")
		if response1.GetValue() == response2.GetValue() {
			result = response1.GetValue()
		} else {
			println("PUT was not updated yet on both replicas, we will call GET again")
			result = get(key, server1, server2, port1, port2)
		}
	} else if success1 {
		println("Success 1 is true")
		result = response1.GetValue()
	} else {
		println("Success 2 is true")
		result = response2.GetValue()
	}

	return result
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
