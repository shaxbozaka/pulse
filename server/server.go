package main

import (
	"context"
	"fmt"
	"log"
	"net"

	controlpb "Pulse/gen/protos/control"
	datapb "Pulse/gen/protos/data"

	"google.golang.org/grpc"
	"io"
)

// Server struct implementing both gRPC services
type Server struct {
	controlpb.UnimplementedTunnelControlServer
	datapb.UnimplementedTunnelDataServer
}

// CreateTunnel handles tunnel creation requests
func (s *Server) CreateTunnel(ctx context.Context, req *controlpb.TunnelRequest) (*controlpb.TunnelResponse, error) {
	tunnelID := fmt.Sprintf("tunnel-%s", req.ClientId)
	publicURL := fmt.Sprintf("static.115.48.21.65.clients.your-server.de:%d", 5000)

	log.Printf("Created tunnel %s -> %s:%d", tunnelID, req.TargetHost, req.TargetPort)
	return &controlpb.TunnelResponse{TunnelId: tunnelID, PublicUrl: publicURL}, nil
}

// ForwardData handles data streaming between client and server
func (s *Server) ForwardData(stream datapb.TunnelData_ForwardDataServer) error {
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("Forwarding data for tunnel %s", packet.TunnelId)
		stream.Send(packet) // Echo back for now
	}
	return nil
}

func main() {
	// Start listening on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create new gRPC server
	grpcServer := grpc.NewServer()
	controlpb.RegisterTunnelControlServer(grpcServer, &Server{})
	datapb.RegisterTunnelDataServer(grpcServer, &Server{})

	log.Println("gRPC Server started on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
