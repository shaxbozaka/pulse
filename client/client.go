package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	datapb "Pulse/gen/protos/data" // Adjust this import path
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	tunnelID          = "tunnel-client-123"
	reconnectBaseTime = 1 * time.Second  // Start with a shorter reconnect time
	reconnectMaxTime  = 30 * time.Second // Match server's max connection age
)

// proxyRequest forwards a request to a local HTTP server and builds a response packet
func proxyRequest(packet *datapb.DataPacket) (*datapb.DataPacket, error) {
	if packet == nil {
		return nil, fmt.Errorf("received nil packet")
	}
	// Ensure we have a valid URL
	if packet.Url == "" {
		return &datapb.DataPacket{
			TunnelId:  packet.TunnelId,
			RequestId: packet.RequestId,
			Type:      datapb.PacketType_RESPONSE,
			Status:    http.StatusBadRequest,
			Headers:   make(map[string]string),
			Data:      []byte("Empty URL in request"),
		}, nil
	}

	// Use the URL from the packet directly since it should already contain the full target URL
	requestURL, err := url.Parse(packet.Url)
	if err != nil {
		return &datapb.DataPacket{
			TunnelId:  packet.TunnelId,
			RequestId: packet.RequestId,
			Type:      datapb.PacketType_RESPONSE,
			Status:    http.StatusBadRequest,
			Headers:   make(map[string]string),
			Data:      []byte(fmt.Sprintf("Invalid URL: %v", err)),
		}, nil
	}

	// Create request with proper method and body
	req, err := http.NewRequest(packet.Method, requestURL.String(), bytes.NewReader(packet.Data))
	if err != nil {
		return &datapb.DataPacket{
			TunnelId:  packet.TunnelId,
			RequestId: packet.RequestId,
			Type:      datapb.PacketType_RESPONSE,
			Status:    http.StatusInternalServerError,
			Headers:   make(map[string]string),
			Data:      []byte(fmt.Sprintf("Error creating request: %v", err)),
		}, nil
	}

	// Copy headers from the original request
	for key, value := range packet.Headers {
		req.Header.Set(key, value)
	}

	// Set up client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		statusCode := http.StatusBadGateway
		if err, ok := err.(net.Error); ok && err.Timeout() {
			statusCode = http.StatusGatewayTimeout
		}
		return &datapb.DataPacket{
			TunnelId:  packet.TunnelId,
			RequestId: packet.RequestId,
			Type:      datapb.PacketType_RESPONSE,
			Status:    int32(statusCode),
			Headers:   make(map[string]string),
			Data:      []byte(fmt.Sprintf("Error forwarding request: %v", err)),
		}, nil
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &datapb.DataPacket{
			TunnelId:  packet.TunnelId,
			RequestId: packet.RequestId,
			Type:      datapb.PacketType_RESPONSE,
			Status:    http.StatusInternalServerError,
			Headers:   make(map[string]string),
			Data:      []byte(fmt.Sprintf("Error reading response body: %v", err)),
		}, nil
	}

	// Create response packet
	responsePacket := &datapb.DataPacket{
		TunnelId:  packet.TunnelId,
		RequestId: packet.RequestId, // Preserve the request ID
		Type:      datapb.PacketType_RESPONSE,
		Status:    int32(resp.StatusCode),
		Headers:   make(map[string]string),
		Data:      respBody,
	}

	// Copy all response headers
	for key, values := range resp.Header {
		responsePacket.Headers[key] = strings.Join(values, ", ")
	}

	return responsePacket, nil
}

