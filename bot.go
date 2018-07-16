package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	cli "github.com/jawher/mow.cli"
	"github.com/utilitywarehouse/uwgolib/log"
)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Listening to !w, !d2, !m and !nfl")
}

type keys struct {
	botKey     string
	dota       string
	geocode    string
	weather    string
	divergrass string
	dictionary string
	omdb       string
	wolfram    string
}

func start(app *cli.Cli) {

	botKey := app.String(cli.StringOpt{
		Name:   "botKey",
		Desc:   "Discord bot key combo",
		Value:  config.botKey,
		EnvVar: "DISCORD_BOT_KEY",
	})

	dg, _ := discordgo.New(fmt.Sprintf("Bot %s", botKey))

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	err := dg.Open()

	if err != nil {
		log.Error("Error opening Discord session: ", err)
	}

	log.Info("Hayley is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
