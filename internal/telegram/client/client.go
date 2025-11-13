package client

import (
	"context"

	"github.com/finlleyl/menty-ping/internal/config"
	"github.com/gotd/td/telegram"
)

type (
	Client struct {
		*telegram.Client
	}
)

func NewTelegramClient(ctx context.Context, config config.TelegramConfig) (*Client, error) {
	client := telegram.NewClient(
		config.AppID,
		config.AppHash,
		telegram.Options{})

	return &Client{
		Client: client,
	}, nil
}
