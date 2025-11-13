package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	app, err := initServ(ctx)
	if err != nil {
		log.Fatalf("init: %v", err)
	}
	defer app.Logger.Sync()

	app.Logger.Info("application started")
	menties, err := app.SheetsService.GetMentyInformation(ctx)
	if err != nil {
		app.Logger.Errorw("failed to get menty information", "error", err)
		return
	}
	app.Logger.Infow("menty information retrieved", "count", len(menties))
	app.Logger.Info("application stopped")
}
