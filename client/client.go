package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	Address string
	Dir     string
}

func New(address string, dir string) (*Client, error) {
	d, err := os.Stat(dir)

	if err != nil {
		return nil, err
	}

	if !d.IsDir() {
		return nil, fmt.Errorf("Path must be to directory not a file")
	}

	return &Client{Address: address, Dir: dir}, nil
}

func (client *Client) Run() {
	if len(client.Dir) < 1 {
		log.Println("Error: file path not provided")
		return
	}

	entries, err := os.ReadDir(client.Dir)

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
		filepath := client.Dir + "/" + entry.Name()

		if !entry.IsDir() {
			sendFile(conn, filepath)
		}

	}
}

func getRelativeFiles(dir string, files []string) []string {
	for index, file := range files {
		files[index] = file[len(dir)+1:]
	}

	return files
}

func getFiles(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Println("Error scanning directory: ", err)
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			childDir := dir + "/" + entry.Name()
			filesInChildDir, err := getFiles(childDir)
			if err != nil {
				return nil, err
			}
			files = append(files, filesInChildDir...)
		} else {
			file := dir + "/" + entry.Name()
			files = append(files, file)
		}
	}

	return files, nil
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
	time.Sleep(time.Second)  // sleeping for now but client needs to keep TCP open until server confirms file received
}
