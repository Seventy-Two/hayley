package teamspeak

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/common/log"
	"github.com/ryanuber/columnize"
	"github.com/seventy-two/hayley/nc"
)

var (
	clidRegexp     = regexp.MustCompile("clid=[^\\s]*")
	typeRegexp     = regexp.MustCompile("client_type=[0-9]")
	nickRegexp     = regexp.MustCompile("client_nickname=[^\\s]*")
	versRegexp     = regexp.MustCompile("client_version=[^\\s]*")
	platformRegexp = regexp.MustCompile("client_platform=[^\\|]*")
	inMutedRegexp  = regexp.MustCompile("client_input_muted=[0-1]")
	outMutedRegexp = regexp.MustCompile("client_output_muted=[0-1]")
)

type tsUser struct {
	cli        string
	clientType string
	nickname   string
	version    string
	platform   string
	inMuted    bool
	outMuted   bool
}

type Service struct {
	Address  string
	Port     string
	Query    string
	Username string
	Password string
}

var serviceConfig *Service

func ts() ([]string, error) {

	resp := []string{fmt.Sprintf("Currently in Teamspeak:")}
	users, err := buildTSUsers()
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	tsString := columnize.SimpleFormat(users)
	resp = append(resp, tsString)
	return resp, nil
}

func buildTSUsers() ([]string, error) {
	s, err := nc.RetrieveString(serviceConfig.Address, serviceConfig.Port, serviceConfig.Query, serviceConfig.Username, serviceConfig.Password)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	var tsUsers []*tsUser
	types := typeRegexp.FindAllString(s, -1)
	nicks := nickRegexp.FindAllString(s, -1)
	vers := versRegexp.FindAllString(s, -1)
	platforms := platformRegexp.FindAllString(s, -1)
	muteds := inMutedRegexp.FindAllString(s, -1)
	deafs := outMutedRegexp.FindAllString(s, -1)
	var mBool []bool
	var dBool []bool

	for _, m := range muteds {
		muted, err := strconv.ParseBool(strings.TrimPrefix(m, "client_input_muted="))
		if err != nil {
			return nil, err
		}
		mBool = append(mBool, muted)
	}

	for _, d := range deafs {
		deaf, err := strconv.ParseBool(strings.TrimPrefix(d, "client_output_muted="))
		if err != nil {
			return nil, err
		}
		dBool = append(dBool, deaf)
	}

	for i := range types {
		u := &tsUser{
			clientType: strings.TrimPrefix(types[i], "client_type="),
			nickname:   strings.TrimPrefix(nicks[i], "client_nickname="),
			version:    strings.TrimPrefix(vers[i], "client_version="),
			platform:   strings.TrimPrefix(platforms[i], "client_platform="),
			inMuted:    mBool[i],
			outMuted:   dBool[i],
		}
		tsUsers = append(tsUsers, u)
	}

	var userList []string
	for _, user := range tsUsers {
		if user.clientType == "1" {
			continue
		}
		mute := ""
		deaf := ""
		if user.inMuted {
			mute = "Mute"
		}
		if user.outMuted {
			deaf = "Deaf"
		}
		userString := strings.Replace(fmt.Sprintf("%s | %s | %s | %s", user.nickname, user.platform, mute, deaf), "\\s", " ", -1)
		userList = append(userList, userString)
	}
	return userList, nil

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
		} else if res == nil {
			s.ChannelMessageSend(m.ChannelID, "Nobody in Teamspeak (0x48.io)")
			s.ChannelMessageSend(m.ChannelID, "sad")
			return
		} else {
			str = strings.Join(res, "\n")
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
