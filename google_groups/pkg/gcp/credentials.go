// package gcp provides utilities for working with GCP
package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kubeflow/internal-acls/google_groups/pkg/gcp/gcs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// CredentialDirPermMode unix permission max suitable for directory storing credentials
	CredentialDirPermMode = 0700
)

// CredentialHelper defines an interface for getting tokens.
type CredentialHelper interface {
	//GetTokenAndConfig() (*TokenAndConfig, error)
	GetTokenSource(ctx context.Context)(oauth2.TokenSource ,error)

	// GetOAuthConfig returns the OAuth2 client configuration
	GetOAuthConfig()(*oauth2.Config)
}

// WebFlowHelper helps get credentials using the webflow.
type WebFlowHelper struct {
	config *oauth2.Config
	Log logr.Logger
}

// NewWebFlowHelper constructs a new web flow helper. oAuthClientFile should be the path to a credentials.json
// downloaded from the API console.
func NewWebFlowHelper(oAuthClientFile string, scopes []string) (*WebFlowHelper, error) {
	var fHelper gcs.FileHelper

	if strings.HasPrefix(oAuthClientFile, "gs://") {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)

		if err != nil {
			return nil, err
		}

		fHelper = &gcs.GcsHelper {
			Ctx: ctx,
			Client: client,
		}
	} else {
		fHelper = &gcs.LocalFileHelper{}
	}

	reader, err := fHelper.NewReader(oAuthClientFile)

	if err != nil {
		return nil, err

	}
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse client secret file to config")
	}
	return &WebFlowHelper{
		config: config,
		Log: zapr.NewLogger(zap.L()),
	}, nil
}

func (h *WebFlowHelper) GetOAuthConfig() (*oauth2.Config) {
	return h.config
}

// GetTokenSource requests a token from the web, then returns the retrieved token.
func (h *WebFlowHelper) GetTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	authURL := h.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// TODO(jlewi): How to open it automatically?
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.Wrapf(err, "Unable to read authorization code")
	}

	tok, err := h.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve token from web: %v")
	}

	return h.config.TokenSource(ctx, tok), nil
}

// TokenCache defines an interface for caching tokens
type TokenCache interface {
	GetToken() (*oauth2.Token, error)
	Save(token *oauth2.Token) error
}

// FileTokenCache implements caching to a file.
type FileTokenCache struct {
	CacheFile string
	Log logr.Logger
}

func (c *FileTokenCache) GetToken() (*oauth2.Token, error) {
	f, err := os.Open(c.CacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Save saves a token to a file path.
func (c *FileTokenCache) Save(token *oauth2.Token) error {
	c.Log.Info("Saving credential", "file", c.CacheFile)

	dir := filepath.Dir(c.CacheFile)

	_, err := os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			c.Log.Info("Create cache directory", "dir", dir)
			err := os.MkdirAll(dir, CredentialDirPermMode)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	f, err := os.OpenFile(c.CacheFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		c.Log.Error(err, "Unable to cache oauth token: %v")
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

// CachedCredentialHelper is a credential helper that will cache the credential.
type CachedCredentialHelper struct {
	CredentialHelper CredentialHelper
	TokenCache TokenCache
	Log logr.Logger
}

func (h *CachedCredentialHelper) GetOAuthConfig() (*oauth2.Config) {
	return h.CredentialHelper.GetOAuthConfig()
}

func (c *CachedCredentialHelper) GetTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	log := c.Log
	// Try the cache.
	tok, err := c.TokenCache.GetToken()

	if err != nil {
		return nil, err
	}

	if tok == nil {
		// Cache is empty so get a token
		ts, err := c.CredentialHelper.GetTokenSource(context.Background())

		if err != nil {
			return nil, err
		}

		// Save the token
		newTok, err := ts.Token()
		tok = newTok
		if err != nil {
			log.Error(err, "Could generate token from token source")
			return ts, err
		}
		err = c.TokenCache.Save(newTok)

		if err != nil {
			log.Error(err, "Could not save token")
		}
	}

	ts := c.CredentialHelper.GetOAuthConfig().TokenSource(ctx, tok)
	return ts, nil
}