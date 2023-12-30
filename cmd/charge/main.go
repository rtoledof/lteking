package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	_ "github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"

	"cubawheeler.io/cmd/order/internal"
)

func init() {
	godotenv.Load()
}

func main() {
	app := internal.New(internal.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println(err)
	}
}
