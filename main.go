package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/compute/metadata"
	"context"
	"fmt"
	"github.com/slack-go/slack"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

func main() {
	slackClient := NewSlackClient()
	bigQueryClient := NewBigQueryClient()

	results := bigQueryClient.queryBill()
	for _, result := range results {
		err := slackClient.sendMessage(result)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

type SlackClient struct {
	*slack.Client
	channelID string
}

func NewSlackClient() *SlackClient {
	token := os.Getenv("SLACK_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")
	api := slack.New(token)
	return &SlackClient{
		api,
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

type BigQueryClient struct {
	*bigquery.Client
}

func NewBigQueryClient() *BigQueryClient {
	ctx := context.Background()
	json := os.Getenv("SERVICE_ACCOUNT_JSON")
	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Fatal(err)
	}

	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(json)))
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &BigQueryClient{
		client,
	}
}

func (c *BigQueryClient) queryBill() []string {
	ctx := context.Background()
	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Fatal(err)
	}
	tableName := "biling"
	splitTableName := os.Getenv("SPLIT_TABLE_NAME")
	referenceTable := formatReferenceTableName(projectID, tableName, splitTableName)
	formattedMonth := convertFormattedFromTime(time.Now())

	query := buildBillQuery(referenceTable, formattedMonth)
	q := c.Query(query)

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

func convertFormattedFromTime(t time.Time) string {
	return t.Format("200601")
}

func formatReferenceTableName(projectID, tableName, splitTableName string) string {
	return fmt.Sprintf("%s.%s.%s", projectID, tableName, splitTableName)
}

func buildBillQuery(referenceTable, formattedMonth string) string {
	return fmt.Sprintf("SELECT "+
		"invoice.month,"+
		"SUM(cost)"+
		"+ SUM(IFNULL((SELECT SUM(c.amount) "+
		"FROM UNNEST(credits) c), 0))"+
		"AS total, (SUM(CAST(cost * 1000000 AS int64)) + SUM(IFNULL((SELECT SUM(CAST(c.amount * 1000000 as int64)) "+
		"FROM UNNEST(credits) c), 0))) / 1000000 "+
		"AS total_exact "+
		"FROM `%s` "+
		"WHERE invoice.month = '%s' "+
		"GROUP BY 1 "+
		"ORDER BY 1 ASC;", referenceTable, formattedMonth)
}
