package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	datapb "Pulse/gen/protos/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// forwardData connects to the gRPC server and handles data forwarding.
func forwardData(client datapb.TunnelDataClient) {
	for {
		stream, err := client.ForwardData(context.Background())
		if err != nil {
			log.Printf("Failed to create stream: %v. Retrying...", err)
			time.Sleep(time.Second * 2)
			continue
		}

		go func() {
			for {
				packet, err := stream.Recv()
				if err == io.EOF {
					log.Println("Stream closed by server")
					return
				}
				if err != nil {
					log.Printf("Error receiving data: %v", err)
					return
				}
				log.Printf("Received packet: %s", string(packet.Data))
				url := fmt.Sprintf("http://localhost:3000/%s", string(packet.Data))
				resp, err := http.Get(url)
				if err != nil {
					log.Printf("Error forwarding request: %v", err)
					continue
				}
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)

				err = stream.Send(&datapb.DataPacket{TunnelId: packet.TunnelId, Data: body})
				if err != nil {
					log.Printf("Error sending response: %v", err)
				}
			}
		}()

		time.Sleep(5 * time.Second)
	}
}

func main() {
	conn, _ := grpc.Dial("static.115.48.21.65.clients.your-server.de:50051", grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time: 30 * time.Second, Timeout: 10 * time.Second, PermitWithoutStream: false,
	}))
	defer conn.Close()

	client := datapb.NewTunnelDataClient(conn)
	forwardData(client)
}
