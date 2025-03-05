package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	controlpb "Pulse/gen/protos/control"
	datapb "Pulse/gen/protos/data"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server implements both gRPC services
type Server struct {
	controlpb.UnimplementedTunnelControlServer
	datapb.UnimplementedTunnelDataServer
	activeClients sync.Map // Stores active clients (key: tunnelID, value: true)
}

// CreateTunnel handles tunnel creation requests
func (s *Server) CreateTunnel(ctx context.Context, req *controlpb.TunnelRequest) (*controlpb.TunnelResponse, error) {
	tunnelID := fmt.Sprintf("tunnel-%s", req.ClientId)
	publicURL := fmt.Sprintf("static.115.48.21.65.clients.your-server.de:%d", 5000)

	// Store active client in sync.Map
	s.activeClients.Store(tunnelID, true)
	log.Printf("Stored active client: %s", tunnelID) // Debugging log

	s.activeClients.Range(func(key, value interface{}) bool {
		log.Printf("Active client found: %v", key) // Print every stored client
		return true
	})
	log.Printf("Created tunnel: %s -> %s:%d", tunnelID, req.TargetHost, req.TargetPort)
	return &controlpb.TunnelResponse{TunnelId: tunnelID, PublicUrl: publicURL}, nil
}

// ForwardData handles bidirectional streaming for relaying packets
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

		// Send received data to HTTP response channel
		response := &datapb.DataPacket{
			TunnelId: packet.TunnelId,
			Data:     packet.Data,
		}

		if err := stream.Send(response); err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}

// clientAvailable checks if at least one active tunnel exists
func (s *Server) clientAvailable() bool {
	var available bool
	s.activeClients.Range(func(key, value interface{}) bool {
		available = true
		return false // Stop iteration as we found an active client
	})
	return available
}

// handleHTTP forwards HTTP requests to the gRPC client if available
func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP request: %s", r.URL.Path)
	log.Print(s.clientAvailable())
	if s.clientAvailable() {
		// Forward request to gRPC client
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect to gRPC server: %v", err)
			http.Error(w, "gRPC connection failed", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := datapb.NewTunnelDataClient(conn)
		stream, err := client.ForwardData(context.Background())
		if err != nil {
			log.Printf("Failed to create gRPC stream: %v", err)
			http.Error(w, "Failed to create gRPC stream", http.StatusInternalServerError)
			return
		}
		defer stream.CloseSend()

		tunnelID := fmt.Sprintf("tunnel-client-%d", r.ContentLength)
		err = stream.Send(&datapb.DataPacket{
			TunnelId: tunnelID,
			Data:     []byte(r.URL.Path),
		})
		if err != nil {
			log.Printf("Failed to send data to gRPC client: %v", err)
			http.Error(w, "Failed to send data", http.StatusInternalServerError)
			return
		}

		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive response from gRPC client: %v", err)
			http.Error(w, "Failed to receive response", http.StatusInternalServerError)
			return
		}

		w.Write(resp.Data)
	} else {
		// Return a direct HTTP response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No active client. Showing local response."))
	}
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
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Start HTTP server
	http.HandleFunc("/", server.handleHTTP)
	log.Println("HTTP Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
