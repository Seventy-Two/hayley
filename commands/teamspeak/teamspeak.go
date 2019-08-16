package teamspeak

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/ryanuber/columnize"
	"github.com/seventy-two/hayley/nc"
)

var (
	clidRegexp     = regexp.MustCompile("clid=[^\\s]*")
	typeRegexp     = regexp.MustCompile("client_type=[0-9]")
	nickRegexp     = regexp.MustCompile("client_nickname=[^\\s]*")
	versRegexp     = regexp.MustCompile("client_version=[^\\s]*")
	platformRegexp = regexp.MustCompile("client_platform=[^\\s]*")
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

	resp := []string{fmt.Sprintf("Currently in Teamspeak (%s):", serviceConfig.Address)}
	users, err := buildTSUsers()
	if err != nil {
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
		return nil, err
	}
	var tsUsers []*tsUser
	users := clidRegexp.FindAllString(s, -1)
	types := typeRegexp.FindAllString(s, -1)
	var wg sync.WaitGroup

	for i, user := range users {
		u := &tsUser{
			cli:        user,
			clientType: strings.TrimPrefix(types[i], "client_type="),
		}
		wg.Add(1)
		tsUsers = append(tsUsers, u)
		go populateUser(u, &wg)
	}
	wg.Wait()

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

func populateUser(u *tsUser, wg *sync.WaitGroup) {
	if u.clientType == "1" {
		wg.Done()
		return
	}
	resp, err := nc.RetrieveString(serviceConfig.Address, serviceConfig.Port, "use %s\nlogin %s %s\n%s\nquit", "1", serviceConfig.Username, serviceConfig.Password, fmt.Sprintf("clientinfo %s", u.cli))
	if err != nil {
		wg.Done()
		return
	}
	u.nickname = strings.TrimPrefix(nickRegexp.FindString(resp), "client_nickname=")
	u.version = strings.TrimPrefix(versRegexp.FindString(resp), "client_version=")
	u.platform = strings.TrimPrefix(platformRegexp.FindString(resp), "client_platform=")
	u.inMuted, err = strconv.ParseBool(strings.TrimPrefix(inMutedRegexp.FindString(resp), "client_input_muted="))
	if err != nil {
		wg.Done()
		return
	}
	u.outMuted, err = strconv.ParseBool(strings.TrimPrefix(outMutedRegexp.FindString(resp), "client_output_muted="))
	if err != nil {
		wg.Done()
		return
	}
	wg.Done()
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
