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
