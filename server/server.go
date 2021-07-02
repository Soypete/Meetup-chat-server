package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ChatServer is the struct upon which the grpc methods are implemented.
type ChatServer struct {
	chat.UnimplementedGatewayConnectorServer
	gwServer *http.Server
}

// Setup created the grpc server for the chat messages.
func Setup(ctx context.Context) *ChatServer {
	cs := ChatServer{}

	conn, err := grpc.DialContext(
		ctx,
		"localhost:9090",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to setup grpc client connection: %w"))
		return nil
	}
	fmt.Println("setup grpc connection")
	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = chat.RegisterGatewayConnectorHandler(ctx, gwmux, conn)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to register gateway: %w"))
		return nil
	}

	fmt.Println("setup gwmux server connection")
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
	cs.gwServer = gwServer
	return &cs
}

// Run administer the server used to handle chat messages.
func (cs *ChatServer) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		fmt.Println(errors.Wrap(err, "cannot setup tcp connection: %w"))
		return nil
	}
	grpcServer := grpc.NewServer()
	chat.RegisterGatewayConnectorServer(grpcServer, cs)

	go func() error {
		err := grpcServer.Serve(lis)
		if err != nil {
			return errors.Wrap(err, "grpc server failure: %w")
		}
		return nil
	}()
	fmt.Println("grpc server listening")

	fmt.Println("Serving gRPC-Gateway on http://localhost:8090")
	err = cs.gwServer.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "server died: %w")
	}
	return nil
}

func (c *ChatServer) SendChat(ctx context.Context, msg *chat.ChatMessage) (*emptypb.Empty, error) {
	// TODO: return recieved message
	fmt.Println(msg.GetText())
	return new(emptypb.Empty), nil
}
