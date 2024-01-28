package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"wallet.io/internal"
)

func init() {
	godotenv.Load()
}

func main() {

	lev := slog.LevelInfo
	if os.Getenv("debug") == "debug" {
		lev = slog.LevelDebug
	}
	appName := os.Getenv("app_name")
	if appName == "" {
		appName = "cubawheeler-wallet"
	}

	logger := internal.NewAppLogger(appName, lev)
	slog.SetDefault(logger)

	app := internal.New(internal.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println(err)
	}
}
