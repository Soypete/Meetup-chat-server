package twitchirc

import (
	"fmt"
	"sync"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/pkg/errors"
	"github.com/soypete/meetup-chat-server/postgres"
	chat "github.com/soypete/meetup-chat-server/protos"
	"golang.org/x/oauth2"
)

const peteTwitchChannel = "soypete01"

// TwitchIRC is used to enforce the methods to interact with twith.
type TwitchIRC interface {
	AppendChat(*chat.ChatMessage)
	PersistChat(v2.PrivateMessage)
}

// IRC Connection to the twitch IRC server.
type IRC struct {
	database postgres.PG
	client   *v2.Client
	wg       *sync.WaitGroup
	tok      *oauth2.Token
	msgQueue chan string
}

// SetupTwitchIRC sets up the IRC, configures oauth, and inits connection functions.
func SetupTwitchIRC(db postgres.PG, wg *sync.WaitGroup) (*IRC, error) {
	irc := &IRC{
		database: db,
		wg:       wg,
		msgQueue: make(chan string),
	}
	wg.Add(1)
	defer wg.Done()
	// TODO: fix go routine for clean shut down and
	// validate non-blocking calls.
	go func() error {
		// TODO error handling? this should shut down...
		err := irc.AuthTwitch()
		if err != nil {
			return fmt.Errorf("failed twitch auth: %w", err)
		}
		return nil
	}()
	return irc, nil
}

// connectIRC gets the auth and connects to the twitch IRC server for channel.
func (irc *IRC) connectIRC() error {
	c := v2.NewClient(peteTwitchChannel, "oauth:"+irc.tok.AccessToken)
	c.Join(peteTwitchChannel)
	c.OnConnect(func() { c.Say(peteTwitchChannel, "grpc twitch bot connected") })
	// TODO: define function that stores message to db
	c.OnPrivateMessage(func(msg v2.PrivateMessage) {
		irc.PersistChat(msg)
		select {
		case msg := <-irc.msgQueue:
			c.Say(peteTwitchChannel, msg)
		default:
		}
		// if len(messages) < 1 {
		// return
		// }
		// for _, msg := range messages {
		// c.Say(peteTwitchChannel, msg)
		// }
		// messages = []string{}
	})
	err := c.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect over IRC")
	}
	irc.client = c
	return nil
}
