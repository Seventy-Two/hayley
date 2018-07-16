package main

import (
	"io/ioutil"
	"log"

	cli "github.com/jawher/mow.cli"
	yaml "gopkg.in/yaml.v2"
)

const (
	appName        = "Hayley"
	appDescription = "Discord chat bot with various commands"
)

var config keys

func main() {

	app := cli.App(appName, appDescription)

	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal("you fucked up")
	}
	err = yaml.Unmarshal(data, &config)

	start(app)
}
