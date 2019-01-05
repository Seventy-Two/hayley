package dota

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"golang.org/x/text/width"

	"github.com/bwmarrin/discordgo"
	"github.com/ryanuber/columnize"

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

func dotamatches(matches []string) (msg []string, err error) {
	data := &LeagueGames{}
	listing := &LeagueListing{}
	getHeroes := &GetHeroes{}
	var radiantNet int
	var direNet int
	var worth int
	heroes := false
	showScore := true
	showTowers := false
	if strings.Contains(strings.Join(matches, ""), "h") {
		heroes = true
	}
	//if strings.Contains(matches[0], "s") {
	//	showScore = true
	//}
	//if strings.Contains(matches[0], "t") {
	//	showTowers = true
	//}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.DotaListingURL, serviceConfig.APIKey), listing)
	if err != nil {
		//		msg = append(msg, fmt.Sprintf("Could not retrieve league listings.",))
		//		return msg, nil
		log.Print(err)
	}
	err = web.GetJSON(fmt.Sprintf(serviceConfig.DotaLeagueURL, serviceConfig.APIKey), data)
	if err != nil {
		log.Print(err)
		msg = append(msg, fmt.Sprintf("Could not retrieve matches."))
		return msg, nil
	}
	web.GetJSON(fmt.Sprintf(serviceConfig.DotaHeroesURL, serviceConfig.APIKey), getHeroes)
	if err != nil {
		msg = append(msg, fmt.Sprintf("Could not retrieve heroes."))
		return msg, nil
	}
	var str []string
	for i := 0; i < len(data.Result.Games); i++ {
		worth = 0
		radiantNet = 0
		direNet = 0
		if (data.Result.Games[i].Spectators >= 1000) || (data.Result.Games[i].LeagueTier == 3 && data.Result.Games[i].Spectators >= 200) {
			var leaguename string
			herostr := []string{""}
			radTower := ""
			direTower := ""
			rad1 := "|"
			rad2 := "|"
			dire1 := "|"
			dire2 := "|"
			rads := []string{""}
			dires := []string{""}
			radiant := data.Result.Games[i].RadiantTeam.TeamName
			dire := data.Result.Games[i].DireTeam.TeamName
			radiantScore := data.Result.Games[i].Scoreboard.Radiant.Score
			direScore := data.Result.Games[i].Scoreboard.Dire.Score
			game := data.Result.Games[i].RadiantSeriesWins + data.Result.Games[i].DireSeriesWins + 1
			viewers := data.Result.Games[i].Spectators
			for k := 0; k < len(listing.Result.Leagues); k++ {
				if data.Result.Games[i].LeagueID == listing.Result.Leagues[k].Leagueid {
					leaguename = listing.Result.Leagues[k].Name
					leaguename = strings.Replace(leaguename, "#DOTA_Item", "", -1)
					leaguename = strings.Replace(leaguename, "_", " ", -1)
					leaguename = strings.TrimSpace(leaguename)
				}
			}

			duration := int(data.Result.Games[i].Scoreboard.Duration)
			t := fmt.Sprintf((time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Add(time.Duration(duration) * time.Second)).Format(timer))

			for j := 0; j < len(data.Result.Games[i].Scoreboard.Radiant.Players); j++ {
				radiantNet += data.Result.Games[i].Scoreboard.Radiant.Players[j].NetWorth
			}
			if heroes {
				for k := 0; k < len(data.Result.Games[i].Scoreboard.Radiant.Picks); k++ {
					rads := getHerofromID(data.Result.Games[i].Scoreboard.Radiant.Picks[k].HeroID, getHeroes)
					rad1 += rads[0] + "|"
					if len(rads) > 1 {
						rad2 += rads[1] + "|"
					} else {
						rad2 += "|"
					}
				}
				rad1 = strings.TrimSuffix(strings.TrimPrefix(rad1, "|"), "|")
				rad2 = strings.TrimSuffix(strings.TrimPrefix(rad2, "|"), "|")
				rads = strings.Split(columnize.Format([]string{rad1, rad2}, &columnize.Config{Glue: " "}), "\n")
			}
			for j := 0; j < len(data.Result.Games[i].Scoreboard.Dire.Players); j++ {
				direNet += data.Result.Games[i].Scoreboard.Dire.Players[j].NetWorth
			}
			if heroes {
				for k := 0; k < len(data.Result.Games[i].Scoreboard.Dire.Picks); k++ {
					dires := getHerofromID(data.Result.Games[i].Scoreboard.Dire.Picks[k].HeroID, getHeroes)
					dire1 += dires[0] + "|"
					if len(dires) > 1 {
						dire2 += dires[1] + "|"
					} else {
						dire2 += "|"
					}
				}
				dire1 = strings.TrimSuffix(strings.TrimPrefix(dire1, "|"), "|")
				dire2 = strings.TrimSuffix(strings.TrimPrefix(dire2, "|"), "|")
				dires = strings.Split(columnize.Format([]string{dire1, dire2}, &columnize.Config{Glue: " "}), "\n")
			}
			worth = radiantNet - direNet
			if showTowers {
				radTower, direTower = towerToString(data.Result.Games[i].Scoreboard.Radiant.TowerState, data.Result.Games[i].Scoreboard.Dire.TowerState)
			}
			if dire == "" {
				dire = "Dire"
			}
			if radiant == "" {
				radiant = "Radiant"
			}

			direTower = strings.TrimSpace(direTower)

			max := len([]rune(direTower))

			if worth != 0 {
				if heroes {
					max = len(dires[0])
					direTower = bloatName(strings.TrimSpace(direTower), max-1)
					dire = bloatName(dire, max-1)
				} else {
					dire = bloatName(dire, max-1)
				}
			}
			if showTowers {
				dire = width.Widen.String(dire)
				radiant = width.Widen.String(radiant)
			}

			str = append(str, fmt.Sprintf("**Dota 2 - %s - Game %d - League: %s - %d viewers**", t, game, leaguename, viewers))

			if showScore && worth != 0 {
				if heroes {
					herostr = append(herostr, fmt.Sprintf("| %s | | %s", rads[0], dires[0]))
					herostr = append(herostr, fmt.Sprintf("| %s | | %s", rads[1], dires[1]))
				}
				if worth > 0 {
					str = append(str, fmt.Sprintf("net %d | %s | %d-%d | %s |",
						worth,
						radiant,
						radiantScore,
						direScore,
						dire))
				} else {
					str = append(str, fmt.Sprintf("| %s | %d-%d | %s | net %.0f",
						radiant,
						radiantScore,
						direScore,
						dire,
						math.Abs(float64(worth))))
				}
				if showTowers {
					// str = append(str, fmt.Sprintf("| %s | | %s | %s Game %d League: %s",
					// 	radTower,
					// 	direTower,
					// 	t,
					// 	game, leaguename))
					str = append(str, fmt.Sprintf("| %s | | %s |",
						radTower,
						direTower))
				}
				str = append(str, herostr...)
				// str = append(str, fmt.Sprintf("%s | %s", herostr, leaguename))

				// str = append(str, fmt.Sprintf("```%s Game %d -- %d viewers -- League: %s",
				// 	t,
				// 	game,
				// 	viewers,
				// 	leaguename))
			} else {
				if heroes {
					herostr = append(herostr, fmt.Sprintf("%s - %s", rad1, dire1))
					herostr = append(herostr, fmt.Sprintf("%s - %s", rad2, dire2))
				}
				str = append(str, fmt.Sprintf("%s - %s",
					radiant,
					dire))

				str = append(str, herostr...)
			}
		}
	}
	if len(str) == 0 {
		msg = append(msg, fmt.Sprintf("No games found."))
		return msg, nil
	}
	return str, nil
}

