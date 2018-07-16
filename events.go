package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Seventy-Two/Hayley/commands/dictionary"
	"github.com/Seventy-Two/Hayley/commands/divegrass"
	"github.com/Seventy-Two/Hayley/commands/dota"
	"github.com/Seventy-Two/Hayley/commands/nfl"
	"github.com/Seventy-Two/Hayley/commands/omdb"
	"github.com/Seventy-Two/Hayley/commands/tvmaze"
	"github.com/Seventy-Two/Hayley/commands/urbandictionary"
	"github.com/Seventy-Two/Hayley/commands/weather"
	"github.com/Seventy-Two/Hayley/commands/wolfram"
	"github.com/Seventy-Two/Hayley/commands/youtube"
	"github.com/bwmarrin/discordgo"
	"github.com/ryanuber/columnize"
	"github.com/seventy-two/Hayley/commands/stocks"
)

var set int
var inprogress = false
var d2interval = time.Duration(3)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Println(fmt.Sprintf("[%5s]: %5s > %s\n", m.ChannelID, m.Author.Username, m.Content))

	matches := strings.Split(m.Content, " ")
	var str string
	var err error

	switch matches[0] {
	case "!nfl":
		res, err := nfl.Nfl()
		if err != nil {
			log.Println(err)
		}
		str = columnize.SimpleFormat(res)
	case "!d2":
		if inprogress == true {
			return
		}
		inprogress = true
		go func() {
			for inprogress {
				dotaStr := ""

				res, err := dota.Dotamatches(matches[1:])
				if err != nil {
					log.Println(err)
				}

				if len(res) == 1 {
					time.Sleep(5 * time.Minute)
					continue
				}

				msgs, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")

				for _, message := range msgs {
					if (message.Author.Bot && strings.Contains(message.Content, "Dota 2")) || strings.Contains(message.Content, "!d2") {

						s.ChannelMessageDelete(m.ChannelID, message.ID)
					}
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

				fmtstr := fmt.Sprintf("```%s```", dotaStr)
				s.ChannelMessageSend(m.ChannelID, fmtstr)

				time.Sleep(d2interval * time.Minute)
			}
			return
		}()

	case "!d2off":
		inprogress = false
	case "!d2interval":
		dur, err := strconv.Atoi(strings.Join(matches, ""))
		if err != nil {
			return
		}
		d2interval = time.Duration(int64(dur))
	case "!w":
		str, err = weather.Weather(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!f":
		str, err = weather.Forecast(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!m":
		str, err = omdb.Omdb(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!ud":
		str, err = urbandictionary.Urban(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!stocks":
		str, err = stocks.GetStock(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!d":
		str, err = dictionary.Dict(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!wotd":
		str, err = dictionary.Wotd()
		if err != nil {
			log.Println(err)
		}
	case "!wa":
		str, err = wolfram.Wolfram(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!tv":
		str, err = tvmaze.Tvmaze(matches[1:])
		if err != nil {
			log.Println(err)
		}
	case "!yt":
		str, err = youtube.Youtube(matches[1:])
		if err != nil {
			log.Println(err)
		}
		fmtstr := fmt.Sprintf("%s", str)
		s.ChannelMessageSend(m.ChannelID, fmtstr)
		return
	case "!foot":
		res, err := divegrass.Divegrass()
		if err != nil {
			log.Println(err)
		}
		str = columnize.SimpleFormat(res)
	default:
		return
	}

	if str != "" {
		fmtstr := fmt.Sprintf("```%s```", str)
		s.ChannelMessageSend(m.ChannelID, fmtstr)
	}
}
