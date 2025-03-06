package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	datapb "Pulse/gen/protos/data" // Update this to your actual proto import path
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// tunnelDataServer handles gRPC and HTTP tunneling
type tunnelDataServer struct {
	datapb.UnimplementedTunnelDataServer
	clients       map[string]chan *datapb.DataPacket // Map tunnelID → Request Channel
	responses     map[string]chan *datapb.DataPacket // Map requestID → Response Channel
	activeStreams map[string]context.CancelFunc      // Map tunnelID → Stream cancel function
	mu            sync.RWMutex                       // Read-Write Mutex for concurrency
	requestIDMu   sync.Mutex                         // Mutex for requestID generation
	requestID     int64                              // Counter for generating unique request IDs
	timeout       time.Duration                      // Timeout for request/response operations
}

// NewTunnelDataServer creates a new server instance
func NewTunnelDataServer() *tunnelDataServer {
	return &tunnelDataServer{
		clients:       make(map[string]chan *datapb.DataPacket),
		responses:     make(map[string]chan *datapb.DataPacket),
		activeStreams: make(map[string]context.CancelFunc),
		timeout:       30 * time.Second, // Increased timeout to match connection settings
	}
}

// ForwardData handles gRPC communication from clients
func (s *tunnelDataServer) ForwardData(stream datapb.TunnelData_ForwardDataServer) error {
	// Get the first packet to identify the client
	packet, err := stream.Recv()
	if err != nil {
		log.Printf("Failed to receive initial packet: %v", err)
		return fmt.Errorf("initial packet error: %v", err)
	}
	if packet == nil {
		log.Println("Received nil initial packet")
		return fmt.Errorf("nil initial packet")
	}
	if packet.TunnelId == "" {
		log.Println("Client sent empty tunnel ID. Rejecting.")
		return fmt.Errorf("invalid tunnel ID")
	}

	tunnelID := packet.TunnelId
	log.Printf("Client attempting to connect with ID: %s", tunnelID)

	// Handle existing connection for this tunnel ID
	s.mu.Lock()
	if cancel, exists := s.activeStreams[tunnelID]; exists {
		log.Printf("Found existing stream for tunnel ID: %s, cleaning up...", tunnelID)
		cancel()
		delete(s.activeStreams, tunnelID)
		// Wait a moment for cleanup
		time.Sleep(100 * time.Millisecond)
	}

	// Create a new context and cancel function for this stream
	streamCtx, cancel := context.WithCancel(context.Background())
	s.activeStreams[tunnelID] = cancel

	// Create a channel for this client
	requestChan := make(chan *datapb.DataPacket, 10)
	s.clients[tunnelID] = requestChan
	s.mu.Unlock()

	// Start a goroutine to monitor stream context
	go func() {
		<-streamCtx.Done()
		s.mu.Lock()
		delete(s.activeStreams, tunnelID)
		s.mu.Unlock()
		log.Printf("Stream context cancelled for tunnel ID: %s", tunnelID)
	}()

	// Ensure cleanup on disconnect
	defer func() {
		s.mu.Lock()
		// Close and remove all response channels for this client
		for id, ch := range s.responses {
			if strings.HasPrefix(id, tunnelID) {
				if ch != nil {
					close(ch)
				}
				delete(s.responses, id)
			}
		}
		// Close and remove the request channel
		if ch := s.clients[tunnelID]; ch != nil {
			close(ch)
		}
		delete(s.clients, tunnelID)

		// Cancel and remove the stream context
		if cancel, exists := s.activeStreams[tunnelID]; exists {
			cancel()
			delete(s.activeStreams, tunnelID)
		}
		s.mu.Unlock()
		log.Printf("Client %s disconnected and cleaned up", tunnelID)
	}()

	// Log successful connection
	log.Printf("✅ Client %s successfully connected and stream established", tunnelID)

	// Handle incoming packets
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Ὂ5 Client %s disconnected (EOF)", tunnelID)
			return nil
		}
		if err != nil {
			log.Printf("❌ Error receiving packet from %s: %v", tunnelID, err)
			return fmt.Errorf("stream receive error: %v", err)
		}
		if packet == nil {
			log.Printf("⚠️ Received nil packet from %s", tunnelID)
			continue
		}

		log.Printf("὎5 Received packet from %s: Type=%v, RequestID=%s", tunnelID, packet.Type, packet.RequestId)

		switch packet.Type {
		case datapb.PacketType_RESPONSE:
			// Validate response packet
			if err := validateResponsePacket(packet); err != nil {
				log.Printf("❌ Invalid response packet from %s: %v", tunnelID, err)
				continue
			}
			log.Printf("✅ Valid response packet received for request %s", packet.RequestId)

			// Get the response channel for this request
			s.mu.RLock()
			respChan, ok := s.responses[packet.RequestId]
			s.mu.RUnlock()

			if ok && respChan != nil {
				// Send response back to HTTP handler with timeout
				select {
				case respChan <- packet:
					log.Printf("Sent response for request %s", packet.RequestId)
				case <-time.After(s.timeout):
					log.Printf("Timeout sending response for request %s", packet.RequestId)
				}
			} else {
				log.Printf("No response channel found for request %s", packet.RequestId)
			}
		case datapb.PacketType_KEEPALIVE:
			log.Printf("Received keepalive from %s", tunnelID)
		default:
			log.Printf("Unknown packet type from %s: %v", tunnelID, packet.Type)
		}
	}
}
func (s *tunnelDataServer) httpHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Get Tunnel ID from query params or header
	tunnelID := r.Header.Get("X-Tunnel-ID")
	if tunnelID == "" {
		tunnelID = r.URL.Query().Get("tunnel_id")
	}
	if tunnelID == "" {
		tunnelID = "tunnel-client-123" // Match the client's default tunnel ID
	}

	// Check if a client is connected first
	s.mu.RLock()
	requestChan, ok := s.clients[tunnelID]
	s.mu.RUnlock()

	if !ok {
		http.Error(w, "No client connected for tunnel ID: "+tunnelID, http.StatusServiceUnavailable)
		return
	}

	// Generate unique request ID
	s.requestIDMu.Lock()
	s.requestID++
	requestID := fmt.Sprintf("%s-%d", tunnelID, s.requestID)
	s.requestIDMu.Unlock()

	// Create response channel for this request
	responseChan := make(chan *datapb.DataPacket, 1)
	s.mu.Lock()
	s.responses[requestID] = responseChan
	s.mu.Unlock()

	// Ensure cleanup of response channel
	defer func() {
		s.mu.Lock()
		delete(s.responses, requestID)
		close(responseChan)
		s.mu.Unlock()
	}()

	// Create the target URL for localhost
	targetURL := "http://localhost:3000" + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	// Create a DataPacket for the HTTP request
	packet := &datapb.DataPacket{
		TunnelId:  tunnelID,
		RequestId: requestID,
		Type:      datapb.PacketType_REQUEST,
		Method:    r.Method,
		Url:       targetURL,
		Headers:   make(map[string]string),
		Data:      body,
	}

	log.Printf("📤 Forwarding request %s to %s", requestID, targetURL)

	// Copy request headers, but skip some headers that might interfere
	for key, values := range r.Header {
		if len(values) > 0 {
			// Skip headers that might interfere with the forwarded request
			switch strings.ToLower(key) {
			case "host", "connection", "keep-alive", "proxy-authenticate", "proxy-authorization", "te", "trailers", "transfer-encoding", "upgrade":
				continue
			default:
				packet.Headers[key] = values[0]
			}
		}
	}

	// Send the request packet to the client with timeout
	select {
	case requestChan <- packet:
		// Successfully sent the packet
	case <-time.After(s.timeout):
		http.Error(w, "Request timed out while sending to client", http.StatusGatewayTimeout)
		return
	}

	// Wait for client response with timeout
	log.Printf("📬 Waiting for response to request %s", requestID)
	select {
	case resp := <-responseChan:
		if resp == nil {
			log.Printf("❌ Received nil response for request %s", requestID)
			http.Error(w, "Invalid response from client", http.StatusBadGateway)
			return
		}
		// Validate response
		if err := validateResponsePacket(resp); err != nil {
			log.Printf("❌ Invalid response for request %s: %v", requestID, err)
			http.Error(w, fmt.Sprintf("Invalid response from client: %v", err), http.StatusBadGateway)
			return
		}
		log.Printf("✅ Received valid response for request %s with status %d", requestID, resp.Status)

		// Set response headers
		for key, value := range resp.Headers {
			w.Header().Set(key, value)
		}

		// Ensure a valid status code, default to 200 if not set
		statusCode := int(resp.Status)
		if statusCode <= 0 {
			log.Printf("⚠️ No status code in response for request %s, defaulting to 200", requestID)
			statusCode = http.StatusOK
		}

		// Write status code and body
		w.WriteHeader(statusCode)
		if len(resp.Data) > 0 {
			if _, err := w.Write(resp.Data); err != nil {
				log.Printf("❌ Error writing response body for request %s: %v", requestID, err)
			} else {
				log.Printf("✅ Successfully wrote response for request %s (%d bytes)", requestID, len(resp.Data))
			}
		} else {
			log.Printf("ℹ️ Empty response body for request %s", requestID)
		}

	case <-time.After(s.timeout):
		http.Error(w, "Response timed out", http.StatusGatewayTimeout)
		return
	}
}

// validateResponsePacket performs validation checks on response packets
func validateResponsePacket(packet *datapb.DataPacket) error {
	if packet == nil {
		return fmt.Errorf("packet is nil")
	}
	if packet.RequestId == "" {
		return fmt.Errorf("missing request ID")
	}
	if packet.Status < 0 {
		return fmt.Errorf("invalid status code: %d", packet.Status)
	}
	return nil
}

func main() {
	// Start gRPC server with keepalive parameters
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on :50051: %v", err)
	}

	// Configure keepalive parameters
	kaep := keepalive.EnforcementPolicy{
		MinTime:             2 * time.Second, // More lenient minimum time between client pings
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Second, // More lenient idle timeout
		MaxConnectionAge:      60 * time.Second, // Allow connections to live longer
		MaxConnectionAgeGrace: 10 * time.Second, // More time for graceful shutdown
		Time:                  5 * time.Second,  // Keep ping interval at 5 seconds
		Timeout:               2 * time.Second,  // More time to wait for ping ack
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	)
	srv := NewTunnelDataServer()
	datapb.RegisterTunnelDataServer(grpcServer, srv)

	// Start HTTP server
	http.HandleFunc("/", srv.httpHandler)
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
