package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	chat "github.com/soypete/meetup-chat-server/protos"
)

func cleanText(message string) string {
	message = strings.ReplaceAll(message, ";", "")
	message = strings.ReplaceAll(message, "DROP", "")
	message = strings.ReplaceAll(message, "SELECT", "")
	message = strings.ReplaceAll(message, "CREATE", "")
	message = strings.ReplaceAll(message, "ALTER", "")
	message = strings.ReplaceAll(message, "INSERT", "")
	message = strings.ReplaceAll(message, "INSERT", "")
	return message
}

// InsertMessage add a single message to the "chat_message" database table.
func (pg *PG) InsertMessage(ctx context.Context, msg *chat.ChatMessage) error {
	query := `INSERT INTO chat_message (username, message_body, source)
			 values ($1, $2, $3)`

	// TODO: add switch for source
	results, err := pg.Client.Exec(query, cleanText(msg.GetUserName()), cleanText(msg.GetText()), msg.GetSource().String())
	if err != nil {
		return errors.Wrap(err, "cannot add message to the db")
	}
	fmt.Println(results)
	return nil
}

// SelectMessages retrieves all the messages that have been stored in the database since
// the last message was recieved. The messages are pulled from the database based on the
// messageID. The message ID is of the postgres serial type and increments sequentially.
func (pg *PG) SelectMessages(lastMessageID int32) ([]*chat.ChatMessage, error) {
	var msgList []*chat.ChatMessage
	// TODO: add deleted at functionality
	// TODO: add banned functionality
	query := `SELECT user_name, message_body, source, created_at 
			  FROM chat_message
			  WHERE id > $1`
	rows, err := pg.Client.Queryx(query, lastMessageID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform query")
	}
	for rows.Next() {
		var msg chat.ChatMessage
		err = rows.Scan(&msg.UserName, &msg.Text, &msg.Source, &msg.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse messages query result")
		}
		msgList = append(msgList, &msg)
	}
	return msgList, nil
}
