package config

import (
	"errors"
	"sync"

	"golang.org/x/oauth2"
)

type (
	Config struct {
		Google   GoogleSheetsConfig
		Telegram TelegramConfig
	}
	GoogleSheetsConfig struct {
		OAuthConfig *oauth2.Config
		TokenStore  TokenStore
	}
	TelegramConfig struct {
		AppID   int
		AppHash string
	}
)

func NewConfig() (*Config, error) {
	var (
		googleCfg   GoogleSheetsConfig
		telegramCfg TelegramConfig
		googleErr   error
		telegramErr error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		googleCfg, googleErr = initGoogleSheetsConfig()
	}()

	go func() {
		defer wg.Done()
		telegramCfg, telegramErr = initTelegramConfig()
	}()

	wg.Wait()

	if err := errors.Join(googleErr, telegramErr); err != nil {
		return nil, err
	}

	return &Config{
		Google:   googleCfg,
		Telegram: telegramCfg,
	}, nil
}
