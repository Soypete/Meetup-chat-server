package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	chat "github.com/soypete/meetup-chat-server/protos"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ChatServer is the struct upon which the grpc methods are implemented.
type ChatServer struct {
	chat.UnimplementedGatewayConnectorServer
}

// NewServer created the grpc server for the chat messages.
func NewServer() *ChatServer {
	return &ChatServer{}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chat.RegisterGatewayConnectorServer(grpcServer, NewServer())

	go func() {
		log.Fatalln(grpcServer.Serve(lis))
	}()
	fmt.Println("grpc server listening")
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:9090",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = chat.RegisterGatewayConnectorHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://localhost:9090")
	log.Fatalln(gwServer.ListenAndServe())
}

func (c *ChatServer) SendChat(ctx context.Context, msg *chat.ChatMessage) (*emptypb.Empty, error) {
	// TODO: return recieved message
	fmt.Println(msg.GetText())
	return new(emptypb.Empty), nil
}
