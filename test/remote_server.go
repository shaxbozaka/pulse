package main

import (
	"io"
	"log"
	"net"
)

const listenAddr = ":80" // Public port for incoming connections

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}
	log.Println("Remote server listening on", listenAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept client:", err)
			continue
		}

		log.Println("New connection from:", clientConn.RemoteAddr())

		// Wait for the local agent to connect
		localConn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to connect to local agent:", err)
			clientConn.Close()
			continue
		}

		log.Println("Forwarding connection...")

		go func() {
			io.Copy(localConn, clientConn) // Forward client → local agent
			io.Copy(clientConn, localConn) // Forward local agent → client
		}()
	}
}
