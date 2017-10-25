package cli

import (
	"log"

	"github.com/venyii/instabot/provider"
	"github.com/venyii/instabot/slack"
)

func Run(watchDog provider.GoodBoy, ig provider.Provider, sc slack.Sender, dryRun bool) {
	for {
		log.Println("Looking for new media...")

		msgs, err := watchDog.Sniff(ig)

		if err != nil {
			log.Println(err)
		} else {
			if len(msgs) == 0 {
				log.Println("No new media found")
			} else {
				if dryRun {
					log.Printf("dry-run enabled, would send %d new messages\n", len(msgs))
				} else {
					log.Printf("Sending %d new messages\n", len(msgs))
					sc.Send(msgs)
				}
			}
		}

		if dryRun {
			break
		}

		watchDog.Sleep()
	}
}
