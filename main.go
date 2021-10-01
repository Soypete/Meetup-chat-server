package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/soypete/meetup-chat-server/postgres"
	server "github.com/soypete/meetup-chat-server/server"
	twitch "github.com/soypete/meetup-chat-server/twitch"
)

const grpcPort = "9090"
const httpPort = "8090"

func main() {
	// setup Database
	db := postgres.ConnectDB()
	err := db.Migrate()
	if err != nil {
		log.Fatalln(err)
	}

	// setup twitch IRC
	wg := new(sync.WaitGroup)
	irc, err := twitch.SetupTwitchIRC(db, wg)
	if err != nil {
		log.Fatalln(err)
	}

	// Configure gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chatServer := server.SetupGrpc(db, irc)

	fmt.Println("gRPCServer is configured and listening on port :9090")

	// TODO: clean shutdown - read about this
	err = chatServer.RunGrpc(ctx, grpcPort, wg)
	if err != nil {
		log.Fatalln(err)
	}

	// Configure Gateway Server
	chatServer.SetupGateway(ctx, httpPort, grpcPort)
	fmt.Println("GatewayServer is configured and running on port :8090")
	err = chatServer.GWServer.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
	wg.Wait()
}
