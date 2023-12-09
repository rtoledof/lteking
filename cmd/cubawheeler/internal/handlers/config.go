package handlers

import (
	"os"
	"strconv"
)

type Config struct {
	Host  string
	Port  int64
	Path  string
	Redis string
	Mongo string

	SMTPServer   string
	SMTPPort     int64
	SMTPUSer     string
	SMTPPassword string
}

func LoadConfig() Config {
	cfg := Config{
		Port:  3000,
		Path:  "./",
		Redis: "redis://localhost:6379",
		Mongo: "mongodb://localhost:27017/",

		SMTPServer:   "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUSer:     "",
		SMTPPassword: "",
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

	if mongoServer := os.Getenv("MONGO_URL"); len(mongoServer) > 0 {
		cfg.Mongo = mongoServer
	}

	if server, exist := os.LookupEnv("SMTP_SERVER"); exist {
		cfg.SMTPServer = server
	}

	if strPort, exist := os.LookupEnv("SMTP_PORT"); exist {
		port, err := strconv.ParseInt(strPort, 10, 8)
		if err == nil {
			cfg.SMTPPort = port
		}
	}

	if user, exist := os.LookupEnv("SMTP_USER"); exist {
		cfg.SMTPUSer = user
	}

	if pass, exist := os.LookupEnv("SMTP_PASS"); exist {
		cfg.SMTPPassword = pass
	}

	return cfg
}
