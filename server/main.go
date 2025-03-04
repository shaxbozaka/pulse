package main

import (
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	pb "Pulse/proto"
	"google.golang.org/grpc"
)

type TunnelServer struct {
	pb.UnimplementedTunnelServiceServer
}

func (s *TunnelServer) CreateTunnel(req *pb.TunnelRequest, stream pb.TunnelService_CreateTunnelServer) error {
	log.Printf("Tunnel created for %s", req.Subdomain)

	for i := 0; i < 10; i++ {
		// Context check to handle client disconnects or timeouts
		select {
		case <-stream.Context().Done():
			log.Println("Client disconnected")
			return stream.Context().Err()
		default:
			// Continue processing normally.
		}

		// Construct streaming response
		resp := &pb.TrafficResponse{
			Payload: []byte(fmt.Sprintf("Tunnel active - Sequence: %d", i)),
		}

		// Send response over the stream
		if err := stream.Send(resp); err != nil {
			log.Printf("Error streaming response (Sequence %d): %v", i, err)
			return fmt.Errorf("stream.Send failed for sequence %d: %w", i, err)
		}

		// Simulate response delay
		time.Sleep(1 * time.Second)
	}

	log.Println("Tunnel stream closed successfully")
	return nil
}
func (s *TunnelServer) StreamTraffic(stream pb.TunnelService_StreamTrafficServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Println("Stream closed")
			return err
		}
		log.Println("Forwarding traffic...")

		// ✅ Correct type: Using `TrafficResponse`
		resp := &pb.TrafficResponse{Payload: req.Payload}
		if err := stream.Send(resp); err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTunnelServiceServer(grpcServer, &TunnelServer{})

	// ✅ Enable reflection
	reflection.Register(grpcServer)

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
