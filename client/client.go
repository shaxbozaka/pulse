package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	datapb "Pulse/gen/protos/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	tunnelID          = "tunnel-client-123"
	reconnectBaseTime = 5 * time.Second  // Initial reconnect delay
	reconnectMaxTime  = 60 * time.Second // Max wait time before retrying
)

// proxyRequest handles the full request forwarding
func proxyRequest(packet *datapb.DataPacket) (*datapb.DataPacket, error) {
	// Extract method and URL
	requestURL, err := url.Parse(fmt.Sprintf("http://localhost:3000/%s", url.QueryEscape(packet.Url)))
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Create new HTTP request
	req, err := http.NewRequest(packet.Method, requestURL.String(), bytes.NewReader(packet.Data))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	for key, value := range packet.Headers {
		req.Header.Set(key, value)
	}

	// Forward request to local server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error forwarding request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Prepare response packet
	responsePacket := &datapb.DataPacket{
		TunnelId: packet.TunnelId,
		Status:   int32(resp.StatusCode),
		Headers:  map[string]string{},
		Data:     respBody,
	}

	// Copy response headers
	for key, values := range resp.Header {
		responsePacket.Headers[key] = values[0] // Use first value
	}

	return responsePacket, nil
}

// forwardData handles gRPC streaming and request forwarding.
func forwardData(client datapb.TunnelDataClient) {
	reconnectDelay := reconnectBaseTime

	for {
		ctx, cancel := context.WithCancel(context.Background())

		// Establish stream
		stream, err := client.ForwardData(ctx)
		if err != nil {
			log.Printf("Failed to create stream: %v. Retrying in %v...", err, reconnectDelay)
			time.Sleep(reconnectDelay)
			reconnectDelay = min(reconnectDelay*2, reconnectMaxTime) // Exponential backoff
			continue
		}

		reconnectDelay = reconnectBaseTime // Reset on success

		// Receive and process requests
		go func() {
			defer cancel()
			for {
				packet, err := stream.Recv()
				if err == io.EOF {
					log.Println("Stream closed by server (EOF)")
					return
				}
				if err != nil {
					log.Printf("Error receiving data: %v", err)
					return
				}

				// Proxy request
				respPacket, err := proxyRequest(packet)
				if err != nil {
					log.Printf("Proxy error: %v", err)
					continue
				}

				// Send response back
				if err := stream.Send(respPacket); err != nil {
					log.Printf("Error sending response: %v", err)
					return
				}
			}
		}()

		// Keepalive pings
		for {
			select {
			case <-ctx.Done():
				log.Println("Stream closed. Reconnecting...")
				break
			default:
				err := stream.Send(&datapb.DataPacket{TunnelId: tunnelID, Data: []byte("keepalive")})
				if err != nil {
					log.Printf("Stream error: %v. Reconnecting...", err)
					break
				}
				time.Sleep(10 * time.Second)
			}
		}

		cancel()
	}
}

// Helper function to get the minimum of two durations
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func main() {
	conn, err := grpc.Dial("static.115.48.21.65.clients.your-server.de:50051",
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := datapb.NewTunnelDataClient(conn)
	forwardData(client)
}
