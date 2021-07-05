package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/soypete/meetup-chat-server/postgres"
	server "github.com/soypete/meetup-chat-server/server"
)

func main() {
	db := postgres.ConnectDB()
	err := db.Migrate()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chatServer := server.SetupGrpc(db)

	fmt.Println("server is configured", chatServer)

	// TODO: clean shutdown - read about this
	err = chatServer.RunGrpc(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	chatServer.SetupGateway(ctx)

	err = chatServer.GWServer.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
