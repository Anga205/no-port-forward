package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	LOCAL_PORT  = 80   // The port your actual server is running on
	TUNNEL_PORT = 9999 // The port that will accept tunnel connections
)

func handleClient(client net.Conn) {
	defer client.Close()

	remote, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", LOCAL_PORT))
	if err != nil {
		log.Printf("Failed to connect to remote server: %v", err)
		return
	}
	defer remote.Close()

	go io.Copy(remote, client)
	io.Copy(client, remote)
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", TUNNEL_PORT))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on port %d", TUNNEL_PORT)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}

		go handleClient(client)
	}
}
