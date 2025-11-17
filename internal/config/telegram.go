package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	TelegramAppID       = "TELEGRAM_APP_ID"
	TelegramAppHash     = "TELEGRAM_APP_HASH"
	TelegramPhoneNumber = "TELEGRAM_PHONE_NUMBER"
	TelegramPassword    = "TELEGRAM_PASSWORD"
)

type TelegramConfig struct {
	AppID   int
	AppHash string
	PhoneNumber string
	Password string
}

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

	phoneNumber := strings.TrimSpace(os.Getenv(TelegramPhoneNumber))
	if phoneNumber == "" {
		return TelegramConfig{}, fmt.Errorf("env %s is required", TelegramPhoneNumber)
	}

	password := strings.TrimSpace(os.Getenv(TelegramPassword))
	if password == "" {
		return TelegramConfig{}, fmt.Errorf("env %s is required", TelegramPassword)
	}

	return TelegramConfig{
		AppID:   appID,
		AppHash: appHash,
		PhoneNumber: phoneNumber,
		Password: password,
	}, nil
}
