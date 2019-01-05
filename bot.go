package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	cli "github.com/jawher/mow.cli"
	"github.com/seventy-two/Hayley/commands/dictionary"
	"github.com/seventy-two/Hayley/commands/divegrass"
	"github.com/seventy-two/Hayley/commands/dota"
	"github.com/seventy-two/Hayley/commands/math"
	"github.com/seventy-two/Hayley/commands/movie"
	"github.com/seventy-two/Hayley/commands/nfl"
	"github.com/seventy-two/Hayley/commands/stocks"
	"github.com/seventy-two/Hayley/commands/tv"
	"github.com/seventy-two/Hayley/commands/urbandictionary"
	"github.com/seventy-two/Hayley/commands/weather"
	"github.com/seventy-two/Hayley/commands/youtube"
)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Listening to prefix !")
}

func start(app *cli.Cli, services *serviceConfig) {
	dg, _ := discordgo.New(fmt.Sprintf("Bot %s", services.discordAPI.APIKey))

	go registerServices(dg, services)

	dg.AddHandler(ready)

	err := dg.Open()

	if err != nil {
		log.Fatalf("Error opening Discord session: %s", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func registerServices(dg *discordgo.Session, services *serviceConfig) {
	if services.dictionaryAPI != nil {
		dictionary.RegisterService(dg, services.dictionaryAPI)
	}
	if services.dotaAPI != nil {
		dota.RegisterService(dg, services.dotaAPI)
	}
	if services.divegrassAPI != nil {
		divegrass.RegisterService(dg, services.divegrassAPI)
	}
	if services.nflAPI != nil {
		nfl.RegisterService(dg, services.nflAPI)
	}
	if services.movieAPI != nil {
		movie.RegisterService(dg, services.movieAPI)
	}
	if services.stocksAPI != nil {
		stocks.RegisterService(dg, services.stocksAPI)
	}
	if services.tvAPI != nil {
		tv.RegisterService(dg, services.tvAPI)
	}
	if services.urbanAPI != nil {
		urbandictionary.RegisterService(dg, services.urbanAPI)
	}
	if services.weatherAPI != nil {
		weather.RegisterService(dg, services.weatherAPI)
	}
	if services.mathAPI != nil {
		math.RegisterService(dg, services.mathAPI)
	}
	if services.youtubeAPI != nil {
		youtube.RegisterService(dg, services.youtubeAPI)
	}
}
