package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func GetWebPage(url string) ([]byte, error) {
	client := http.DefaultClient
	// check http.Client initialized
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot create request %w", err)
	}
	request.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "value",
	})
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, fmt.Errorf("failure in Do request:\n %w ---\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New("not a 200 status code")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Cannot parse response body: %w", err)
	}
	defer resp.Body.Close()
	return body, nil
}
func parseAuthCode(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := req.FormValue("code")
	fmt.Fprint(os.Stdout, code)
}

func main() {

	wg := new(sync.WaitGroup)
	http.HandleFunc("/oauth/redirect", parseAuthCode)
	go http.ListenAndServe("localhost:8081", nil)

	var tok *oauth2.Token
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("TWITCH_ID"),
		ClientSecret: os.Getenv("TWITCH_SECRET"),
		Scopes:       []string{"chat:read", "chat:edit"},
		RedirectURL:  "http://localhost:8081/oauth/redirect",
		Endpoint:     twitch.Endpoint}

	wg.Add(1)
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	go func() {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		var code string
		_, err := fmt.Scan(&code)
		if err != nil {
			log.Fatal(err)
		}
		tok, err = conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		_ = conf.Client(ctx, tok)
		wg.Done()
	}()
	wg.Wait()
	c := v2.NewClient("soypete01", "oauth:"+tok.AccessToken)
	c.Join("soypete01")
	c.OnConnect(func() { c.Say("soypete01", "hello twitches") })
	c.OnPrivateMessage(func(msg v2.PrivateMessage) { fmt.Println(msg.Message) })
	err := c.Connect()
	if err != nil {
		panic(err)
	}

}
