package main

import (
	datapb "Pulse/gen/protos/data"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
)

type Client struct {
	datapb.UnimplementedTunnelDataServer
}

// ForwardData receives data from the server and forwards it to localhost
// forwardData handles data forwarding between the server and the local application.
func forwardData(client datapb.TunnelDataClient) {
	stream, err := client.ForwardData(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.CloseSend()

	log.Println("Client waiting for requests...")

	for {
		// Receive request from the server
		packet, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			continue
		}

		log.Printf("Received request for %s, forwarding to localhost", packet.Data)

		// Forward request to the local server (localhost:3000)
		resp, err := http.Get(fmt.Sprintf("http://localhost:3000%s", string(packet.Data)))
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			continue
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			continue
		}

		// Send response back to the gRPC server
		err = stream.Send(&datapb.DataPacket{
			TunnelId: packet.TunnelId,
			Data:     body,
		})
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
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
	forwardData(client) // Now correctly implemented
}
