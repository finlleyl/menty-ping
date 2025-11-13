package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	TelegramAppID = "TELEGRAM_APP_ID"
	TelegramAppHash = "TELEGRAM_APP_HASH"
)

func initTelegramConfig() (TelegramConfig, error) {
	appIDRaw := strings.TrimSpace(os.Getenv(TelegramAppID))
	if appIDRaw == "" {
		return TelegramConfig{}, fmt.Errorf("env %s is required", TelegramAppID)
	}

	appID, err := strconv.Atoi(appIDRaw)
	if err != nil {
		return TelegramConfig{}, fmt.Errorf("parse %s: %w", TelegramAppID, err)
	}

	appHash := strings.TrimSpace(os.Getenv(TelegramAppHash))
	if appHash == "" {
		return TelegramConfig{}, fmt.Errorf("env %s is required", TelegramAppHash)
	}

	return TelegramConfig{
		AppID:       appID,
		AppHash:     appHash,
	}, nil
}
