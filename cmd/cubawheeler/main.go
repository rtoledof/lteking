package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-faker/faker/v4"
	_ "github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"

	"cubawheeler.io/cmd/cubawheeler/internal/handlers"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/database"
)

func init() {
	godotenv.Load()
	loadDatabase()
}

func main() {
	app := handlers.New(handlers.LoadConfig(database.Db))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println(err)
	}
}

func loadDatabase() {
	database.InitDb()
	database.Db.AutoMigrate(&cubawheeler.Client{})
	database.Db.AutoMigrate(&cubawheeler.User{})
	database.Db.AutoMigrate(&cubawheeler.Profile{})
	database.Db.AutoMigrate(&cubawheeler.Ads{})
	database.Db.AutoMigrate(&cubawheeler.Plan{})
	database.Db.AutoMigrate(&cubawheeler.Vehicle{})
	database.Db.AutoMigrate(&cubawheeler.Location{})
	database.Db.AutoMigrate(&cubawheeler.Trip{})

	seedData()
}

func seedData() {
	var clients = []cubawheeler.Client{}
	for i := 0; i < 5; i++ {
		var client cubawheeler.Client
		faker.FakeData(&client)
		clients = append(clients, client)
	}
	var users = []cubawheeler.User{}
	var profiles []cubawheeler.Profile
	for i := 0; i < 10; i++ {
		var user cubawheeler.User
		user.ID = cubawheeler.NewID().String()
		faker.FakeData(&user)
		faker.FakeData(&user.Profile)
		user.Profile.UserID = user.ID
		profiles = append(profiles, user.Profile)
		users = append(users, user)
	}

	database.Db.Save(clients)
	database.Db.Save(users)
	database.Db.Save(profiles)
}
