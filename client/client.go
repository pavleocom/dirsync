package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	Address string
}

func New(address string) *Client {
	return &Client{Address: address}
}

func (client *Client) Run(dir string) {

	if len(dir) < 1 {
		log.Println("Error: file path not provided")
		return
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Println("Error scanning directory: ", err)
		return
	}

	conn, err := net.Dial("tcp", client.Address)

	defer conn.Close()

	if err != nil {
		log.Println("Error connecting to server: ", err)
		return
	}

	for _, entry := range entries {
		filepath := dir + "/" + entry.Name()

		if !entry.IsDir() {
			sendFile(conn, filepath)
		}

	}
}

func sendFile(conn net.Conn, filePath string) {
	conn.Write([]byte{0x1E}) //sending delimiter 0x1E

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		log.Println("Error reading file info: ", err)
		return
	}

	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()

	_, err = conn.Write([]byte(fmt.Sprintf("%s\n", fileName)))

	if err != nil {
		log.Println("Error sending file name: ", err)
		return
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%d\n", fileSize)))

	if err != nil {
		log.Println("Error sending file size: ", err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}

	defer file.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading file: ", err)
			continue
		}

		if n > 0 {
			conn.Write(buffer[:n])
		}
	}

	conn.Write([]byte{0x1E}) //sending delimiter 0x1E
}
