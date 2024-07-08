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