// forwardData manages the gRPC stream and reconnection logic
func forwardData(client datapb.TunnelDataClient) {
	reconnectDelay := reconnectBaseTime

	for {
		// Create context for the stream
		ctx, cancel := context.WithCancel(context.Background())
		log.Printf("὎1 Establishing stream with server...")
		stream, err := client.ForwardData(ctx)
		if err != nil {
			log.Printf("❌ Failed to create stream: %v. Retrying in %v...", err, reconnectDelay)
			time.Sleep(reconnectDelay)
			reconnectDelay = min(reconnectDelay*2, reconnectMaxTime)
			cancel()
			continue
		}

		// Send initial packet to establish the tunnel
		initialPacket := &datapb.DataPacket{
			TunnelId: tunnelID,
			Type:     datapb.PacketType_KEEPALIVE,
			Data:     []byte("initial connection"),
			Url:      "",                      // Ensure URL is empty for initial packet
			Method:   "",                      // Ensure method is empty for initial packet
			Headers:  make(map[string]string), // Initialize empty headers
		}
		if err := stream.Send(initialPacket); err != nil {
			log.Printf("❌ Failed to send initial packet: %v. Retrying...", err)
			cancel()
			continue
		}
		log.Printf("ὑ0 Initial connection packet sent successfully")

		// Reset delay and log successful connection
		reconnectDelay = reconnectBaseTime
		log.Printf("✅ Successfully connected to server with tunnel ID: %s", tunnelID)

		// Handle incoming packets in a goroutine
		go func() {
			defer cancel()
			for {
				packet, err := stream.Recv()
				if err == io.EOF {
					log.Println("🔌 Stream closed by server (EOF)")
					return
				}
				if err != nil {
					log.Printf("⚠️ Error receiving data: %v", err)
					return
				}

				log.Printf("📦 Received packet: Type=%v, URL=%s", packet.Type, packet.Url)

				switch packet.Type {
				case datapb.PacketType_REQUEST:
					log.Printf("➡️ Processing REQUEST: %s", packet.Url)
					respPacket, err := proxyRequest(packet)
					if err != nil {
						log.Printf("❌ Proxy error: %v", err)
						continue
					}
					if err := stream.Send(respPacket); err != nil {
						log.Printf("❌ Error sending response: %v", err)
						return
					}
				case datapb.PacketType_KEEPALIVE:
					log.Println("💓 Received KEEPALIVE")
				case datapb.PacketType_RESPONSE:
					log.Println("⚠️ Unexpected RESPONSE packet received, ignoring")
				}
			}
		}()

		// Send keepalive packets periodically (match server's ping interval)
		keepaliveTicker := time.NewTicker(5 * time.Second)
		defer keepaliveTicker.Stop()
		log.Printf("ὐd Connected to server with tunnel ID: %s", tunnelID)

		for {
			select {
			case <-ctx.Done():
				log.Println("🔄 Stream closed. Reconnecting...")
				return
			case <-keepaliveTicker.C:
				err := stream.Send(&datapb.DataPacket{
					TunnelId: tunnelID,
					Type:     datapb.PacketType_KEEPALIVE,
					Data:     []byte("keepalive"),
				})
				if err != nil {
					log.Printf("❌ Stream error: %v. Reconnecting...", err)
					return // Exit the loop to trigger reconnection
				}
				log.Printf("Ὁ3 Sent keepalive")
			}
		}

		cancel()
	}
}

// min returns the smaller of two durations
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func main() {
	serverAddress := "localhost:9000" // Connect to remote server

	// Configure keepalive parameters to match server's expectations
	kacp := keepalive.ClientParameters{
		Time:                5 * time.Second, // Send pings every 5 seconds
		Timeout:             2 * time.Second, // Wait 2 seconds for ping ack
		PermitWithoutStream: true,            // Allow pings without active streams
	}

	// Configure connection parameters
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(), // Use WithTransportCredentials for production
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{"service": "TunnelData"}],
				"waitForReady": true,
				"retryPolicy": {
					"maxAttempts": 5,
					"initialBackoff": "1s",
					"maxBackoff": "10s",
					"backoffMultiplier": 2,
					"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`),
	}

	log.Printf("὎1 Connecting to server at %s...", serverAddress)

	// Connect with retry
	var conn *grpc.ClientConn
	var err error
	for attempt := 1; ; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err = grpc.DialContext(ctx, serverAddress, dialOpts...)
		cancel()

		if err == nil {
			log.Printf("✅ Successfully connected to server after %d attempt(s)", attempt)
			break
		}

		backoff := time.Duration(math.Min(float64(attempt*attempt), 30)) * time.Second
		log.Printf("❌ Failed to connect (attempt %d): %v. Retrying in %v...", attempt, err, backoff)
		time.Sleep(backoff)
	}
	defer conn.Close()

	client := datapb.NewTunnelDataClient(conn)
	forwardData(client)
}
