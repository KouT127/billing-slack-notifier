package module

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
	"os"
	"testing"
)

func TestSecretManagerClient_AccessSecret(t *testing.T) {
	projectNo := os.Getenv("PROJECT_NO")
	json := os.Getenv("SERVICE_ACCOUNT_JSON")
	m, err := NewSecretManager(projectNo, option.WithCredentialsJSON([]byte(json)))
	if err != nil {
		t.Fatalf("%v", err)
	}
	type fields struct {
		projectNo string
		Client    *secretmanager.Client
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
		{
			name: "should access secret",
			fields: fields{
				projectNo: m.projectNo,
				Client:    m.Client,
			},
			args: args{
				keyName: "TEST_SECRET",
			},
			want:    "SECRET",
			wantErr: false,
		},
		{
			name: "should not access secret when not exists key",
			fields: fields{
				projectNo: m.projectNo,
				Client:    m.Client,
			},
			args: args{
				keyName: "NOT_EXISTS",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SecretManagerClient{
				projectNo: tt.fields.projectNo,
				Client:    tt.fields.Client,
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
