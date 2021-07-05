package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/soypete/meetup-chat-server/postgres"
	chat "github.com/soypete/meetup-chat-server/protos"
	"github.com/soypete/meetup-chat-server/server"
	"google.golang.org/grpc"
)

func setupTests(t *testing.T) (*grpc.ClientConn, postgres.PG, context.Context, error) {

	rand.Seed(time.Now().Unix())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	db := postgres.ConnectDB()
	err := db.Migrate()
	if err != nil {
		return nil, postgres.PG{}, ctx, errors.Wrap(err, "test setup failed to config db")
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		return nil, postgres.PG{}, ctx, errors.Wrap(err, "test setup failed to setup grpcClient")
	}

	t.Cleanup(func() { conn.Close() })
	chatServer := server.SetupGrpc(db)
	err = chatServer.RunGrpc(ctx)
	if err != nil {
		return nil, postgres.PG{}, ctx, errors.Wrap(err, "test setup failed to setup grpcserver")
	}
	// Configure Gateway Server
	chatServer.SetupGateway(ctx)
	err = chatServer.GWServer.ListenAndServe()
	if err != nil {
		return nil, postgres.PG{}, ctx, errors.Wrap(err, "test setup failed to setup gateway")
	}
	t.Cleanup(func() { chatServer.GWServer.Close() })
	return conn, db, ctx, nil
}

func testSetup(t *testing.T) {
	_, db, _, err := setupTests(t)
	if err != nil {
		t.Error(err)
	}
	err = db.Client.Ping()
	if err != nil {
		t.Error(err)
	}
}

func TestServer(t *testing.T) {
	t.Run("test configs", testSetup)
	// post to server via grpc
	t.Run("test succesful publish to grpc client", testClientPublishGRPC)
	// post to server via http
	t.Run("test succesful publish to http client", testClientPublishHTTP)

}

func getText() string {
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

func testClientPublishGRPC(t *testing.T) {
	conn, db, ctx, err := setupTests(t)
	if err != nil {
		t.Error(err)
	}
	client := chat.NewGatewayConnectorClient(conn)
	for i := 0; i <= rand.Intn(20); i++ {
		// generate random sentences
		sentence := getText()
		// make grpc call
		msg := chat.ChatMessage{
			UserName: "tester",
			Text:     sentence,
		}
		// response body is empty

		resp, err := client.SendChat(ctx, &msg)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(resp)
	}
	// check db
	fmt.Println(db.Client.Query(`SELECT COUNT(ID) FROM chat_message`))
}

func testClientPublishHTTP(t *testing.T) {
	_, db, _, err := setupTests(t)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i <= rand.Intn(20); i++ {
		// generate random sentences
		sentence := getText()
		msg := fmt.Sprintf(`{"user_name": "tester","text":"%s"}`, sentence)
		payload, err := json.Marshal(msg)
		if err != nil {
			t.Error(err)
			return
		}

		req := httptest.NewRequest("POST", "localhost:8090/chat/postmessage", bytes.NewBuffer(payload))
		httpClient := http.DefaultClient
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != 200 {
			t.FailNow()
		}
		resp.Body.Close()
	}
	fmt.Println(db.Client.Query(`SELECT COUNT(ID) FROM chat_message`))
}
