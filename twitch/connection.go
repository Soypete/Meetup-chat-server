package twitchirc

import (
	"fmt"
	"sync"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/pkg/errors"
	"github.com/soypete/meetup-chat-server/postgres"
	"golang.org/x/oauth2"
)

const peteTwitchChannel = "soypete01"

// IRC Connection to the twitch IRC server.
type IRC struct {
	Database postgres.PG
	client   *v2.Client
	WG       *sync.WaitGroup
	tok      *oauth2.Token
}

// SetupIRC gets the auth and connects to the twitch IRC server for channel.
func (irc *IRC) SetupIRC() error {
	c := v2.NewClient(peteTwitchChannel, "oauth:"+irc.tok.AccessToken)
	c.Join(peteTwitchChannel)
	c.OnConnect(func() { c.Say(peteTwitchChannel, "grpc twitch bot connected") })
	// TODO: define function that stores message to db
	c.OnPrivateMessage(func(msg v2.PrivateMessage) { fmt.Println(msg.Message) })
	err := c.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect over IRC")
	}
	irc.client = c
	return nil
}
