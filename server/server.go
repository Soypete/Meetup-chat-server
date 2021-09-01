package server

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/soypete/meetup-chat-server/postgres"
	chat "github.com/soypete/meetup-chat-server/protos"
	twitch "github.com/soypete/meetup-chat-server/twitch"
	"google.golang.org/grpc"
)

// ChatServer is the struct upon which the grpc methods are implemented.
type ChatServer struct {
	chat.UnimplementedGatewayConnectorServer
	GWServer     *http.Server
	database     postgres.PG
	twitchClient twitch.TwitchIRC
}

// SetupGrpc created the grpc server for the chat messages.
func SetupGrpc(db postgres.PG, tc twitch.TwitchIRC) *ChatServer {
	cs := ChatServer{
		database:     db,
		twitchClient: tc,
	}
	return &cs
}

// SetupGateway creates the Rest server via grpc connection/
func (cs *ChatServer) SetupGateway(ctx context.Context, port string, grpcPort string) error {
	conn, err := grpc.DialContext(
		ctx,
		"localhost:"+grpcPort,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to setup grpc client connection: %w")
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	if err = chat.RegisterGatewayConnectorHandler(ctx, gwmux, conn); err != nil {
		return errors.Wrap(err, "failed to register gateway: %w")
	}

	cs.GWServer = &http.Server{
		Addr:    "localhost:" + port,
		Handler: gwmux,
	}

	return nil
}

// RunGrpc administer the server used to handle chat messages.
func (cs *ChatServer) RunGrpc(ctx context.Context, port string, wg *sync.WaitGroup) error {
	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		return errors.Wrap(err, "cannot setup tcp connection: %w")
	}

	grpcServer := grpc.NewServer()
	chat.RegisterGatewayConnectorServer(grpcServer, cs)

	wg.Add(1)
	go func() error {
		defer wg.Done()
		err = grpcServer.Serve(lis)
		if err != nil {
			return errors.Wrap(err, "grpc server failure: %w")
		}
		return nil
	}()
	if err != nil {
		return err
	}
	return nil
}
