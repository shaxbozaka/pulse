package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	remoteServer = "localhost:9000" // Simulated remote server
	localService = "localhost:3000" // The actual local HTTP service
	maxRetries   = 5
	retryDelay   = 2 * time.Second
)

// handleConnection manages bidirectional data transfer between connections
func handleConnection(remote, local net.Conn) {
	defer remote.Close()
	defer local.Close()

	log.Printf("🔗 Tunnel established: %s ↔️ %s", remote.RemoteAddr(), local.RemoteAddr())

	var wg sync.WaitGroup
	wg.Add(2)

	// Forward remote → local
	go func() {
		defer wg.Done()
		n, err := io.Copy(local, remote)
		if err != nil {
			log.Printf("❌ Error forwarding remote → local: %v", err)
		} else {
			log.Printf("✅ Forwarded %d bytes from remote → local", n)
		}
	}()

	// Forward local → remote
	go func() {
		defer wg.Done()
		n, err := io.Copy(remote, local)
		if err != nil {
			log.Printf("❌ Error forwarding local → remote: %v", err)
		} else {
			log.Printf("✅ Forwarded %d bytes from local → remote", n)
		}
	}()

	wg.Wait()
	log.Printf("🔓 Connection closed: %s ↔️ %s", remote.RemoteAddr(), local.RemoteAddr())
}

// connectWithRetry attempts to establish a connection with retries
func connectWithRetry(addr string) (net.Conn, error) {
	for i := 0; i < maxRetries; i++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			return conn, nil
		}
		log.Printf("⚠️ Failed to connect to %s (attempt %d/%d): %v", addr, i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

func main() {
	log.Printf("🚀 Starting tunnel client...")
	log.Printf("ℹ️ Will connect to remote server at %s", remoteServer)
	log.Printf("ℹ️ Will forward to local service at %s", localService)

	// Connect to remote server
	remoteConn, err := connectWithRetry(remoteServer)
	if err != nil {
		log.Fatalf("❌ Failed to connect to remote server: %v", err)
	}
	log.Printf("✅ Connected to remote server at %s", remoteConn.RemoteAddr())

	// Connect to local service
	localConn, err := connectWithRetry(localService)
	if err != nil {
		remoteConn.Close()
		log.Fatalf("❌ Failed to connect to local service: %v", err)
	}
	log.Printf("✅ Connected to local service at %s", localConn.RemoteAddr())

	// Handle the connection
	handleConnection(remoteConn, localConn)
}
