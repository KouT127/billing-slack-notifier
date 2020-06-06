package module

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManagerClient struct {
	projectNo string
	*secretmanager.Client
}

func NewSecretManager(projectNo string, opts ...option.ClientOption) (*SecretManagerClient, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, xerrors.Errorf("failed to create secretmanager client: %v", err)
	}
	return &SecretManagerClient{
		projectNo,
		client,
	}, nil
}

func (c *SecretManagerClient) AccessSecret(keyName string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.projectNo, keyName)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := c.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", xerrors.Errorf("failed to access secret version: %v", err)
	}
	return string(result.Payload.Data), nil
}
