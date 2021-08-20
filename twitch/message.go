package twitchirc

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/pkg/errors"
	chat "github.com/soypete/meetup-chat-server/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SendChat sends a message to the twitch chat through the IRC connection.
// This function currently uses the hard coded "soypete01" user and sents a
// text message.
func (irc *IRC) SendChat(msg *chat.ChatMessage) error {

	// TODO: add chat bot account and user name
	irc.client.Say(peteTwitchChannel, msg.GetText())
	return nil
}

func (irc *IRC) PersistChat(msg v2.PrivateMessage) {
	insertMessage := &chat.ChatMessage{
		UserName:  msg.User.Name,
		Text:      msg.Message,
		Timestamp: timestamppb.New(msg.Time),
		Source:    chat.Source_TWITCH,
	}
	// FIXME: Should we add context to the IRC
	err := irc.Database.InsertMessage(context.Background(), insertMessage)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to insert message in db"))
	}
	fmt.Println(msg.Message)
}
