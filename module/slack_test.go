package module

import (
	"cloud.google.com/go/bigquery"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/slack-go/slack"
	"reflect"
	"testing"
	"time"
)

func TestBigQueryClient_FindBill(t *testing.T) {
	type fields struct {
		Client *bigquery.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &BigQueryClient{
				Client: tt.fields.Client,
			}
			if got := c.FindBill(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindBill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigQueryClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *BigQueryClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBigQueryClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBigQueryClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigQueryClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSecretManager(t *testing.T) {
	tests := []struct {
		name    string
		want    *SecretManagerClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSecretManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSecretManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecretManager() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSlackClient(t *testing.T) {
	type args struct {
		token     string
		channelID string
	}
	tests := []struct {
		name string
		args args
		want *SlackClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlackClient(tt.args.token, tt.args.channelID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlackClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecretManagerClient_AccessSecret(t *testing.T) {
	type fields struct {
		Client *secretmanager.Client
	}
	type args struct {
		keyName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecretManagerClient{
				Client: tt.fields.Client,
			}
			got, err := c.AccessSecret(tt.args.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccessSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AccessSecret() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_buildBillQuery(t *testing.T) {
	succeedQuery := "SELECT " +
		"invoice.month," +
		"SUM(cost)" +
		"+ SUM(IFNULL((SELECT SUM(c.amount) " +
		"FROM UNNEST(credits) c), 0))" +
		"AS total, (SUM(CAST(cost * 1000000 AS int64)) + SUM(IFNULL((SELECT SUM(CAST(c.amount * 1000000 as int64)) " +
		"FROM UNNEST(credits) c), 0))) / 1000000 " +
		"AS total_exact " +
		"FROM `referenceTable` " +
		"WHERE invoice.month = '202006' " +
		"GROUP BY 1 " +
		"ORDER BY 1 ASC;"

	type args struct {
		referenceTable string
		formattedMonth string
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{
			name: "should build query",
			args: args{
				referenceTable: "referenceTable",
				formattedMonth: "202006",
			},
			want: succeedQuery,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildBillQuery(tt.args.referenceTable, tt.args.formattedMonth); got != tt.want {
				t.Errorf("buildBillQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertFormattedFromTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should convert format",
			args: args{t: time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)},
			want: "202006",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertFormattedFromTime(tt.args.t); got != tt.want {
				t.Errorf("convertFormattedFromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatReferenceTableName(t *testing.T) {
	type args struct {
		projectID      string
		tableName      string
		splitTableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should format reference table name",
			args: args{
				projectID:      "projectID",
				tableName:      "tableName",
				splitTableName: "splitTableName",
			},
			want: "projectID.tableName.splitTableName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatReferenceTableName(tt.args.projectID, tt.args.tableName, tt.args.splitTableName); got != tt.want {
				t.Errorf("formatReferenceTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
