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
	instance  *slack.Client
	channelId string
}

func NewClient(token, channel string) Client {
	c := &Client{instance: slack.New(token)}
	c.resolveChannelId(channel)

	return *c
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
	_, _, err := c.instance.PostMessage(c.channelId, m.Text, m.SlackMessageParams())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) resolveChannelId(name string) error {
	channels, err := c.instance.GetChannels(true)
	if err != nil {
		return err
	}

	for _, channel := range channels {
		if channel.Name == name {
			c.channelId = channel.ID
			return nil
		}
	}

	groups, err := c.instance.GetGroups(true)
	if err != nil {
		return err
	}

	for _, group := range groups {
		if group.Name == name {
			c.channelId = group.ID
			return nil
		}
	}

	return errors.New(fmt.Sprintf("could not resolve channel '%s'", name))
}
