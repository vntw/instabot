package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/venyii/instabot/cache"
	"github.com/venyii/instabot/cfg"
	"github.com/venyii/instabot/instagram"
	"github.com/venyii/instabot/slack"
)

var (
	cacheStrat  cache.Cache
	slackClient slack.Client
	configFile  = flag.String("config", "config.json", "The config.json path")
	dryRun      = flag.Bool("dry-run", false, "Only show which media would be posted with the current cache")
	dummy       = flag.Bool("dummy", false, "Dummy data mode?")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	config, err := cfg.NewConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}
	printOptions(config)

	cacheStrat = cache.NewFileCache()
	slackClient = slack.NewClient(config.SlackToken, config.SlackChannel)

	for {
		log.Println("Looking for new media...")
		observeInstagram(*config)

		if *dryRun {
			break
		}

		waitTime := calcWaitTime(config.WaitTime)
		log.Printf("Sleeping for %v minutes...\n", waitTime)
		time.Sleep(waitTime)
	}
}

func observeInstagram(config cfg.Config) {
	p := instagram.NewProvider(config, cacheStrat)
	msgs, err := p.Latest(*dryRun, *dummy)

	if err != nil {
		log.Println(err)
	}

	if *dryRun {
		log.Printf("%#v\n", msgs)
		return
	}

	slackClient.Send(msgs)
}

func calcWaitTime(waitTime uint8) time.Duration {
	min, max := 0, 60
	randomSeconds := rand.Intn(max-min) + min
	variance := time.Duration(randomSeconds) * time.Second

	return time.Duration(waitTime)*time.Minute + variance
}

func printOptions(config *cfg.Config) {
	var proxy string
	if config.Proxy != "" {
		proxy = config.Proxy
	} else {
		proxy = "-"
	}

	log.Println("")
	log.Println("-- Options ---")
	log.Printf("WaitTime:\t%dm\n", config.WaitTime)
	log.Printf("Channel:\t%s\n", config.SlackChannel)
	log.Printf("Proxy:\t%s\n", proxy)
	log.Println("")
}
