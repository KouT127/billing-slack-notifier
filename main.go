package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/slack-go/slack"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	client := NewSlackClient()

	results := queryBill()
	for _, result := range results {
		err := client.sendMessage(result)
		if err != nil {
			log.Fatalf("%v", err)
		}
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

func queryBill() []string {
	ctx := context.Background()
	json := os.Getenv("SERVICE_ACCOUNT_JSON")
	projectID := os.Getenv("PROJECT_ID")

	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(json)))
	if err != nil {
		log.Fatalf("%v", err)
	}

	query := fmt.Sprintf("SELECT "+
		"invoice.month,"+
		"SUM(cost)"+
		"+ SUM(IFNULL((SELECT SUM(c.amount) "+
		"FROM UNNEST(credits) c), 0))"+
		"AS total, (SUM(CAST(cost * 1000000 AS int64)) + SUM(IFNULL((SELECT SUM(CAST(c.amount * 1000000 as int64)) "+
		"FROM UNNEST(credits) c), 0))) / 1000000 "+
		"AS total_exact "+
		"FROM `%s.biling.gcp_billing_export_v1_01C98C_7E00E6_392CCC` "+
		"GROUP BY 1 "+
		"ORDER BY 1 ASC;", projectID)

	q := client.Query(query)

	rows, err := q.Read(ctx)
	if err != nil {
		log.Fatal(err)
	}

	results := make([]string, rows.TotalRows)
	idx := 0
	for {
		var values []bigquery.Value
		err := rows.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		results[idx] = fmt.Sprintf("%s　利用金額: %f", values[0], values[1])
		idx++
	}

	return results
}
