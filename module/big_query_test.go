package module

import (
	"testing"
	"time"
)

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
