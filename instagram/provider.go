package instagram

import (
	"fmt"
	"net/http"

	"github.com/venyii/instabot/provider/cache"
	"github.com/venyii/instabot/slack"
)

const (
	name     = "Instagram"
	baseUrl  = "https://www.instagram.com"
	maxNodes = 12
)

type Provider struct {
	cache      cache.Cache
	httpClient http.Client
	username   string
	msgColor   string
}

func NewProvider(c cache.Cache, hc http.Client, username string, msgColor string) *Provider {
	return &Provider{
		cache:      c,
		httpClient: hc,
		username:   username,
		msgColor:   msgColor,
	}
}

func (p Provider) Cache() cache.Cache {
	return p.cache
}

func (p Provider) Latest() ([]slack.Message, error) {
	var data Data
	var err error

	f := NewFetcher(p.httpClient, p.username)

	data, err = f.fetchUserData()

	if err != nil {
		return nil, err
	}

	up := data.UserProfile()
	var newMsgs []slack.Message

	for _, node := range up.MediaNodes(maxNodes) {
		newMsgs = append(newMsgs, p.createMessage(up, node))
	}

	return newMsgs, nil
}

func (p Provider) createMessage(u User, n Node) slack.Message {
	m := slack.NewMessage(name, u.Username)
	m.Text = n.DetailUrl()
	m.Author = u.Fullname
	m.AuthorIcon = u.ProfilePic
	m.AuthorLink = u.ProfileUrl()
	m.ImageUrl = n.Src
	m.Date = int64(n.Date)
	m.Color = p.msgColor
	m.Fallback = fmt.Sprintf("New Instagram post by %s", u.Username)

	return *m
}
