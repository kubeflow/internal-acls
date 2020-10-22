package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SecretCache implements a cache for an OAuth2 credential using GCP secret manager
type SecretCache struct {
	client *secretmanager.Client
	Project string
	Secret string
	Version string
	Log logr.Logger

}

func NewSecretCache(project string,  secret string, version string) (*SecretCache, error) {
	c := &SecretCache{
		Project: project,
		Secret:  secret,
		Version: version,
		Log:     zapr.NewLogger(zap.L()),
	}

	client, err := secretmanager.NewClient(context.Background())

	if err != nil {
		return nil, err
	}

	c.client = client
	return c, nil
}

func (c *SecretCache) GetToken() (*oauth2.Token, error) {
	payload, err := c.loadSecret()

	if err != nil {
		return nil, err
	}

	// No secret has been saved yet.
	if payload == nil {
		return nil, nil
	}

	tok := &oauth2.Token{}
	err = json.Unmarshal(payload, tok)
	return tok, err
}

// Save saves a token to a file path.
func (c *SecretCache) Save(token *oauth2.Token) error {
	log := c.Log
	log.Info("Saving credential to secret manager", "project", c.Project, "secret", c.Secret)

	// Create the request to create the secret.
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", c.Project),
		SecretId: c.Secret,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	ctx := context.Background()

	_, err := c.client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		status, ok := status.FromError(err)

		if !ok {
			log.Error(err, "Error creating secret.", "project", c.Project, "secret", c.Secret)
			return err
		}

		if status.Code() == codes.AlreadyExists {
			log.Info("Secret exists", "project", c.Project, "secret", c.Secret)
		} else {
			log.Error(err, "Error creating secret.", "project", c.Project, "secret", c.Secret)
			return err
		}
	}

	payload, err := json.Marshal(token)
	if err != nil {
		return err
	}


	// Build the request.
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%v/secrets/%v", c.Project, c.Secret),
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	version, err := c.client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return err
	}

	log.Info("Stored token in secret manager,", "version", version)
	return nil
}

func (c *SecretCache) loadSecret() ([]byte, error) {
	log := c.Log
	name := fmt.Sprintf("projects/%v/secrets/%v/versions/%s", c.Project, c.Secret, c.Version)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	ctx := context.Background()

	// Call the API.
	result, err := c.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		status, ok := status.FromError(err)

		if ok {
			if status.Code() == codes.NotFound {
				log.Info("No secret exists containing cached token", "secret", name)
				return nil, nil
			}

			if status.Code() == codes.FailedPrecondition {
				log.Info("Latest version of secret is not valid a new secret will be created.", "secret", name, "status_message", status.Message())
				return nil, nil
			}
		}
		return nil, errors.Wrapf(err, "failed to access secret %v", name)
	}

	return result.Payload.Data, nil
}