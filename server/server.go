package main

import (
	"fmt"
	"io"
	"log"
	"net"

	datapb "Pulse/gen/protos/data" // Adjust this import path
	"google.golang.org/grpc"
)

type tunnelDataServer struct {
	datapb.UnimplementedTunnelDataServer
}

// ForwardData implements the bidirectional streaming RPC
func (s *tunnelDataServer) ForwardData(stream datapb.TunnelData_ForwardDataServer) error {
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client closed the stream")
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving packet: %v", err)
		}

		// Handle packets based on their type
		switch packet.Type {
		case datapb.PacketType_REQUEST:
			// Echo back a simple response for demonstration
			response := &datapb.DataPacket{
				TunnelId: packet.TunnelId,
				Type:     datapb.PacketType_RESPONSE,
				Status:   200,
				Headers:  map[string]string{"Content-Type": "text/plain"},
				Data:     []byte("Response from server"),
			}
			if err := stream.Send(response); err != nil {
				return fmt.Errorf("error sending response: %v", err)
			}
		case datapb.PacketType_KEEPALIVE:
			// Echo back the keepalive packet
			if err := stream.Send(packet); err != nil {
				return fmt.Errorf("error sending keepalive: %v", err)
			}
		default:
			log.Printf("Received unknown packet type: %v", packet.Type)
		}
	}
}

func main() {
	// Listen on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on :50051: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()
	datapb.RegisterTunnelDataServer(s, &tunnelDataServer{})

	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
