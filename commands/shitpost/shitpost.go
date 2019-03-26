package shitpost

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// RegisterService will reg shitpost
func RegisterService(dg *discordgo.Session) {
	dg.AddHandler(invokeCommand)
}

func invokeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if hasWord(m.Content, "same") {
		s.ChannelMessageSend(m.ChannelID, "same")
	}

	if hasWord(m.Content, "brexit") {
		s.ChannelMessageSend(m.ChannelID, brexitCountdown())
	}

	if hasWord(m.Content, "flac") {
		s.ChannelMessageSend(m.ChannelID, flac())
	}

	if hasWord(m.Content, "linux") && !strings.Contains(strings.ToLower(m.Content), "gnu") {
		s.ChannelMessageSend(m.ChannelID, rms())
	}
}

func hasWord(s, match string) bool {
	fields := strings.Fields(s)
	for _, field := range fields {
		if strings.ToLower(field) == match {
			return true
		}
	}
	return false
}
