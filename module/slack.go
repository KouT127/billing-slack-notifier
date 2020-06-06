package module

import (
	"fmt"
	"github.com/slack-go/slack"
)

type SlackClient struct {
	*slack.Client
	ChannelID string
}

func NewSlackClient(token, channelID string) *SlackClient {
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
	channel, _, _, err := c.SendMessage(c.ChannelID, opts...)
	if err != nil {
		return err
	}
	fmt.Printf("channel: %s", channel)
	return nil
}
