package internal

import (
	"fmt"
	"os"
	"strconv"
)

type ServiceDiscover struct {
	AuthServiceURL string
}

type DB struct {
	Host     string
	Port     int64
	User     string
	Pass     string
	Options  string
	Database string
}

func (db DB) ConnectionString() string {
	connectionString := fmt.Sprintf("mongodb://%s:%d/?retryWrites=true&w=majority", db.Host, db.Port)
	if len(db.User) > 0 {
		connectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d/?retryWrites=true&w=majority", db.User, db.Pass, db.Host, db.Port)
	}
	return connectionString
}

type Config struct {
	Host  string
	Port  int64
	Path  string
	Redis string
	DB    DB

	JWTPrivateKey string
}

func LoadConfig() Config {
	cfg := Config{
		Port:  3000,
		Path:  "./",
		Redis: "redis://localhost:6379",
		DB: DB{
			Host:     "localhost",
			Port:     27017,
			Database: "orders",
			Options:  "retryWrites=true&w=majority",
		},
	}

	if redisAddr, exist := os.LookupEnv("REDIS_ADDR"); exist {
		cfg.Redis = redisAddr
	}
	if serverHost, exist := os.LookupEnv("SERVER_HOST"); exist {
		cfg.Host = serverHost
	}
	if serverPort, exist := os.LookupEnv("SERVER_PORT"); exist {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.Port = int64(port)
		}
	}

	if path, exist := os.LookupEnv("DB_PATH"); exist {
		cfg.Path = path
	}

	if mongoServer := os.Getenv("MONGO_HOST"); len(mongoServer) > 0 {
		cfg.DB.Host = mongoServer
	}

	if port := os.Getenv("MONGO_PORT"); len(port) > 0 {
		var err error
		cfg.DB.Port, err = strconv.ParseInt(port, 10, 16)
		if err != nil {
			cfg.DB.Port = 27017
		}
	}

	if mongoDB := os.Getenv("MONGO_DB_NAME"); len(mongoDB) > 0 {
		cfg.DB.Database = mongoDB
	}

	if mongoUser := os.Getenv("MONGO_USER"); len(mongoUser) > 0 {
		cfg.DB.User = mongoUser
	}

	if mongoPass := os.Getenv("MONGO_PASS"); len(mongoPass) > 0 {
		cfg.DB.Pass = mongoPass
	}

	if key, exist := os.LookupEnv("JWT_SECRET_KEY"); exist {
		cfg.JWTPrivateKey = key
	}
	if cfg.JWTPrivateKey == "" {
		panic("JWT_SECRET_KEY is not set")
	}

	return cfg
}