func towerToString(rad int, dire int) (radTower string, direTower string) {
	towerUp := "♜"
	towerDown := "♖"
	ancient := "♚"
	radstr := ancient + fmt.Sprintf("%011b", int64(rad))
	radstr = organise(radstr)
	radstr = strings.Replace(radstr, "0", towerDown, -1)
	radstr = strings.Replace(radstr, "1", towerUp, -1)

	direstr := ancient + fmt.Sprintf("%011b", int64(dire))
	direstr = organise(direstr)
	var tempstr string
	for _, v := range direstr {
		tempstr = string(v) + tempstr // because allow runes and unicode
	}
	direstr = tempstr
	direstr = strings.Replace(direstr, "0", towerDown, -1)
	direstr = strings.Replace(direstr, "1", towerUp, -1)

	return radstr, direstr
}

func organise(in string) (out string) {
	// volvo returns towers grouped by top/mid/bot, when we want towers grouped by tier
	tempin := []rune(in)
	var tempout []rune
	tempout = append(tempout, tempin[0])
	tempout = append(tempout, ' ')
	tempout = append(tempout, tempin[1])
	tempout = append(tempout, tempin[2])
	tempout = append(tempout, ' ')
	tempout = append(tempout, tempin[3])
	tempout = append(tempout, tempin[6])
	tempout = append(tempout, tempin[9])
	tempout = append(tempout, ' ')
	tempout = append(tempout, tempin[4])
	tempout = append(tempout, tempin[7])
	tempout = append(tempout, tempin[10])
	tempout = append(tempout, ' ')
	tempout = append(tempout, tempin[5])
	tempout = append(tempout, tempin[8])
	tempout = append(tempout, tempin[11])
	out = string(tempout)
	return out
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

func bloatName(inStr string, max int) string {
	if len([]rune(inStr)) < max {
		inStr = strings.Repeat(" ", max-len([]rune(inStr))) + inStr
	}
	return inStr
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
		dotaStr := ""
		res, err := dotamatches(matches[1:])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("an error occured (%s)", err))
			return
		}

		if len(res) < 4 {
			dotaStr = dotaStr + columnize.SimpleFormat(res)
		} else {
			var game []string
			for _, line := range res {
				if strings.Contains(line, "Dota 2") {
					dotaStr += columnize.Format(game, &columnize.Config{
						NoTrim: true,
					})
					dotaStr += line + "\n"
					game = []string{}
				} else {
					game = append(game, line)
				}
			}
			dotaStr += columnize.Format(game, &columnize.Config{
				NoTrim: true,
			})
		}

		fmtstr := fmt.Sprintf("%s", dotaStr)
		s.ChannelMessageSend(m.ChannelID, fmtstr)

	}
}
