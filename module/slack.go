package module

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

type SlackClient struct {
	*slack.Client
	channelID string
}

func NewSlackClient() *SlackClient {
	token, err := accessSecretVersion("SLACK_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	channelID, err := accessSecretVersion("CHANNEL_ID")
	if err != nil {
		log.Fatal(err)
	}

	api := slack.New(token)
	return &SlackClient{
		api,
		channelID,
	}
}

func (c *SlackClient) NotifyMessage(msg string) error {
	opts := []slack.MsgOption{
		slack.MsgOptionText(msg, true),
	}
	channel, _, _, err := c.SendMessage(c.channelID, opts...)
	if err != nil {
		return err
	}
	fmt.Printf("channel: %s", channel)
	return nil
}
