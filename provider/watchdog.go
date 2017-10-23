package provider

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/venyii/instabot/slack"
)

type WatchDog struct {
	baseWaitTime uint8
}

func NewWatchDog(baseWaitTime uint8) WatchDog {
	return WatchDog{baseWaitTime: baseWaitTime}
}

func (wd WatchDog) Sniff(p Provider) ([]slack.Message, error) {
	msgs, err := p.Latest()
	if err != nil {
		return []slack.Message{}, err
	}

	if len(msgs) == 0 {
		return msgs, nil
	}

	var newMsgs []slack.Message
	lastDate, err := p.Cache().ReadLastDate()
	if err != nil {
		log.Fatal(err)
	}

	if lastDate == -1 {
		// Nothing cached, post only the last img and build cache file
		newMsgs = msgs[:1]
	} else {
		for _, msg := range msgs {
			if msg.Date > lastDate {
				newMsgs = append(newMsgs, msg)
			} else {
				break
			}
		}
	}

	newLastDate := msgs[0].Date
	if newLastDate != lastDate {
		if err := p.Cache().WriteLastDate(newLastDate); err != nil {
			return []slack.Message{}, errors.New("could not write cache file " + err.Error())
		}
	}

	return newMsgs, nil
}

func (wd WatchDog) Sleep() {
	min, max := 0, 60
	randomSeconds := rand.Intn(max-min) + min
	variance := time.Duration(randomSeconds) * time.Second

	waitTime := time.Duration(wd.baseWaitTime)*time.Minute + variance

	log.Printf("Sleeping for %v minutes...\n", waitTime)
	time.Sleep(waitTime)
}
