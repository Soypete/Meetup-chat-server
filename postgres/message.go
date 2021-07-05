package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	chat "github.com/soypete/meetup-chat-server/protos"
)

const (
	portalMessage  = "portal" // the moderator portal
	twitchMessage  = "twitch"
	discordMessage = "discord"
)

func (pg *PG) InsertMessage(ctx context.Context, msg *chat.ChatMessage) error {
	query := `INSERT INTO chat_message (username, message_body, source)
			 values ($1, $2, $3)`

	// TODO: add switch for source
	result, err := pg.Client.Exec(query, msg.GetUserName(), msg.GetText(), portalMessage)
	if err != nil {
		fmt.Println(errors.Wrap(err, "cannot add message to the db"))
	}
	fmt.Println(result)
	return nil
}
