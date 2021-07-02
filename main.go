package main

import (
	"context"
	"fmt"
	"log"

	server "github.com/soypete/meetup-chat-server/server"
)

func main() {
	ctx := context.Background()
	chatServer := server.Setup(ctx)

	fmt.Println("server is configured")
	// TODO: clean shutdown with channel listener

	err := chatServer.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}

}
