package provider

import (
	"github.com/venyii/instabot/provider/cache"
	"github.com/venyii/instabot/slack"
)

type Provider interface {
	// A list of messages sorted by date descending
	Latest() ([]slack.Message, error)
	// The provider specific cache
	Cache() cache.Cache
}
