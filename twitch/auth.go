package twitchirc

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

func parseAuthCode(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := req.FormValue("code")
	fmt.Fprint(os.Stdout, code)
}

// AuthTwitch use oauth2 protocol to retrieve oauth2 token for twitch IRC.
// _NOTE_: this has not been tested on long standing projects.
func (irc *IRC) AuthTwitch() error {
	http.HandleFunc("/oauth/redirect", parseAuthCode)
	go http.ListenAndServe("localhost:8081", nil)

	ctx := context.Background()
	conf := &oauth2.Config{
		// TODO: use const for the following.
		ClientID:     os.Getenv("TWITCH_ID"),
		ClientSecret: os.Getenv("TWITCH_SECRET"),
		Scopes:       []string{"chat:read", "chat:edit", "channel:moderate"},
		RedirectURL:  "http://localhost:8081/oauth/redirect",
		Endpoint:     twitch.Endpoint,
	}
	// TODO: fix error handling with channels?
	irc.wg.Add(1)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	go func() {
		defer irc.wg.Done()
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		var code string
		_, err := fmt.Scan(&code)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot get input from standard in: %w", err))
		}

		irc.tok, err = conf.Exchange(ctx, code)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to get token with auth code: %w", err))
		}
		_ = conf.Client(ctx, irc.tok)
		// is this overkill?
		fmt.Println("start twitch connection")
		err = irc.connectIRC()
		if err != nil {
			fmt.Println(fmt.Errorf("failed to connect over IRC: %w", err))
		}
		fmt.Println("connection completed")
	}()
	return nil
}
