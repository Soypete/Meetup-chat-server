package twitchirc

import chat "github.com/soypete/meetup-chat-server/protos"

func (irc *IRC) SendChat(msg *chat.ChatMessage) {
	// TODO: add chat bot account and user name
	irc.client.Say(peteTwitchChannel, msg.GetText())
}
