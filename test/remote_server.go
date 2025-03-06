package main

import (
	"io"
	"log"
	"net"
	"time"
)

const (
	listenAddr = ":80" // Public port for incoming connections
	timeout    = 30 * time.Second
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("🔗 New connection from %s", conn.RemoteAddr())

	// Buffer for reading HTTP request
	buf := make([]byte, 4096)
	for {
		// Read request
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("❌ Error reading from connection: %v", err)
			}
			return
		}

		// Connect to local service
		localConn, err := net.DialTimeout("tcp", "localhost:3000", 5*time.Second)
		if err != nil {
			log.Printf("❌ Failed to connect to local service: %v", err)
			// Send error response back to client
			errorResp := []byte("HTTP/1.1 502 Bad Gateway\r\nConnection: close\r\n\r\nFailed to connect to local service")
			conn.Write(errorResp)
			return
		}

		// Forward the request we've read
		_, err = localConn.Write(buf[:n])
		if err != nil {
			log.Printf("❌ Error writing to local service: %v", err)
			localConn.Close()
			return
		}

		// Read response from local service
		respBuf := make([]byte, 4096)
		for {
			n, err := localConn.Read(respBuf)
			if err != nil {
				if err != io.EOF {
					log.Printf("❌ Error reading from local service: %v", err)
				}
				break
			}

			// Forward response back to client
			_, err = conn.Write(respBuf[:n])
			if err != nil {
				log.Printf("❌ Error writing response to client: %v", err)
				break
			}
		}

		localConn.Close()
		log.Printf("✅ Request processed successfully")
	}
}

func main() {
	// Create listener
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("❌ Failed to listen on %s: %v", listenAddr, err)
	}
	defer listener.Close()

	log.Printf("🚀 Remote server listening on %s", listenAddr)
	log.Printf("ℹ️ Will forward requests to local service at localhost:3000")

	for {
		// Accept client connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("❌ Failed to accept connection: %v", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}
