package main

import (
	"log"
	"net"

	"github.com/Soypete/Meetup-chat-server/protos"

	"google.golang.org/grpc"
)

// ChatServer is the struct upon which the grpc methods are implemented.
type ChatServer struct{}

// NewServer created the grpc server for the chat messages.
func NewServer() ChatServer {
	return ChatServer{}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	protos.RegisterGatewatConnectorServer(grpcServer, NewServer())

	grpcServer.Serve(lis)
}

func (c *ChatServer) SendChat(msg protos.ChatMessage) error {
	return nil
}
