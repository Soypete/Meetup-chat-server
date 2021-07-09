package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/soypete/meetup-chat-server/postgres"
	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/grpc"
)

// helper functions
func getText() string {
	rand.Seed(time.Now().Unix())
	answers := []string{
		"Follow soypete01 on twitch",
		"Check out my meetups",
		"Follow Me on Twitter",
		"Do you want to see more of my dogs?",
		"Say hi in chat",
		"Do you want me to work on the cloud technologies",
	}
	return answers[rand.Intn(len(answers))]
}

func setupTestDB(t *testing.T) (postgres.PG, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to setup testDB %s\n", err.Error())
	}

	// delay cleanup until after the test completes
	t.Cleanup(func() { db.Close() })

	// implement sqlx funtions
	dbx := sqlx.NewDb(db, "postgres")
	return postgres.PG{Client: dbx}, mock
}
func setupContextAndConnection(t *testing.T) (context.Context, *grpc.ClientConn) {
	// define timeout and canel func
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	// define grpc connections
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		t.Fatalf("failed to dial grpc conn")
	}
	t.Cleanup(func() { conn.Close() })

	return ctx, conn
}

// TestServer executes TableTestFunctions
func TestServer(t *testing.T) {
	// t.Run("test-grpc-client", testClientPublishGRPC)
	// post to server via http
	t.Run("test-http-client", testClientPublishHTTP)

}

func testClientPublishGRPC(t *testing.T) {
	// setup db mock
	pgClient, mock := setupTestDB(t)
	ctx, conn := setupContextAndConnection(t)

	// configure grpcServer
	chatServer := SetupGrpc(pgClient)
	err := chatServer.RunGrpc(ctx)
	if err != nil {
		t.Error(err)
	}

	// create testing client
	client := chat.NewGatewayConnectorClient(conn)

	// Mock db call
	sentence := getText()
	mock.ExpectExec("INSERT INTO chat_message").WithArgs("tester", sentence, "portal").WillReturnResult(sqlmock.NewResult(1, 1))

	// create grpc message
	msg := chat.ChatMessage{
		UserName: "tester",
		Text:     sentence,
	}
	// call function
	_, err = client.SendChat(ctx, &msg)
	if err != nil {
		t.Error(err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}

}

func testClientPublishHTTP(t *testing.T) {
	// setup db mock
	pgClient, mock := setupTestDB(t)

	ctx, _ := setupContextAndConnection(t)
	// configure grpcServer
	chatServer := SetupGrpc(pgClient)
	err := chatServer.RunGrpc(ctx)
	if err != nil {
		t.Error(err)
	}

	// Configure Gateway Server
	chatServer.SetupGateway(ctx)
	go chatServer.GWServer.ListenAndServe()
	fmt.Println("GatewayServer is configured and running on port :8090")

	// generate random sentence
	sentence := getText()
	// create grpc message
	msg := chat.ChatMessage{
		UserName: "tester",
		Text:     sentence,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec("INSERT INTO chat_message").WithArgs("tester", sentence, "portal").WillReturnResult(sqlmock.NewResult(1, 1))
	req, err := http.NewRequest("POST", "http://localhost:8090/chat/postmessage", bytes.NewBuffer(payload))
	if err != nil {
		t.Error(err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status code not expected %d\n", resp.StatusCode)

	}
	resp.Body.Close()
	chatServer.GWServer.Close()
}
