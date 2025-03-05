package main

import (
	"context"
	"fmt"
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return stream.Send(&datapb.DataPacket{
				TunnelId: packet.TunnelId,
				Data:     []byte("Error reading response body"),
			})
		}

		// Send response back to server
		stream.Send(&datapb.DataPacket{
			TunnelId: packet.TunnelId,
			Data:     body,
		})
	}
	return nil
}

// Handle HTTP requests and forward them to the client
// Handle HTTP requests and forward them to the client
func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP request: %s", r.URL.Path)

	// Connect to the client via gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Connect to client
	if err != nil {
		log.Printf("Failed to connect to client: %v", err)
		http.Error(w, "Failed to connect to client", http.StatusInternalServerError)
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

	// Send actual HTTP request data to client
	err = stream.Send(&datapb.DataPacket{
		TunnelId: "tunnel-client-123",
		Data:     []byte(r.URL.Path), // Sending the HTTP request path
	})
	if err != nil {
		log.Printf("Failed to send data to client: %v", err)
		http.Error(w, "Failed to send data to client", http.StatusInternalServerError)
		return
	}

	// Receive response from client
	resp, err := stream.Recv()
	if err != nil {
		log.Printf("Failed to receive data from client: %v", err)
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
