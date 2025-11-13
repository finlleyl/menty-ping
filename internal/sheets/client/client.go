package sheets

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/finlleyl/menty-ping/internal/config"
)

type Client interface {
	HTTPClient(ctx context.Context) (*http.Client, error)
	SaveToken(ctx context.Context, token *oauth2.Token) error
}

func New(config *oauth2.Config, store config.TokenStore) Client {
	return &googleClient{oauthConfig: config, tokens: store}
}

type googleClient struct {
	oauthConfig *oauth2.Config
	tokens      config.TokenStore
}

func (c *googleClient) HTTPClient(ctx context.Context) (*http.Client, error) {
	tok, err := c.tokens.Load()
	if err != nil {
		return nil, fmt.Errorf("load oauth token: %w", err)
	}
	if tok != nil && !tok.Valid() {
		if tok, err = c.refresh(ctx, tok); err != nil {
			return nil, fmt.Errorf("refresh oauth token: %w", err)
		}
	}
	return c.oauthConfig.Client(ctx, tok), nil
}

func (c *googleClient) SaveToken(ctx context.Context, token *oauth2.Token) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}
	if !token.Valid() {
		var err error
		if token, err = c.refresh(ctx, token); err != nil {
			return err
		}
	}
	return c.tokens.Save(token)
}

func (c *googleClient) refresh(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	tok, err := c.oauthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, fmt.Errorf("refresh oauth token: %w", err)
	}
	if err = c.tokens.Save(tok); err != nil {
		return nil, fmt.Errorf("persist refreshed token: %w", err)
	}
	return tok, nil
}
