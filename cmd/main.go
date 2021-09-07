package main

import (
	"fmt"

	riak "github.com/tpjg/goriakpbc"
)

func main() {
	client := riak.New("127.0.0.1:8087")
	err := client.Connect()
	if err != nil {
		fmt.Println("Cannot connect, is Riak running?")
		return
	}

	bucket, _ := client.Bucket("tstriak")
	obj := bucket.New("tstobj")
	obj.ContentType = "application/json"
	obj.Data = []byte("{'field':'value'}")
	obj.Store()

	fmt.Printf("Stored an object in Riak, vclock = %v\n", obj.Vclock)

	client.Close()
}
