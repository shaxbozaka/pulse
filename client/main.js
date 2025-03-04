import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import net from 'net';
import { program } from 'commander';

program
.option('-s, --server <url>', 'gRPC server address')
.option('-p, --port <number>', 'Local port to bind', '3000')
.parse(process.argv);

const options = program.opts();
const SERVER_ADDRESS = options.server;
const LOCAL_PORT = parseInt(options.port, 10);

if (!SERVER_ADDRESS) {
console.error('Error: gRPC server address is required');
process.exit(1);
}

const PROTO_PATH = './proto/bridge.proto';
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
keepCase: true,
longs: String,
enums: String,
defaults: true,
oneofs: true,
});

const bridgeProto = grpc.loadPackageDefinition(packageDefinition).bridge;
const client = new bridgeProto.BridgeService(SERVER_ADDRESS, grpc.credentials.createInsecure());

const server = net.createServer((socket) => {
console.log(`Local client connected on port ${LOCAL_PORT}`);

socket.on('data', (data) => {
client.forwardData({ data: data.toString() }, (err, response) => {
if (err) {
console.error('Error forwarding data:', err);
return;
}
socket.write(response.data);
});
});

socket.on('close', () => {
console.log('Local client disconnected');
});
});

server.listen(LOCAL_PORT, () => {
console.log(`Listening on localhost:${LOCAL_PORT}`);
});
