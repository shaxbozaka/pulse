package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	datapb "Pulse/gen/protos/data"
	"google.golang.org/grpc"
)

// Client struct to handle gRPC communication
type Client struct {
	datapb.UnimplementedTunnelDataServer
}

// ForwardData receives data from the server and forwards it to localhost
func (c *Client) ForwardData(stream datapb.TunnelData_ForwardDataServer) error {
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("Received request for %s, forwarding to localhost", packet.Data)

		// Forward to actual application running on localhost:8080
		resp, err := http.Get(fmt.Sprintf("http://localhost:3000%s", string(packet.Data)))
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			return stream.Send(&datapb.DataPacket{
				TunnelId: packet.TunnelId,
				Data:     []byte("Error forwarding request"),
			})
		}
		defer resp.Body.Close()

		// Read response
		body, _ := io.ReadAll(resp.Body)

		// Send response back to server
		stream.Send(&datapb.DataPacket{
			TunnelId: packet.TunnelId,
			Data:     body,
		})
	}
	return nil
}

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("65.21.48.115:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := datapb.NewTunnelDataClient(conn)

	// Keep connection open
	stream, err := client.ForwardData(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Keep the client running
	for {
		err := stream.Send(&datapb.DataPacket{TunnelId: "tunnel-client-123", Data: []byte("Ping")})
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}
	}
}
