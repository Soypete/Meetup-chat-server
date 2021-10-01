package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	authToken := os.Getenv("DISCORD_SECRET")
	discord, err := discordgo.New("twitch-chat-sync " + authToken)
	if err != nil {
		log.Fatalln(err)
	}

	discord.AddHandler(sendResponse)

	// this intend is recieving messages
	// TODO: investigate intents
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// opens websocket connection
	err = discord.Open()
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening connection: %w", err))
		return
	}

	defer discord.Close()

}

func sendResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content != "" {
		s.ChannelMessageSend(m.ChannelID, "chat")
	}

}
