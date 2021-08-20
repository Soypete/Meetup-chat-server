package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendChat is called by the client to send a chat message to the server. The message is then
// stored in the database.
func (c *ChatServer) SendChat(ctx context.Context, msg *chat.ChatMessage) (*emptypb.Empty, error) {
	// TODO: add user to db
	err := c.database.InsertMessage(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert message to DB")
	}

	// send message to twitch over IRC connection
	err = c.twitchClient.SendChat(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send message to twitch")
	}
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
	return &chat.Chats{Messages: msgList}, nil
}
