package internal

import (
	"fmt"
	"os"
	"strconv"
)

type DB struct {
	Host     string
	Port     int64
	Username string
	Password string
	Database string
}

func (db DB) ConnectionString() string {
	connectionString := fmt.Sprintf("mongodb://%s:%d/?retryWrites=true&w=majority", db.Host, db.Port)
	if len(db.Username) > 0 {
		connectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d/?retryWrites=true&w=majority", db.Username, db.Password, db.Host, db.Port)
	}
	return connectionString
}

func (db DB) WithUsernamePassword(username, password string) DB {
	db.Username = username
	db.Password = password
	return db
}

func (db DB) WithDatabase(database string) DB {
	db.Database = database
	return db
}

type Config struct {
	Host          string
	Port          int64
	DB            DB
	MongoDatabase string

	JWTPrivateKey string
}

func LoadConfig() Config {
	cfg := Config{
		Port: 3000,
		DB: DB{
			Host:     "localhost",
			Port:     27017,
			Database: "wallet",
		},
	}

	if serverHost, exist := os.LookupEnv("SERVER_HOST"); exist {
		cfg.Host = serverHost
	}
	if serverPort, exist := os.LookupEnv("SERVER_PORT"); exist {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.Port = int64(port)
		}
	}

	if mongoHost, exist := os.LookupEnv("MONGO_HOST"); exist {
		cfg.DB.Host = mongoHost
	}

	if mongoPort, exist := os.LookupEnv("MONGO_PORT"); exist {
		if port, err := strconv.ParseUint(mongoPort, 10, 16); err == nil {
			cfg.DB.Port = int64(port)
		}
	}

	if mongoDBName := os.Getenv("MONGO_DB_NAME"); len(mongoDBName) > 0 {
		cfg.MongoDatabase = mongoDBName
	}

	if mongoUser := os.Getenv("MONGO_USER"); len(mongoUser) > 0 {
		cfg.DB.Username = mongoUser
	}

	if mongoPass := os.Getenv("MONGO_PASS"); len(mongoPass) > 0 {
		cfg.DB.Password = mongoPass
	}

	if key, exist := os.LookupEnv("JWT_SECRET_KEY"); exist {
		cfg.JWTPrivateKey = key
	}

	return cfg
}
