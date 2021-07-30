package twitchirc

import (
	"fmt"
	"sync"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const peteTwitchChannel = "soypete01"

// IRC Connection to the twitch IRC server.
type IRC struct {
	client *v2.Client
	WG     *sync.WaitGroup
	mutex  *sync.Mutex
	tok    *oauth2.Token
	// tokChan chan *oath2.token
}

// SetupIRC gets the auth and connects to the twitch IRC server for channel.
func (irc *IRC) SetupIRC() error {
	c := v2.NewClient(peteTwitchChannel, "oauth:"+irc.tok.AccessToken)
	c.Join(peteTwitchChannel)
	c.OnConnect(func() { c.Say(peteTwitchChannel, "hello twitches") })
	c.OnPrivateMessage(func(msg v2.PrivateMessage) { fmt.Println(msg.Message) })
	err := c.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect over IRC")
	}
	irc.client = c
	return nil
}
