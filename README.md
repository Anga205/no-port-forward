# Node.js TCP/HTTP Tunnel

This project sets up a TCP/HTTP tunnel using Node.js. The tunnel allows you to forward network traffic from a client machine to a server with a public IP, making it appear as if the service is running locally on the client machine.

## Features

- **TCP Tunnel**: Forward any TCP traffic.
- **HTTP Tunnel**: Forward HTTP traffic.

## Prerequisites

- Node.js installed on both the client and server machines.

## Setup

### Server Setup

1. **Install Node.js**: Ensure Node.js is installed on your server.
2. **Save Server Code**: Create a file named `server.js` and paste the following code:

    ```javascript
    const net = require('net');

    const LOCAL_PORT = 80;  // The port your actual server is running on
    const TUNNEL_PORT = 9999;  // The port that will accept tunnel connections

    const server = net.createServer((clientSocket) => {
        const remoteSocket = net.createConnection({ host: 'localhost', port: LOCAL_PORT });

        clientSocket.on('data', (data) => {
            remoteSocket.write(data);
        });

        remoteSocket.on('data', (data) => {
            clientSocket.write(data);
        });

        clientSocket.on('close', () => {
            remoteSocket.end();
        });

        remoteSocket.on('close', () => {
            clientSocket.end();
        });
    });

    server.listen(TUNNEL_PORT, () => {
        console.log(`Server listening on port ${TUNNEL_PORT}`);
    });
    ```

3. **Run the Server**: Execute the following command to start the server:

    ```sh
    node server.js
    ```

### Client Setup

1. **Install Node.js**: Ensure Node.js is installed on your client machine.
2. **Save Client Code**: Create a file named `client.js` and paste the following code:

    ```javascript
    const net = require('net');

    const LOCAL_PORT = 8080;  // Port on the client machine that will be used to access the service
    const REMOTE_HOST = 'your.server.public.ip';  // Public IP of your server
    const REMOTE_PORT = 9999;  // The port on the server that the tunnel listens on

    const server = net.createServer((clientSocket) => {
        const remoteSocket = net.createConnection({ host: REMOTE_HOST, port: REMOTE_PORT });

        clientSocket.on('data', (data) => {
            remoteSocket.write(data);
        });

        remoteSocket.on('data', (data) => {
            clientSocket.write(data);
        });

        clientSocket.on('close', () => {
            remoteSocket.end();
        });

        remoteSocket.on('close', () => {
            clientSocket.end();
        });
    });

    server.listen(LOCAL_PORT, () => {
        console.log(`Client listening on port ${LOCAL_PORT}`);
    });
    ```

3. **Run the Client**: Execute the following command to start the client (replace `your.server.public.ip` with your server's actual public IP):

    ```sh
    node client.js
    ```

## Usage

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

