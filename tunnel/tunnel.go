package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "0.0.0.0:3000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Continuously read and process messages from the server
	reader := bufio.NewReader(conn)
	for {
		// Read the HTTP request line from the server
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed connection")
				return
			}
			fmt.Println("Error reading request line:", err)
			return
		}
		fmt.Println("Received request line:", requestLine)

		// Parse the request line
		requestParts := strings.Fields(requestLine)
		if len(requestParts) < 3 {
			fmt.Println("Invalid request line received:", requestLine)
			continue
		}
		method := requestParts[0]
		uri := requestParts[1]
		proto := strings.TrimSpace(requestParts[2])

		fmt.Println("Method:", method)
		fmt.Println("URI:", uri)
		fmt.Println("Protocol:", proto)

		// Read the headers
		headers := make(http.Header)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Server closed connection")
					return
				}
				fmt.Println("Error reading headers:", err)
				return
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			headerParts := strings.SplitN(line, ": ", 2)
			if len(headerParts) < 2 {
				fmt.Println("Invalid header line received:", line)
				continue
			}
			headers.Add(headerParts[0], headerParts[1])
		}

		fmt.Println("Headers:")
		for name, values := range headers {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}

		// Determine the body length or if it’s chunked
		var body io.Reader
		if lengthStr := headers.Get("Content-Length"); lengthStr != "" {
			length, err := parseContentLength(lengthStr)
			if err != nil {
				fmt.Println("Error parsing Content-Length:", err)
				return
			}
			body = io.LimitReader(reader, length)
		} else if headers.Get("Transfer-Encoding") == "chunked" {
			body = &chunkedReader{reader}
		} else {
			// No body is expected, set body to an empty reader
			body = io.NopCloser(strings.NewReader(""))
		}

		// Read the body
		var bodyBuilder strings.Builder
		_, err = io.Copy(&bodyBuilder, body)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading body:", err)
			return
		}
		fmt.Println("Body:", bodyBuilder.String())

		// Create a new HTTP request
		req, err := http.NewRequest(method, "http://127.0.0.1:6000"+uri, strings.NewReader(bodyBuilder.String()))
		if err != nil {
			fmt.Println("Error creating new request:", err)
			return
		}
		req.Header = headers
		req.Proto = proto // Set the protocol version

		// Print the created request for debugging
		fmt.Println("Created HTTP request:")
		fmt.Println(req.Method, req.URL, req.Proto)
		for name, values := range req.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}

		// Send the request to 127.0.0.1:6000
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request to 127.0.0.1:6000:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		fmt.Println("Response body from 127.0.0.1:6000:", string(respBody))

		// Send the response body back to the server
		_, err = conn.Write([]byte(string(respBody) + "\n"))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}

// Helper function to parse Content-Length header
func parseContentLength(value string) (int64, error) {
	var length int64
	_, err := fmt.Sscanf(value, "%d", &length)
	return length, err
}

// Helper type to handle chunked transfer encoding
type chunkedReader struct {
	rd *bufio.Reader
}

func (cr *chunkedReader) Read(p []byte) (int, error) {
	// Read chunk size
	line, err := cr.rd.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, io.EOF
	}

	var size int64
	_, err = fmt.Sscanf(line, "%x", &size)
	if err != nil {
		return 0, err
	}
	if size == 0 {
		return 0, io.EOF
	}

	// Read chunk data
	n, err := cr.rd.Read(p[:size])
	if err != nil && err != io.EOF {
		return n, err
	}

	// Discard CRLF after chunk
	_, err = cr.rd.ReadString('\n')
	return n, err
}
