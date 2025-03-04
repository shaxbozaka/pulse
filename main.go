package main

import (
	"log"
	"net"

	pb "github.com/shaxbozaka/pulse/github.com/shaxbozaka/pulse/proto"
	"google.golang.org/grpc"
)

type TunnelServer struct {
	pb.UnimplementedTunnelServiceServer
}

func (s *TunnelServer) CreateTunnel(req *pb.TunnelRequest, stream pb.TunnelService_StreamTrafficServer) error {
	log.Printf("Tunnel created for %s", req.Subdomain)
	// Store tunnel info and wait for traffic
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
		// Process and send response
		resp := &pb.TrafficResponse{Payload: req.Payload}
		stream.Send(resp)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTunnelServiceServer(grpcServer, &TunnelServer{})
	log.Println("gRPC server running on port 50051")
	grpcServer.Serve(listener)
}
