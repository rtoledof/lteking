package mongo

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	database = os.Getenv("MONGO_DB_NAME")
}

var database string

type Collections string

func (c Collections) String() string {
	return string(c)
}

type DB struct {
	client   *mongo.Client
	database string
}

func NewDB(serverURL, database string) *DB {
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

	return &DB{client: client, database: database}
}

func (db *DB) Collection(name Collections) *mongo.Collection {
	return db.client.Database(db.database).Collection(name.String())
}

func (db *DB) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}
