package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func main() {
	client := NewSlackClient()
	err := client.sendMessage("New Test")
	if err != nil {
		log.Fatalf("%v", err)
	}
}

type SlackClient struct {
	slack.Client
	channelID string
}

func NewSlackClient() *SlackClient {
	token := os.Getenv("SLACK_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")
	api := slack.New(token)
	return &SlackClient{
		*api,
		channelID,
	}
}

func (c *SlackClient) sendMessage(msg string) error {
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
