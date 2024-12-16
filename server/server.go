package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type Server struct {
	Address string
	Dir     string
}

func New(address string, dir string) (*Server, error) {
	d, err := os.Stat(dir)

	if err != nil {
		return nil, err
	}

	if !d.IsDir() {
		return nil, fmt.Errorf("Path must be to directory not a file")
	}

	return &Server{Address: address, Dir: dir}, nil
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", server.Address)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Printf("Server listening for connections on %s...\n", server.Address)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		go server.HandleConnection(conn)
	}
}

func (server *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		err := receiveFile(server, reader)

		if err == io.EOF {
			break
		}
	}
}

func receiveFile(server *Server, reader *bufio.Reader) error {
	_, err := reader.ReadBytes(byte(0x1E))

	if err == io.EOF {
		return err
	}

	if err != nil {
		log.Println("Error reading delimiter: ", err)
		return err
	}

	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading file name: ", err)
		return err
	}

	fileName = fileName[:len(fileName)-1]

	fileSizeStr, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading file size: ", err)
		return err
	}

	fileSizeStr = fileSizeStr[:len(fileSizeStr)-1]

	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		log.Println("Error parsing file size: ", err)
		return err
	}

	fmt.Printf("Receiving file: %s (%d bytes)\n", fileName, fileSize)

	newFile, err := os.Create(server.Dir + "/" + fileName)
	if err != nil {
		log.Println("Error creating file: ", err)
		return err
	}

	defer newFile.Close()

	buffer := make([]byte, 1024)
	var leftOverBytes []byte

	for {
		if len(leftOverBytes) > 0 {
			fmt.Println(leftOverBytes)
			reader = bufio.NewReader(io.MultiReader(bytes.NewBuffer(leftOverBytes), reader))
			leftOverBytes = nil
		}

		n, err := reader.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading from connection: ", err)
		}

		data := buffer[:n]

		eofIndex := bytes.IndexByte(data, byte(0x1E))

		if eofIndex != -1 {
			_, err := newFile.Write(data[:eofIndex])

			if err != nil {
				log.Print("Error writing into file when delimiter found: ", err)
				return err
			}

			leftOverBytes = data[eofIndex+1:]
			break
		}

		_, err = newFile.Write(data)

		if err != nil {
			log.Println("Erro writing to file: ", err)
			return err
		}
	}

	fmt.Printf("File %s received succesfully!\n", fileName)

	return nil
}
