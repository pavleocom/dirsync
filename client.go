package main

import (
	"log"
	"os"

	"github.com/pavleocom/dirsync/client"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Must provide file as argument")
	}

	file := os.Args[1]

	clientInstance := client.New("127.0.0.1:7007")
	clientInstance.Run(file)

}
