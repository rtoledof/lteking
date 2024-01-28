package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"

	"cubawheeler.io/cmd/internal"
	"cubawheeler.io/cmd/rider/internal/handlers"
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
		appName = "cubawheeler-rider"
	}

	logger := internal.NewAppLogger(appName, lev)
	slog.SetDefault(logger)

	app := handlers.New(handlers.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println(err)
	}
}