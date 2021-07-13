package main

import (
	"context"
	"fmt"

	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/grpc"
)

func main() {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := chat.NewGatewayConnectorClient(conn)
	ctx := context.Background()
	msg := chat.ChatMessage{
		UserName: "tester",
		Text:     "hello there",
	}
	_, err = client.SendChat(ctx, &msg)
	fmt.Println(err)
}
