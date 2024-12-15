package main

import (
	"log"
	"os"

	"github.com/pavleocom/dirsync/server"
)

func main() {

	syncDir := "."

	if len(os.Args) > 1 {
		syncDir = os.Args[1]
	}

	serverInstance, err := server.New(syncDir)

	if err != nil {
		log.Println("Error starting server: ", err)
		return
	}

	serverInstance.Run()
}
