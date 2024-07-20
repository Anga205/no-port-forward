package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	connections = make(map[net.Conn]bool)
	mutex       sync.Mutex
)

func main() {
	go startSocketServer()
	startHTTPServer()
}

func startSocketServer() {
	// Listen on port 3000 for socket connections
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println("Error setting up the socket server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Socket server is listening on port 3000")

	for {
		// Accept a connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		mutex.Lock()
		connections[conn] = true
		mutex.Unlock()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		mutex.Lock()
		delete(connections, conn)
		mutex.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		// Read data from the connection
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by peer")
			} else {
				fmt.Println("Error reading from connection:", err)
			}
			return
		}
		// Print received message
		fmt.Print("Received message: ", message)

		// Close the connection if the message is "exit"
		if strings.TrimSpace(message) == "exit" {
			fmt.Println("Closing connection as 'exit' was received")
			return
		}
	}
}

func startHTTPServer() {
	http.HandleFunc("/", handleHTTPRequest)
	fmt.Println("HTTP server is listening on port 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error setting up the HTTP server:", err)
	}
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	// Serialize the request to send to the client
	var requestBuilder strings.Builder
	requestBuilder.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.RequestURI, r.Proto))
	for name, values := range r.Header {
		for _, value := range values {
			requestBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", name, value))
		}
	}
	requestBuilder.WriteString("\r\n")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	requestBuilder.Write(body)
	fmt.Println("Serialized request:", requestBuilder.String())

	// Send the serialized request to all active connections
	mutex.Lock()
	var responseBody strings.Builder
	for conn := range connections {
		fmt.Println("Sending request to connection")
		_, err := conn.Write([]byte(requestBuilder.String() + "\n"))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			conn.Close()
			delete(connections, conn)
			continue
		}

		// Read the response from the client
		connReader := bufio.NewReader(conn)
		for {
			line, err := connReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Connection closed by client")
					break
				}
				fmt.Println("Error reading from connection:", err)
				break
			}
			responseBody.WriteString(line)
		}
	}
	mutex.Unlock()

	// Respond to the HTTP request with the received message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBody.String()))
}
