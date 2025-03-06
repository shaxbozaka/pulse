package main

import (
	"io"
	"log"
	"net"
)

const (
	remoteServer = "static.115.48.21.65.clients.your-server.de:9000" // Remote server
	localService = "localhost:3000"                                  // Local HTTP server
)

func main() {
	// Connect to remote server
	conn, err := net.Dial("tcp", remoteServer)
	if err != nil {
		log.Fatalf("Failed to connect to remote server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to remote server, waiting for requests...")

	for {
		localConn, err := net.Dial("tcp", localService)
		if err != nil {
			log.Println("Failed to connect to local service:", err)
			continue
		}

		go func() {
			io.Copy(localConn, conn)
			io.Copy(conn, localConn)
		}()
	}
}
