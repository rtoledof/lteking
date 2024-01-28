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
	return fmt.Sprintf("mongodb://%s:%d/?retryWrites=true&w=majority", db.Host, db.Port)
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
	Host    string
	Port    int64
	Redis   string
	RedisDB int
	Mongo   DB
	MongoDB string

	JWTPrivateKey string

	SMTPServer   string
	SMTPPort     int64
	SMTPUSer     string
	SMTPPassword string

	WalletApi string
}

func DefaultConfig() Config {
	return Config{
		Port:  3000,
		Redis: "redis://localhost:6379",
		Mongo: DB{
			Host:     "localhost",
			Port:     27017,
			Database: "identity",
		},
	}
}

func LoadConfig() Config {
	cfg := Config{
		Port:  3000,
		Redis: "redis://localhost:6379",
		Mongo: DB{
			Host:     "localhost",
			Port:     27017,
			Database: "identity",
		},

		SMTPServer: "smtp.gmail.com",
		SMTPPort:   587,
	}

	if redisAddr, exist := os.LookupEnv("REDIS_ADDR"); exist {
		cfg.Redis = redisAddr
	}
	if redisDB, exist := os.LookupEnv("REDIS_DB"); exist {
		var err error
		cfg.RedisDB, err = strconv.Atoi(redisDB)
		if err != nil {
			cfg.RedisDB = 0
		}
	}
	if serverHost, exist := os.LookupEnv("SERVER_HOST"); exist {
		cfg.Host = serverHost
	}
	if serverPort, exist := os.LookupEnv("SERVER_PORT"); exist {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.Port = int64(port)
		}
	}

	if mongoServer := os.Getenv("MONGO_HOST"); len(mongoServer) > 0 {
		cfg.Mongo.Host = mongoServer
	}

	if mongoPort := os.Getenv("MONGO_PORT"); len(mongoPort) > 0 {
		var err error
		cfg.Mongo.Port, err = strconv.ParseInt(mongoPort, 10, 16)
		if err != nil {
			cfg.Mongo.Port = 27017
		}
	}

	if database := os.Getenv("MONGO_DB_NAME"); len(database) > 0 {
		cfg.MongoDB = database
	}

	if server := os.Getenv("SMTP_SERVER"); len(server) > 0 {
		cfg.SMTPServer = server
	}

	if strPort := os.Getenv("SMTP_PORT"); len(strPort) > 0 {
		port, err := strconv.ParseInt(strPort, 10, 8)
		if err == nil {
			cfg.SMTPPort = port
		}
	}

	if user := os.Getenv("SMTP_USER"); len(user) > 0 {
		cfg.SMTPUSer = user
	}

	if pass := os.Getenv("SMTP_PASS"); len(pass) > 0 {
		cfg.SMTPPassword = pass
	}

	if privateKey, exist := os.LookupEnv("JWT_SECRET"); exist {
		cfg.JWTPrivateKey = privateKey
	}
	if cfg.JWTPrivateKey == "" {
		panic("JWT_SECRET is not set")
	}

	if walletApi, exist := os.LookupEnv("WALLET_API"); exist {
		cfg.WalletApi = walletApi
	}

	if cfg.WalletApi == "" {
		panic("WALLET_API is not set")
	}
	return cfg
}
