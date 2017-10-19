package slack

import (
	"encoding/json"
	"strconv"

	"github.com/nlopes/slack"
)

type Message struct {
	Username string

	Fallback string
	Text     string

	Author     string
	AuthorLink string
	AuthorIcon string

	ImageUrl     string
	ProviderName string
	Color        string

	Date int64
}

func NewMessage(provider string, username string) *Message {
	return &Message{ProviderName: provider, Username: username}
}

func (m Message) SlackMessageParams() slack.PostMessageParameters {
	p := slack.NewPostMessageParameters()
	p.AsUser = true
	p.Attachments = append(p.Attachments, slack.Attachment{
		Fallback:   m.Fallback,
		Color:      m.Color,
		ImageURL:   m.ImageUrl,
		AuthorName: m.Author,
		AuthorLink: m.AuthorLink,
		AuthorIcon: m.AuthorIcon,
		Footer:     m.ProviderName,
		Ts:         json.Number(strconv.FormatInt(m.Date, 10)),
	})

	return p
}
