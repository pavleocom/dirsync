package main

import (
	"log"
	"os"

	"github.com/pavleocom/dirsync/client"
)

func main() {

	dir := "."

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	c, err := client.New("127.0.0.1:7007", dir)

	if err != nil {
		log.Fatal("Error starting client: ", err)
		return
	}

	c.Run()
}
