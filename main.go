package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/soypete/meetup-chat-server/postgres"
	server "github.com/soypete/meetup-chat-server/server"
)

const grpcPort = "9090"
const httpPort = "8090"

func main() {
	db := postgres.ConnectDB()
	err := db.Migrate()
	if err != nil {
		log.Fatalln(err)
	}
	// Configure gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chatServer := server.SetupGrpc(db)

	fmt.Println("gRPCServer is configured and listening on port :9090")

	// TODO: clean shutdown - read about this
	err = chatServer.RunGrpc(ctx, grpcPort)
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
}
