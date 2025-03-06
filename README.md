# Pulse - HTTP Tunneling System

Pulse is a robust tunneling system that allows you to expose your local HTTP services to the internet securely.

## Quick Start

### 1. Deploy the Remote Server

You have two options to deploy the remote server:

#### Option A: Using Docker

```bash
# Build the Docker image
docker build -t pulse-server .

# Run the server
docker run -d -p 50051:50051 --name pulse-server pulse-server
```

#### Option B: Direct Deployment

```bash
# Build the server binary
go build -o pulse-server server/server.go

# Run the server
./pulse-server
```

### 2. Run the Client

The client connects to your remote server and forwards requests to your local service:

```bash
# Build the client
go build -o pulse-client client/client.go

# Run the client (replace SERVER_IP with your remote server's IP)
./pulse-client -server SERVER_IP:50051 -local localhost:3000
```

## Configuration

### Server Configuration
- Default port: 50051 (gRPC)
- Environment variables:
  - `PULSE_PORT`: Override default port
  - `PULSE_TIMEOUT`: Connection timeout (default: 30s)

### Client Configuration
- Default local service: localhost:3000
- Default tunnel ID: tunnel-client-123
- Environment variables:
  - `PULSE_SERVER`: Remote server address
  - `PULSE_LOCAL_SERVICE`: Local service address
  - `PULSE_TUNNEL_ID`: Custom tunnel ID

## Security Considerations

1. In production:
   - Enable TLS encryption
   - Implement authentication
   - Use secure tunnel IDs
   - Monitor connections

2. Firewall rules:
   - Allow incoming traffic on port 50051 (server)
   - Allow outgoing connections to your remote server (client)

## Troubleshooting

1. Connection Issues:
   ```bash
   # Check if local service is running
   curl localhost:3000

   # Verify server is accessible
   telnet SERVER_IP 50051

   # Check logs
   docker logs pulse-server  # if using Docker
   ```

2. Common Problems:
   - Port 3000 already in use: Change local service port
   - Connection refused: Check firewall rules
   - Timeout: Check network connectivity

## Development

1. Prerequisites:
   - Go 1.24+
   - Protocol Buffers compiler
   - Docker (optional)

2. Build from source:
   ```bash
   # Generate protobuf code
   protoc --go_out=. --go-grpc_out=. protos/*.proto

   # Build both components
   ./build.sh
   ```

## Support

For issues and feature requests, please create an issue in the repository.
