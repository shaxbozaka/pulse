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

	s.activeClients.Store(tunnelID, true)
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

		log.Printf("Received packet: type=%s, tunnel=%s, data=%s", packet.Type, packet.TunnelId, string(packet.Data))

		// Only echo KEEPALIVE packets
		if packet.Type == datapb.PacketType_KEEPALIVE {
			response := &datapb.DataPacket{
				TunnelId: packet.TunnelId,
				Type:     datapb.PacketType_KEEPALIVE,
				Data:     packet.Data,
			}
			if err := stream.Send(response); err != nil {
				log.Printf("Error sending keepalive response: %v", err)
				return err
			}
		}
		// REQUEST and RESPONSE packets are handled by handleHTTP, not echoed here
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
	if !s.clientAvailable() {
		http.Error(w, "No active clients", http.StatusServiceUnavailable)
		return
	}

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

	tunnelID := "tunnel-client-123" // Use a consistent tunnel ID
	err = stream.Send(&datapb.DataPacket{
		TunnelId: tunnelID,
		Type:     datapb.PacketType_REQUEST,
		Data:     []byte(r.URL.Path),
		Method:   r.Method,
		Url:      r.URL.String(),
	})
	if err != nil {
		log.Printf("Failed to send request to gRPC client: %v", err)
		http.Error(w, "Failed to send data", http.StatusInternalServerError)
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		log.Printf("Failed to receive response from gRPC client: %v", err)
		http.Error(w, "Failed to receive response", http.StatusInternalServerError)
		return
	}

	if resp.Type != datapb.PacketType_RESPONSE {
		log.Printf("Unexpected packet type: %s", resp.Type)
		http.Error(w, "Invalid response type", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(int(resp.Status))
	for key, value := range resp.Headers {
		w.Header().Set(key, value)
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
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	http.HandleFunc("/", server.handleHTTP)
	log.Println("HTTP Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
