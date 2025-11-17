package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/finlleyl/menty-ping/internal/config"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

type (
	Client struct {
		*telegram.Client
	}
)

const (
	TelegramSession = "TELEGRAM_SESSION"
)

func getSessionPath() (string, error) {
	path := os.Getenv(TelegramSession)

	if path == "" {
		return "", fmt.Errorf("env %s is required", TelegramSession)
	}
	return path, nil
}

func NewTelegramClient(config config.TelegramConfig) (*Client, error) {
	path, err := getSessionPath()
	if err != nil {
		return nil, fmt.Errorf("get session path: %w", err)
	}
	client := telegram.NewClient(
		config.AppID,
		config.AppHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: path},
		})

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) WithTelegram(
	ctx context.Context,
	cfg config.TelegramConfig,
	fn func(ctx context.Context) error,
) error {
	return c.Client.Run(ctx, func(ctx context.Context) error {
		if err := c.authTelegramClient(ctx, cfg); err != nil {
			return fmt.Errorf("auth telegram client: %w", err)
		}
		return fn(ctx)
	})
}

func (c *Client) authTelegramClient(ctx context.Context, config config.TelegramConfig) error {
	status, err := c.Auth().Status(ctx)
	if err != nil {
		return fmt.Errorf("get auth status: %w", err)
	}
	if !status.Authorized {
		flow := auth.NewFlow(
			auth.Constant(
				config.PhoneNumber,
				config.Password,
				auth.CodeAuthenticatorFunc(getCodeFunc),
			),
			auth.SendCodeOptions{
				AllowFlashCall: false,
				CurrentNumber:  true,
			},
		)

		if err := flow.Run(ctx, c.Auth()); err != nil {
			return fmt.Errorf("auth flow: %w", err)
		}
	}
	return nil
}

func getCodeFunc(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read code: %w", err)
	}
	return strings.TrimSpace(code), nil
}
