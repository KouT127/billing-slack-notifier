package module

import (
	"github.com/slack-go/slack"
	"testing"
)

func TestSlackClient_NotifyMessage(t *testing.T) {
	type fields struct {
		Client    *slack.Client
		ChannelID string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SlackClient{
				Client:    tt.fields.Client,
				ChannelID: tt.fields.ChannelID,
			}
			if err := c.NotifyMessage(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("NotifyMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
