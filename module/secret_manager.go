package module

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/KouT127/billing-slack-notifier/config"
	"golang.org/x/xerrors"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManagerClient struct {
	*secretmanager.Client
}

func NewSecretManager() (*SecretManagerClient, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx, config.GCPClientOptions...)
	if err != nil {
		return nil, xerrors.Errorf("failed to create secretmanager client: %v", err)
	}
	return &SecretManagerClient{
		client,
	}, nil
}

func (c *SecretManagerClient) AccessSecret(keyName string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", config.ProjectNo, keyName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := c.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}
	return string(result.Payload.Data), nil
}
