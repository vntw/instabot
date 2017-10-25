package slack

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/nlopes/slack"
)

type Sender interface {
	Send(msgs []Message)
}

type Client struct {
	instance *slack.Client
	channel  string
}

func NewClient(token, channel string) Client {
	return Client{
		instance: slack.New(token),
		channel:  channel,
	}
}

func (c Client) Send(msgs []Message) {
	var wg sync.WaitGroup

	for _, msg := range msgs {
		log.Printf("Sending: %s - %s\n", msg.Username, msg.ImageUrl)
		wg.Add(1)

		go func(m Message) {
			defer wg.Done()

			if err := c.sendMsg(m); err != nil {
				log.Println("Error sending msg:", err)
			}
		}(msg)
	}

	wg.Wait()
}

func (c Client) sendMsg(m Message) error {
	var err error

	channelId, err := c.resolveChannelId()
	if err != nil {
		return err
	}

	_, _, err = c.instance.PostMessage(channelId, m.Text, m.SlackMessageParams())
	if err != nil {
		return err
	}

	return nil
}

func (c Client) resolveChannelId() (string, error) {
	channels, err := c.instance.GetChannels(true)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == c.channel {
			return channel.ID, nil
		}
	}

	groups, err := c.instance.GetGroups(true)
	if err != nil {
		return "", err
	}

	for _, group := range groups {
		if group.Name == c.channel {
			return group.ID, nil
		}
	}

	return "", errors.New(fmt.Sprintf("could not resolve channel '%s'", c.channel))
}
