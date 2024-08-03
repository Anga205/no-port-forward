# ❗THIS PROJECT IS NO LONGER MAINTAINED❗


## TCP/HTTP Tunnel

This project sets up a TCP/HTTP tunnel using Go. The tunnel allows you to forward network traffic from a client machine to a server with a public IP, making it appear as if the service is running locally on the client machine.

## Features

- **TCP Tunnel**: Forward any TCP traffic.
- **HTTP Tunnel**: Forward HTTP traffic.

## Prerequisites

- Go compiler installed on both the client and server machines.

## Usage

On server:
```bash
go run service.go
```

on open-network proxy
```bash
go run proxy.go
```


After setting up and running both the server and client, you can test the tunnel:

1. **Open a Browser**: Navigate to `http://localhost:8080` on your client machine.
2. **Using cURL**: Execute the following command on your client machine:

    ```sh
    curl http://localhost:8080
    ```

The traffic will be tunneled through the client to the server and then forwarded to the actual service running on the server's local port (80).

## Troubleshooting

- Ensure that the server and client machines can reach each other over the network.
- Verify that the ports are correctly configured and not blocked by firewalls.

