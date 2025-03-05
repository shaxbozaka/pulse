package main

import (
	"context"
	"log"
	"time"

	controlpb "Pulse/gen/protos/control"
	datapb "Pulse/gen/protos/data"

	"google.golang.org/grpc"
	"io"
)

func main() {
	// Connect to the gRPC server
	//conn, err := grpc.Dial("static.115.48.21.65.clients.your-server.de:50051", grpc.WithInsecure())
	conn, err := grpc.Dial("static.115.48.21.65.clients.your-server.de:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	controlClient := controlpb.NewTunnelControlClient(conn)
	dataClient := datapb.NewTunnelDataClient(conn)

	// Create a tunnel
	createTunnel(controlClient)

	// Send and receive data
	forwardData(dataClient)
}

// Call CreateTunnel RPC
func createTunnel(client controlpb.TunnelControlClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &controlpb.TunnelRequest{
		ClientId:   "client-123",
		TargetHost: "localhost",
		TargetPort: 8080,
	}

	res, err := client.CreateTunnel(ctx, req)
	if err != nil {
		log.Fatalf("CreateTunnel failed: %v", err)
	}

	log.Printf("Tunnel created: ID=%s, PublicURL=%s", res.TunnelId, res.PublicUrl)
}

// Call ForwardData RPC
func forwardData(client datapb.TunnelDataClient) {
	stream, err := client.ForwardData(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

	// Send data
	go func() {
		for i := 0; i < 5; i++ {
			packet := &datapb.DataPacket{
				TunnelId: "tunnel-client-123",
				Data:     []byte("Hello from client!"),
			}

			if err := stream.Send(packet); err != nil {
				log.Fatalf("Failed to send data: %v", err)
			}

			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	// Receive data
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive data: %v", err)
		}

		log.Printf("Received data: %s", packet.Data)
	}
}
