#!/bin/bash

echo "🔨 Building Pulse components..."

# Build remote server
echo "Building remote server..."
go build -o bin/pulse-server server/server.go

# Build client
echo "Building client..."
go build -o bin/pulse-client client/client.go

echo "✅ Build complete! Binaries are in the bin/ directory"
