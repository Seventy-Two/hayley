package dota

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/seventy-two/Cara/web"
)

// Service is configuration for the Dota service
type Service struct {
	DotaLeagueURL  string
	DotaListingURL string
	DotaMatchURL   string
	DotaHeroesURL  string
	APIKey         string
}

var serviceConfig *Service

const timer = "15:04:05"

type match struct {
	league        string
	viewers       int
	clock         string
	roshan        string
	radiant       string
	radiantScore  int
	radiantNet    int
	radiantHeroes []string
	radiantWins   int
	dire          string
	direScore     int
	direNet       int
	direHeroes    []string
	direWins      int
}

func dotamatches(params []string) ([]*match, error) {
	var matches []*match
	data := &LeagueGames{}
	listing := &LeagueListing{}
	getHeroes := &GetHeroes{}

	heroes := false

	if strings.Contains(strings.Join(params, ""), "h") {
		heroes = true
	}

	err := web.GetJSON(fmt.Sprintf(serviceConfig.DotaListingURL), listing)
	if err != nil {
		return nil, err
	}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.DotaLeagueURL, serviceConfig.APIKey), data)
	if err != nil {
		return nil, err
	}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.DotaHeroesURL, serviceConfig.APIKey), getHeroes)
	if err != nil {
		return nil, err
	}

	for _, game := range data.Result.Games {
		if (game.Spectators >= 1000) || (game.LeagueTier == 3 && game.Spectators >= 200) {
			m := &match{}
			m.radiant = game.RadiantTeam.TeamName
			m.dire = game.DireTeam.TeamName
			m.radiantScore = game.Scoreboard.Radiant.Score
			m.direScore = game.Scoreboard.Dire.Score
			m.radiantWins = game.RadiantSeriesWins
			m.direWins = game.DireSeriesWins
			m.viewers = game.Spectators
			if game.Scoreboard.RoshanRespawnTimer == 0 {
				m.roshan = "Up"
			} else {
				m.roshan = "Killed"
			}
			for k := 0; k < len(listing.Infos); k++ {
				if game.LeagueID == listing.Infos[k].LeagueID {
					m.league = listing.Infos[k].Name
				}
			}

			duration := int(game.Scoreboard.Duration)
			m.clock = fmt.Sprintf((time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Add(time.Duration(duration) * time.Second)).Format(timer))

			for j := 0; j < len(game.Scoreboard.Radiant.Players); j++ {
				m.radiantNet += game.Scoreboard.Radiant.Players[j].NetWorth
			}

			for j := 0; j < len(game.Scoreboard.Dire.Players); j++ {
				m.direNet += game.Scoreboard.Dire.Players[j].NetWorth
			}

			if heroes {
				for k := 0; k < len(game.Scoreboard.Radiant.Picks); k++ {
					m.radiantHeroes = append(m.radiantHeroes, strings.Join(getHerofromID(game.Scoreboard.Radiant.Picks[k].HeroID, getHeroes), " "))
				}
				for k := 0; k < len(game.Scoreboard.Dire.Picks); k++ {
					m.direHeroes = append(m.direHeroes, strings.Join(getHerofromID(game.Scoreboard.Dire.Picks[k].HeroID, getHeroes), " "))
				}

			}

			if m.dire == "" {
				m.dire = "Dire"
			}
			if m.radiant == "" {
				m.radiant = "Radiant"
			}
			matches = append(matches, m)
		}
	}

	if len(matches) == 0 {
		return nil, nil
	}

	sort.Slice(matches[:], func(i, j int) bool {
		return matches[i].viewers > matches[j].viewers
	})

	return matches, nil
}

func getHerofromID(id int, heroes *GetHeroes) (out []string) {
	if id == 0 {
		return []string{"PICK"}
	}
	out = getShortHero(id)
	if len(out) > 0 {
		return out
	}
	for m := 0; m < len(heroes.Result.Heroes); m++ {
		if heroes.Result.Heroes[m].ID == id {
			out = []string{heroes.Result.Heroes[m].LocalizedName}
			return out
		}
	}
	return nil
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

	switch matches[0] {
	case "!d2":
		res, err := dotamatches(matches[1:])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("an error occured (%s)", err))
			return
		}

		if res == nil {
			s.ChannelMessageSend(m.ChannelID, "No games")
			return
		}

		for _, topGame := range res {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title: topGame.league,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Series Score",
						Value:  fmt.Sprintf("%d - %d", topGame.radiantWins, topGame.direWins),
						Inline: true,
					},
					{
						Name:   "Time",
						Value:  topGame.clock,
						Inline: true,
					},
					{
						Name:   "Dota TV Viewers",
						Value:  fmt.Sprintf("%d", topGame.viewers),
						Inline: true,
					},
					{
						Name:   topGame.radiant,
						Value:  fmt.Sprintf("%d kills | %dg", topGame.radiantScore, topGame.radiantNet),
						Inline: true,
					},
					{
						Name:   topGame.dire,
						Value:  fmt.Sprintf("%d kills | %dg", topGame.direScore, topGame.direNet),
						Inline: true,
					},
					{
						Name:   "Roshan",
						Value:  topGame.roshan,
						Inline: true,
					},
				},
			})
		}
	}
}
