package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
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

func TestClientPublishGRPC(t *testing.T) {
	// setup db mock
	pgClient, mock := setupTestDB(t)
	ctx, conn := setupContextAndConnection(t)

	// configure grpcServer
	chatServer := SetupGrpc(pgClient)
	err := chatServer.RunGrpc(ctx, "9090")
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

func TestClientPublishHTTP(t *testing.T) {
	// setup db mock
	pgClient, mock := setupTestDB(t)

	ctx, _ := setupContextAndConnection(t)

	// configure grpcServer
	chatServer := SetupGrpc(pgClient)
	err := chatServer.RunGrpc(ctx, "9091")
	if err != nil {
		// this is setup step
		t.Fatalf("cannot setup grpc %v", err)
	}

	// Configure Gateway Server
	chatServer.SetupGateway(ctx, "8091", "9091")
	go chatServer.GWServer.ListenAndServe()
	fmt.Println("GatewayServer is configured and running on port :8091")

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
	req, err := http.NewRequest("POST", "http://localhost:8091/chat/postmessage", bytes.NewBuffer(payload))
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

func TestIntegration(t *testing.T) {
	// only run in integration env
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("Skipping integration tests")
	}
	t.Run("test-grpc-client", testMessagePublishGRPCEndToEnd)
	// post to server via http
	t.Run("test-http-client", testMessagePublishHTTPEndToEnd)

}
func testMessagePublishGRPCEndToEnd(t *testing.T) {
	db := postgres.ConnectDB()
	ctx, conn := setupContextAndConnection(t)

	// configure grpcServer
	chatServer := SetupGrpc(db)
	err := chatServer.RunGrpc(ctx, "9090")
	if err != nil {
		t.Error(err)
	}

	// create testing client
	client := chat.NewGatewayConnectorClient(conn)
	rownum := rand.Intn(20)
	for i := 0; i < rownum; i++ {
		sentence := getText()

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
	}
	row := db.Client.QueryRow(`SELECT COUNT(ID) FROM chat_message`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Error(err)
	}
	if count != rownum {
		t.Errorf("db out of sync: expected row count: %d, actuial row count: %d", rownum, count)
	}

}

func testMessagePublishHTTPEndToEnd(t *testing.T) {
	db := postgres.ConnectDB()

	ctx, _ := setupContextAndConnection(t)
	// configure grpcServer

	// wait to unbind address
	chatServer := SetupGrpc(db)
	err := chatServer.RunGrpc(ctx, "9091")
	if err != nil {
		// this is setup step
		t.Fatalf("cannot setup grpc %v", err)
	}

	// Configure Gateway Server
	chatServer.SetupGateway(ctx, "8091", "9091")
	go chatServer.GWServer.ListenAndServe()
	fmt.Println("GatewayServer is configured and running on port :8090")

	rownum := rand.Intn(20)
	for i := 0; i < rownum; i++ {
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

		req, err := http.NewRequest("POST", "http://localhost:8091/chat/postmessage", bytes.NewBuffer(payload))
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
	}
	row := db.Client.QueryRow(`SELECT COUNT(ID) FROM chat_message`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Error(err)
	}
	if count != rownum {
		t.Errorf("db out of sync: expected row count: %d, actuial row count: %d", rownum, count)
	}
	chatServer.GWServer.Close()
}
