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

/func forwardData(client datapb.TunnelDataClient) {
	stream, err := client.ForwardData(context.Background()) // Persistent connection
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.CloseSend() // Close only when exiting

	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by server")
			break
		}
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			continue
		}

		log.Printf("Received request for %s", packet.Data)

		// Forward request
		resp, err := http.Get(fmt.Sprintf("http://localhost:3000%s", string(packet.Data)))
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			continue
		}

		// Send response
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
