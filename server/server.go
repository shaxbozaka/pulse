package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	controlpb "Pulse/gen/protos/control"
	datapb "Pulse/gen/protos/data"

	"google.golang.org/grpc"
)

// Server struct implementing both gRPC services
type Server struct {
	controlpb.UnimplementedTunnelControlServer
	datapb.UnimplementedTunnelDataServer
	client datapb.TunnelDataClient // Store gRPC client to forward requests
}

// CreateTunnel handles tunnel creation requests
func (s *Server) CreateTunnel(ctx context.Context, req *controlpb.TunnelRequest) (*controlpb.TunnelResponse, error) {
	tunnelID := fmt.Sprintf("tunnel-%s", req.ClientId)
	publicURL := fmt.Sprintf("static.115.48.21.65.clients.your-server.de:%d", 5000)

	log.Printf("Created tunnel %s -> %s:%d", tunnelID, req.TargetHost, req.TargetPort)
	return &controlpb.TunnelResponse{TunnelId: tunnelID, PublicUrl: publicURL}, nil
}

// ForwardData handles streaming between client and server
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

// Handle incoming HTTP requests and forward to the client
func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP request: %s", r.URL.Path)

	// Forward the request to the client via gRPC
	stream, err := s.client.ForwardData(context.Background())
	if err != nil {
		http.Error(w, "Failed to create gRPC stream", http.StatusInternalServerError)
		return
	}

	// Send request to client
	err = stream.Send(&datapb.DataPacket{
		TunnelId: "tunnel-client-123",
		Data:     []byte(r.URL.Path), // Just sending the path for now
	})
	if err != nil {
		http.Error(w, "Failed to send data to client", http.StatusInternalServerError)
		return
	}

	// Receive response from client
	resp, err := stream.Recv()
	if err != nil {
		http.Error(w, "Failed to receive data from client", http.StatusInternalServerError)
		return
	}

	// Send response back to HTTP client
	w.Write(resp.Data)
}

func main() {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{}
	controlpb.RegisterTunnelControlServer(grpcServer, server)
	datapb.RegisterTunnelDataServer(grpcServer, server)

	go func() {
		log.Println("gRPC Server started on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Start HTTP Server
	http.HandleFunc("/", server.handleHTTP)
	log.Println("HTTP Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
