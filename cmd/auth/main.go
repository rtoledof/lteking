package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"

	"cubawheeler.io/cmd/auth/internal"
	mi "cubawheeler.io/cmd/internal"
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
		appName = "cubawheeler-order"
	}

	logger := mi.NewAppLogger(appName, lev)
	slog.SetDefault(logger)

	app := internal.New(internal.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println(err)
	}
}
