package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/soypete/meetup-chat-server/postgres"
	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ChatServer is the struct upon which the grpc methods are implemented.
type ChatServer struct {
	chat.UnimplementedGatewayConnectorServer
	GWServer *http.Server
	database postgres.PG
}

// SetupGrpc created the grpc server for the chat messages.
func SetupGrpc(db postgres.PG) *ChatServer {
	cs := ChatServer{
		database: db,
	}
	return &cs
}

// SetupGateway creates the Rest server via grpc connection/
func (cs *ChatServer) SetupGateway(ctx context.Context) error {
	conn, err := grpc.DialContext(
		ctx,
		"localhost:9090",
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
		Addr:    ":8090",
		Handler: gwmux,
	}

	return nil
}

// RunGrpc administer the server used to handle chat messages.
func (cs *ChatServer) RunGrpc(ctx context.Context) error {
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return errors.Wrap(err, "cannot setup tcp connection: %w")
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
	return nil
}

// SendChat is called by the client to send a chat message to the server. The message is then
// stored in the database.
func (c *ChatServer) SendChat(ctx context.Context, msg *chat.ChatMessage) (*emptypb.Empty, error) {
	fmt.Println(msg.GetText())
	err := c.database.InsertMessage(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert message to DB")
	}
	// TODO: add user to db
	return new(emptypb.Empty), nil
}

// GetChat is used by a client to retrieve messages that the sever has collected. The client supplies the last
// messageID that the recieved and the server returns all the messages send after that last ID. It's are sequential
// so it just has to return when the messageID is larger than the last MessageID.
func (c *ChatServer) GetChat(ctx context.Context, request *chat.RetrieveChatMessages) (*chat.Chats, error) {
	fmt.Println(request)
	msgList, err := c.database.SelectMessages(request.GetLastMessageId())
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve caht messages")
	}
	return &chat.Chats{messages: msgList}, nil
}
