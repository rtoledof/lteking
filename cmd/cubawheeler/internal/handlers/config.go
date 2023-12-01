package handlers

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/gorm"

	"cubawheeler.io/pkg/bolt"
)

type Config struct {
	Host  string
	Port  int64
	Path  string
	Redis string
	DB    *gorm.DB
}

func LoadConfig(db *gorm.DB) Config {
	cfg := Config{
		Port:  3000,
		Path:  "./",
		Redis: "redis://localhost:6379",
		DB:    db,
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
	db, err := bolt.Open(fmt.Sprintf("%s/cubawheeler.db", cfg.Path))
	if err != nil {
		panic(err)
	}
	cfg.DB = db

	return cfg
}
