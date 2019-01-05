package nfl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ryanuber/columnize"
	"github.com/seventy-two/Cara/web"
	"github.com/seventy-two/Hayley/service"
)

var serviceConfig *service.Service

func nfl() ([]string, error) {
	m := map[string]game{}
	err := web.GetJSON(serviceConfig.TargetURL, &m)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, g := range m {
		if g.Qtr == nil {
			continue
		}
		out = append(out, createGameString(&g))
	}
	return out, nil

}

func createGameString(g *game) string {
	home := getTeamName(g.Home.Abbr)
	away := getTeamName(g.Away.Abbr)
	homeScore := strconv.Itoa(g.Home.Score.T)
	awayScore := strconv.Itoa(g.Away.Score.T)

	if g.Posteam == g.Home.Abbr {
		homeScore = homeScore + " 🏈"
	} else {
		awayScore = "🏈 " + awayScore
	}
	down := ""
	switch g.Down {
	case 1:
		down = "1st Down"
	case 2:
		down = "2nd Down"
	case 3:
		down = "3rd Down"
	case 4:
		down = "4th Down"
	}

	ballAt := ""
	if g.Bp != 0 {
		ballAt = "| Ball at " + strconv.Itoa(g.Bp)
	}
	return awayScore + " - " + away + " @ " + home + " - " + homeScore + " | " + g.Clock + " " + *g.Qtr + ballAt + down + " | " + g.Media.Tv
}

func getTeamName(team string) string {
	switch team {
	case "ARI":
		return "Arizona Cardinals"
	case "ATL":
		return "Atlanta Falcons"
	case "CAR":
		return "Carolina Panthers"
	case "CHI":
		return "Chicago Bears"
	case "DAL":
		return "Dallas Cowboys"
	case "DET":
		return "Detroit Lions"
	case "GB":
		return "Green Bay Packers"
	case "MIN":
		return "Minnesota Vikings"
	case "NO":
		return "New Orleans Saints"
	case "NYG":
		return "New York Giants"
	case "PHI":
		return "Philadelphia Eagles"
	case "LAR":
		return "Los Angeles Rams"
	case "SF":
		return "San Fransisco 49ers"
	case "SEA":
		return "Seattle Seahawks"
	case "TB":
		return "Tampa Bay Buccaneers"
	case "WAS":
		return "Washington Redskins"
	case "BAL":
		return "Baltimore Ravens"
	case "BUF":
		return "Buffalo Bills"
	case "CIN":
		return "Cincinnati Bengals"
	case "CLE":
		return "Cleveland Browns"
	case "DEN":
		return "Denver Broncos"
	case "HOU":
		return "Houston Texans"
	case "IND":
		return "Indianapolis Colts"
	case "JAC":
		return "Jacksonville Jaguars"
	case "KC":
		return "Kansas City Chiefs"
	case "MIA":
		return "Miami Dolphins"
	case "NE":
		return "New England Patriots"
	case "NYJ":
		return "New York Jets"
	case "OAK":
		return "Oakland Raiders"
	case "PIT":
		return "Pittsburgh Steelers"
	case "LAC":
		return "Los Angeles Chargers"
	case "TEN":
		return "Tennessee Titans"
	default:
		return team
	}
}

func RegisterService(dg *discordgo.Session, config *service.Service) {
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
	case "!nfl":
		res, err := nfl()

		if err != nil {
			str = fmt.Sprintf("an error occured (%s)", err)
		} else {
			str = columnize.SimpleFormat(res)
		}

		if str != "" {
			fmtstr := fmt.Sprintf("```%s```", str)
			s.ChannelMessageSend(m.ChannelID, fmtstr)
		}
	}
}
