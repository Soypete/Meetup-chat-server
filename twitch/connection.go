package main

import (
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

func parseAuthCode(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL)
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := req.FormValue("code")
	fmt.Fprint(w, code)
}

func main() {

	wg := new(sync.WaitGroup)
	// mutex := new(sync.Mutex)
	http.HandleFunc("/oauth/redirect", parseAuthCode)
	go http.ListenAndServe("localhost:8081", nil)

	var tok *oauth2.Token
	// ctx := context.Background()
	conf := &oauth2.Config{
		ClientID: os.Getenv("TWITCH_ID"),
		// ClientSecret: "TWITCH_SECRET",
		Scopes:      []string{"chat:edit"},
		RedirectURL: "http://localhost:8081/oauth/redirect",
		Endpoint:    twitch.Endpoint}

	wg.Add(1)
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	go func() {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))

		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		// var code string
		// _, err = fmt.Scan(&code)
		// if err != nil {
		// log.Fatal(err)
		// }
		// fmt.Println(code)
		// mutex.Lock()
		// tok, err = conf.Exchange(ctx, code)
		// if err != nil {
		// log.Fatal(err)
		// }
		// _ = conf.Client(ctx, tok)
		// mutex.Unlock()
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
