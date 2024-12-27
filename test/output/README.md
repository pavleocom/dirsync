# Dirsync

A lightweight app to sync files from a source directory (on the client) to a target directory (on the server) over a TCP connection. Ideal for non-production use, such as syncing files from a main machine to a Raspberry Pi.

## Usage

### Server
  
Build the server: `go build -o server server.go`  
Run the server: `./server test/output` (Files received will be saved in the `test/output` directory.)  

### Client

Edit the server IP and PORT in `client.go`.  
Build the client: `go build -o client client.go`  
Run the client: `./client test/input` (All files in the `test/input` directory will be sent to the server; nested directories are not supported.)  

## Future Improvements

- Support syncing files in nested directories.  
- Add TLS for secure file transfer.  
- Convert the client into a daemon to monitor directory changes and send updated files.  
- Escape delimiters to avoid potential bugs.  
