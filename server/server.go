package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

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
	publicURL := fmt.Sprintf("%s:%d", os.Getenv("SERVER_IP"), 5000) // Use env variable

	log.Printf("Created tunnel %s -> %s:%d", tunnelID, req.TargetHost, req.TargetPort)
	return &controlpb.TunnelResponse{TunnelId: tunnelID, PublicUrl: publicURL}, nil
}

// ForwardData handles incoming gRPC stream
func (s *Server) ForwardData(stream datapb.TunnelData_ForwardDataServer) error {
	for {
		packet, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed the stream")
				return nil
			}
			log.Printf("Error receiving data: %v", err)
			return err
		}

		log.Printf("Received data from client: %s", string(packet.Data))

		err = stream.Send(&datapb.DataPacket{
			TunnelId: packet.TunnelId,
			Data:     []byte(fmt.Sprintf("Server received: %s", packet.Data)),
		})
		if err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}

// Handle HTTP requests and forward them to the client
func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP request: %s", r.URL.Path)

	// Reuse existing gRPC client
	stream, err := s.client.ForwardData(context.Background())
	if err != nil {
		log.Printf("Failed to create gRPC stream: %v", err)
		http.Error(w, "Failed to create gRPC stream", http.StatusInternalServerError)
		return
	}

	err = stream.Send(&datapb.DataPacket{TunnelId: "tunnel-client-123", Data: []byte(r.URL.Path)})
	if err != nil {
		log.Printf("Failed to send data to client: %v", err)
		http.Error(w, "Failed to send data to client", http.StatusInternalServerError)
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		log.Printf("Failed to receive data from client: %v", err)
		http.Error(w, "Failed to receive data from client", http.StatusInternalServerError)
		return
	}

	w.Write(resp.Data)
}

func main() {
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

	http.HandleFunc("/", server.handleHTTP)
	log.Println("HTTP Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
