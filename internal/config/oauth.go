package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	DefaultSheetsScope = "https://www.googleapis.com/auth/spreadsheets.readonly"
	SheetsCreds        = "SHEETS_CREDS"
	SheetsTokenPath    = "SHEETS_TOKEN_PATH"
)

type TokenStore interface {
	Load() (*oauth2.Token, error)
	Save(*oauth2.Token) error
}

type fileTokenStore struct {
	path string
}

func NewFileTokenStore(path string) TokenStore {
	return &fileTokenStore{path: resolveTokenPath(path)}
}

func (s *fileTokenStore) Load() (*oauth2.Token, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("open token file %q: %w", s.path, err)
	}
	defer f.Close()
	tok := &oauth2.Token{}
	if err = json.NewDecoder(f).Decode(tok); err != nil {
		return nil, fmt.Errorf("decode token file %q: %w", s.path, err)
	}
	return tok, nil
}

func (s *fileTokenStore) Save(token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("ensure token dir %q: %w", filepath.Dir(s.path), err)
	}
	f, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open token file %q for writing: %w", s.path, err)
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(token); err != nil {
		return fmt.Errorf("encode token to %q: %w", s.path, err)
	}
	return nil
}

func LoadOAuthConfig(path string, scopes ...string) (*oauth2.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read oauth credentials from %q: %w", path, err)
	}
	if len(scopes) == 0 {
		scopes = []string{DefaultSheetsScope}
	}
	cfg, err := google.ConfigFromJSON(data, scopes...)
	if err != nil {
		return nil, fmt.Errorf("parse oauth credentials %q: %w", path, err)
	}
	return cfg, nil
}

func initGoogleSheetsConfig() (GoogleSheetsConfig, error) {
	credsPath := strings.TrimSpace(os.Getenv(SheetsCreds))
	if credsPath == "" {
		credsPath = resolveDataPath("creds.json")
	}
	tokenPath := strings.TrimSpace(os.Getenv(SheetsTokenPath))
	if tokenPath == "" {
		tokenPath = resolveDataPath("token.json")
	}

	oauthCfg, err := LoadOAuthConfig(credsPath)
	if err != nil {
		return GoogleSheetsConfig{}, fmt.Errorf("load sheets oauth config: %w", err)
	}

	return GoogleSheetsConfig{
		OAuthConfig: oauthCfg,
		TokenStore:  NewFileTokenStore(tokenPath),
	}, nil
}

func resolveTokenPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return resolveDataPath("token.json")
	}
	if filepath.IsAbs(path) {
		return path
	}
	return resolveDataPath(filepath.Base(path))
}

func resolveDataPath(name string) string {
	return filepath.Join("..", "data", name)
}
