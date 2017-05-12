package instagram

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/venyii/instabot/cache"
	"github.com/venyii/instabot/cfg"
	"github.com/venyii/instabot/slack"
)

const (
	name     = "Instagram"
	baseUrl  = "https://www.instagram.com"
	maxNodes = 12
)

type Provider struct {
	cfg   cfg.Config
	cache cache.Cache
}

func NewProvider(cfg cfg.Config, c cache.Cache) *Provider {
	return &Provider{
		cfg:   cfg,
		cache: c,
	}
}

func (p Provider) Latest(dryRun bool, dummy bool) ([]slack.Message, error) {
	var data Data
	var err error

	f := NewFetcher(p.cfg.ProxyUrl())

	if dummy {
		data, err = f.ExtractDataFromHtml(getDummyHtml())
	} else {
		data, err = f.FetchUserData(p.cfg.Username)
	}

	if err != nil {
		return nil, err
	}

	lastDate := p.cache.ReadLastDate()
	newMsgs := p.findNewMsgs(data.UserProfile(), lastDate)

	log.Printf("[IG] New: %d\n", len(newMsgs))

	if dryRun || len(newMsgs) == 0 {
		return newMsgs, nil
	}

	newLastDate := newMsgs[0].Date

	if newLastDate != lastDate {
		if err := p.cache.WriteLastDate(newLastDate); err != nil {
			log.Println("Could not write cache file")
			log.Fatalln(err)
		}
	}

	return newMsgs, nil
}

func (p Provider) findNewMsgs(u User, lastDate int64) []slack.Message {
	nodes := u.MediaNodes(maxNodes)

	if len(nodes) == 0 {
		return []slack.Message{}
	}

	if lastDate == -1 {
		// Nothing cached, post only the last img and build cache file
		return []slack.Message{p.createMessage(u, nodes[:1][0])}
	}

	var newMsgs []slack.Message

	for _, node := range nodes {
		if node.Date > lastDate {
			newMsgs = append(newMsgs, p.createMessage(u, node))
		}
	}

	return newMsgs
}

func (p Provider) createMessage(u User, n Node) slack.Message {
	m := slack.NewMessage(name)
	m.Text = n.DetailUrl()
	m.Author = u.Username
	m.AuthorIcon = u.ProfilePic
	m.AuthorLink = u.ProfileUrl()
	m.ImageUrl = n.Src
	m.Date = int64(n.Date)
	m.Color = "#fddd4a"
	m.Fallback = fmt.Sprintf("New Instagram post by %s", u.Username)

	return *m
}

func getDummyHtml() []byte {
	dummyResp, err := ioutil.ReadFile("instagram/example/response.html")
	if err != nil {
		log.Fatal(err)
	}

	return dummyResp
}
