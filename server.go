package main

import (
	"log"
	"os"

	"github.com/pavleocom/dirsync/server"
)

func main() {

	dir := "."

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	s, err := server.New("0.0.0.0:7007", dir)

	if err != nil {
		log.Println("Error starting server: ", err)
		return
	}

	s.Run()
}
