package provider

import (
	"github.com/venyii/instabot/slack"
)

type Provider interface {
	// A list of messages sorted by date descending
	Latest(dryRun bool, dummy bool) ([]slack.Message, error)
}
