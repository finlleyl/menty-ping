package main

import (
	"context"
	"fmt"
	"log"

	"github.com/finlleyl/menty-ping/internal/config"
	"github.com/finlleyl/menty-ping/internal/logger"
	sheetsclient "github.com/finlleyl/menty-ping/internal/sheets/client"
	sheetsservice "github.com/finlleyl/menty-ping/internal/sheets/service"
	telegramclient "github.com/finlleyl/menty-ping/internal/telegram/client"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type App struct {
	Logger         *zap.SugaredLogger
	SheetsService  *sheetsservice.Service
	TelegramClient *telegramclient.Client
}

func initServ(ctx context.Context) (*App, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("godotenv: skipped loading .env: %v", err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Errorw("failed to load config", "error", err)
		return nil, fmt.Errorf("init config: %w", err)
	}
	logger.Infow("configuration loaded")

	oauthCfg, store := cfg.Google.OAuthConfig, cfg.Google.TokenStore
	client := sheetsclient.New(oauthCfg, store)
	logger.Infow("sheets oauth client initialized")

	httpClient, err := client.HTTPClient(ctx)
	if err != nil {
		logger.Errorw("failed to create sheets http client", "error", err)
		return nil, fmt.Errorf("init sheets client: %w", err)
	}
	logger.Infow("sheets http client ready")

	sheetsService, err := sheetsservice.NewSheetsService(ctx, httpClient, logger)
	if err != nil {
		logger.Errorw("failed to initialize sheets service", "error", err)
		return nil, fmt.Errorf("init sheets service: %w", err)
	}
	logger.Infow("sheets service initialized")

	telegramClient, err := telegramclient.NewTelegramClient(cfg.Telegram)
	if err != nil {
		logger.Errorw("failed to initialize telegram client", "error", err)
		return nil, fmt.Errorf("init telegram client: %w", err)
	}
	logger.Infow("telegram client initialized")

	return &App{
		Logger:         logger,
		SheetsService:  sheetsService,
		TelegramClient: telegramClient,
	}, nil
}
