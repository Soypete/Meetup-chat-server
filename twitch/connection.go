package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

// IRC Connection to the twitch IRC server.
type IRC struct {
	client *v2.Client
}

func main() {
	var tok *oauth2.Token
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID: os.Getenv("TWITCH_ID"),
		// ClientSecret: "TWITCH_SECRET",
		Scopes:      []string{"chat:edit"},
		RedirectURL: "http://localhost",
		Endpoint:    twitch.Endpoint}

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	wg.Add(1)
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	go func() {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		var code string
		_, err := fmt.Scan(&code)
		if err != nil {
			log.Fatal(err)
		}
		mutex.Lock()
		tok, err = conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		_ = conf.Client(ctx, tok)
		mutex.Unlock()
		wg.Done()
	}()
	wg.Wait()
	fmt.Println(tok.TokenType, tok.AccessToken)
	c := v2.NewClient("soypete01", tok.AccessToken)
	c.Join("soypete01")
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	msgs, err := c.Userlist("soypete01")
	fmt.Println(msgs, err)
}
