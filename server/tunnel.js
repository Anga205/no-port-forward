import net from 'net';
const fs = require('fs');
const readline = require('readline');

const LOCAL_PORT = 8080;  // Port on the client machine that will be used to access the service

let serverDetails;

try {
    // Attempt to read server_details.json
    const data = fs.readFileSync('server_details.json');
    serverDetails = JSON.parse(data);
} catch (error) {
    // If server_details.json is not found, prompt user for input and store it in the json
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });

    rl.question('Enter the public IP of your server: ', (host) => {
        rl.question('Enter the port on the server that the tunnel listens on: ', (port) => {
            serverDetails = {
                REMOTE_HOST: host,
                REMOTE_PORT: parseInt(port)
            };

            fs.writeFileSync('server_details.json', JSON.stringify(serverDetails));
            rl.close();
        });
    });
}

const server1 = net.createServer((clientSocket) => {
    const remoteSocket = net.createConnection({ host: serverDetails.REMOTE_HOST, port: serverDetails.REMOTE_PORT });

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
