package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/venyii/instabot/cfg"
	"github.com/venyii/instabot/cli"
	"github.com/venyii/instabot/instagram"
	"github.com/venyii/instabot/provider"
	"github.com/venyii/instabot/provider/cache"
	"github.com/venyii/instabot/slack"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	configFile := flag.String("config", "config.json", "The config.json path")
	dryRun := flag.Bool("dry-run", false, "Only show which media would be posted with the current cache")
	flag.Parse()

	config, err := cfg.NewConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	printOptions(config, *dryRun)

	slackClient := slack.NewClient(config.SlackToken, config.SlackChannel)
	watchDog := provider.NewWatchDog(config.WaitTime)
	ig := createIgProvider(config)

	cli.Run(watchDog, ig, slackClient, *dryRun)
}

func createIgProvider(config *cfg.Config) provider.Provider {
	var hc *http.Client
	if config.ProxyUrl() != nil {
		hc = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(config.ProxyUrl())}}
	} else {
		hc = &http.Client{}
	}

	return instagram.NewProvider(
		cache.NewFileCache("instagram", config.Username),
		*hc,
		config.Username,
		config.MsgColor,
	)
}

func printOptions(config *cfg.Config, dryRun bool) {
	var proxy string
	var dryRunEnabled string
	if config.Proxy != "" {
		proxy = config.Proxy
	} else {
		proxy = "-"
	}

	if dryRun {
		dryRunEnabled = "Yes"
	} else {
		dryRunEnabled = "No"
	}

	log.Println("")
	log.Println("-- Options ---")
	log.Printf("Username:\t%s\n", config.Username)
	log.Printf("MsgColor:\t%s\n", config.MsgColor)
	log.Printf("WaitTime:\t%dm\n", config.WaitTime)
	log.Printf("Channel:\t%s\n", config.SlackChannel)
	log.Printf("Proxy:\t%s\n", proxy)
	log.Printf("DryRun:\t%s\n", dryRunEnabled)
	log.Println("")
}
