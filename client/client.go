package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	datapb "Pulse/gen/protos/data"
	"google.golang.org/grpc"
)

// ForwardData connects to the gRPC server and listens for data to forward
func forwardData(client datapb.TunnelDataClient) {
	stream, err := client.ForwardData(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Use a goroutine to receive messages
	go func() {
		for {
			packet, err := stream.Recv()
			if err == io.EOF {
				log.Println("Stream closed by server (EOF)")
				return
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}

			log.Printf("Received request for %s, forwarding to localhost", packet.Data)

			// Forward request to local server
			resp, err := http.Get(fmt.Sprintf("http://localhost:3000/%s", string(packet.Data)))
			if err != nil {
				log.Printf("Error forwarding request: %v", err)
				continue
			}
			defer resp.Body.Close()

			// Read response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v", err)
				continue
			}

			// Send response back to server
			err = stream.Send(&datapb.DataPacket{
				TunnelId: packet.TunnelId,
				Data:     body,
			})
			if err != nil {
				log.Printf("Error sending response: %v", err)
			}
		}
	}()

	// Keep sending heartbeat messages to keep the stream alive
	for {
		err := stream.Send(&datapb.DataPacket{TunnelId: "tunnel-client-123", Data: []byte("Ping")})
		if err == io.EOF {
			log.Println("Stream closed by server (EOF)")
			break
		}
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}
		time.Sleep(5 * time.Second) // Avoid spamming
	}
}

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("static.115.48.21.65.clients.your-server.de:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := datapb.NewTunnelDataClient(conn)

	log.Println("Client started, listening for incoming requests...")
	forwardData(client)
}
