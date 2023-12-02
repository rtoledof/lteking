package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func NewDB(serverURL string) *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(serverURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic("unable to connect with mongo server")
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic("database server not available")
	}

	return &DB{client: client}
}
