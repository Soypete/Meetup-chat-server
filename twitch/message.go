package twitchirc

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/gempir/go-twitch-irc/v2"
	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TODO: use channels cause they fit the use case better
var messages []string

// AppendChat sends a message to the twitch chat through the IRC connection.
// This function currently uses the hard coded "soypete01" user and sents a
// text message.
func (irc *IRC) AppendChat(msg *chat.ChatMessage) {
	messages = append(messages, msg.GetText())
}

// PresistChat is used to handle PrivateMessages received from twitch IRC.
// When private message is received it is transformed to standard proto definition,
// and then stored in the DB. The v2.OnPrivateMessage(func) handler that calls this
// method does not have error handling, so error is just logged.
func (irc *IRC) PersistChat(msg v2.PrivateMessage) {
	insertMessage := &chat.ChatMessage{
		UserName:  msg.User.Name,
		Text:      msg.Message,
		Timestamp: timestamppb.New(msg.Time),
		Source:    chat.Source_TWITCH,
	}
	// FIXME: Should we add context to the IRC
	err := irc.database.InsertMessage(context.Background(), insertMessage)
	if err != nil {
		log.Println(fmt.Errorf("failed to insert message in db: %w", err))
	}
}
