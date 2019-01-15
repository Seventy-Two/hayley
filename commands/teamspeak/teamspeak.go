package teamspeak

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/seventy-two/Hayley/nc"
)

type Service struct {
	Address  string
	Port     string
	Query    string
	Username string
	Password string
}

var serviceConfig *Service

var tsRegexp = regexp.MustCompile("client_nickname=[^\\s]*")

func ts() ([]string, error) {
	s, err := nc.RetrieveString(serviceConfig.Address, serviceConfig.Port, serviceConfig.Query, serviceConfig.Username, serviceConfig.Password)
	if err != nil {
		return nil, err
	}
	users := tsRegexp.FindAllString(s, -1)
	resp := []string{fmt.Sprintf("Currently in Teamspeak (%s):", serviceConfig.Address)}
	for _, user := range users {
		user = strings.Replace(user, "client_nickname=", "", -1)
		resp = append(resp, user)
	}
	return resp, nil
}

func RegisterService(dg *discordgo.Session, config *Service) {
	serviceConfig = config
	dg.AddHandler(invokeCommand)
}

func invokeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	matches := strings.Split(m.Content, " ")
	var str string

	switch matches[0] {
	case "!ts":
		res, err := ts()

		if err != nil {
			str = fmt.Sprintf("an error occured (%s)", err)
		} else {
			str = strings.Join(res, "\n")
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
