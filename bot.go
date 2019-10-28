package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seventy-two/hayley/commands/beer"
	"github.com/seventy-two/hayley/commands/shitpost"
	"github.com/seventy-two/hayley/commands/teamspeak"

	"github.com/bwmarrin/discordgo"
	cli "github.com/jawher/mow.cli"
	"github.com/seventy-two/hayley/commands/dictionary"
	"github.com/seventy-two/hayley/commands/divegrass"
	"github.com/seventy-two/hayley/commands/dota"
	"github.com/seventy-two/hayley/commands/math"
	"github.com/seventy-two/hayley/commands/movie"
	"github.com/seventy-two/hayley/commands/nfl"
	"github.com/seventy-two/hayley/commands/quotes"
	"github.com/seventy-two/hayley/commands/siege"
	"github.com/seventy-two/hayley/commands/stocks"
	"github.com/seventy-two/hayley/commands/tv"
	"github.com/seventy-two/hayley/commands/urbandictionary"
	"github.com/seventy-two/hayley/commands/weather"
	"github.com/seventy-two/hayley/commands/youtube"
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

	dg.AddHandler(logger)

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
	if services.siegeAPI != nil {
		siege.RegisterService(dg, services.siegeAPI)
	}
	if services.teamspeakAPI != nil {
		teamspeak.RegisterService(dg, services.teamspeakAPI)
	}
	if services.beerAPI != nil {
		beer.RegisterService(dg, services.beerAPI)
	}
	if services.quotesAPI != nil {
		quotes.RegisterService(dg, services.quotesAPI)
	}

	shitpost.RegisterService(dg)

}

func logger(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.ID + " | " + m.Author.Username + " | " + m.Content)
}
