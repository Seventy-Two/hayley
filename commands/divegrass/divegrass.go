package divegrass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ryanuber/columnize"
	"github.com/seventy-two/Hayley/service"
)

var serviceConfig *service.Service

func divegrass() ([]string, error) {
	data := &Fixtures{}
	client := &http.Client{}
	frames := []string{"p1", "n2"}
	var str []string
	for _, frame := range frames { // Loop through the leagues we want
		url := fmt.Sprintf(serviceConfig.TargetURL, url.QueryEscape(frame))
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("X-Auth-Token", "2bb6957d70554e388044373fb073671d")
		req.Header.Add("X-Response-Control", `minified`)
		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = json.Unmarshal(body, data)
		if err != nil {
			str = append(str, fmt.Sprintf("There was a problem with your request"))
			return str, nil
		}
		if data.Count <= 0 && len(str) == 0 {
			continue
		}
		for j := 0; j < data.Count; j++ {
			date := data.Fixtures[j].Date
			//			if date.Date() != time.Now().Date() {
			loc, err := time.LoadLocation("Europe/London")
			if err != nil {
				return nil, err
			}
			fmtdate := fmt.Sprintf("%s", date.In(loc).Format("Mon 2 Jan 15:04"))
			//			} else {
			//				hour,min,_ := date.Clock()
			//				fmtdate := fmt.Sprintf("%d:%d", hour ,min)
			//			}
			hName := data.Fixtures[j].HomeTeamName
			aName := data.Fixtures[j].AwayTeamName
			hScore := 0
			aScore := 0
			if data.Fixtures[j].Result.GoalsHomeTeam != nil {
				hScore = *data.Fixtures[j].Result.GoalsHomeTeam
			}
			if data.Fixtures[j].Result.GoalsAwayTeam != nil {
				aScore = *data.Fixtures[j].Result.GoalsAwayTeam
			}
			if data.Fixtures[j].Status == "IN_PLAY" {
				fmtdate = "Live"
			}
			if data.Fixtures[j].Status == "FINISHED" {
				fmtdate = "Final"
			}
			str = append(str, fmt.Sprintf("%s |%s | %d - %d | %s", fmtdate, hName, hScore, aScore, aName))
		}
	}
	if len(str) == 0 {
		return nil, nil
	}

	return str, nil
}

func RegisterService(dg *discordgo.Session, config *service.Service) {
	serviceConfig = config
	dg.AddHandler(invokeCommand)
}

func invokeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	var str string

	matches := strings.Split(m.Content, " ")

	switch matches[0] {
	case "!foot":
		res, err := divegrass()
		if err == nil {
			str = columnize.SimpleFormat(res)
		} else {
			str = fmt.Sprintf("an error occured (%s)", err)
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
