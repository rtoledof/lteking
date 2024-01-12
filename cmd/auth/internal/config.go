package internal

import (
	"os"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Amqp struct {
	Connection string
	Queue      string
	Consumer   string
	AutoAsk    bool
	Exclusive  bool
	NoLocal    bool
	NoWait     bool
	Arg        amqp.Table
}

type Ably struct {
	ApiKey           string
	ApiSubscriperKey string
}

type Config struct {
	Host    string
	Port    int64
	Path    string
	Redis   string
	RedisDB int
	Mongo   string
	MongoDB string
	Amqp    Amqp

	JWTPrivateKey string

	SMTPServer   string
	SMTPPort     int64
	SMTPUSer     string
	SMTPPassword string

	PusherAppId   string
	PusherKey     string
	PusherSecret  string
	PusherCluster string
	PusherSecure  bool

	BeansInterest string
	BeansSecret   string

	Ably Ably

	WalletApi string
}

func LoadConfig() Config {
	cfg := Config{
		Port:  3000,
		Path:  "./",
		Redis: "redis://localhost:6379",
		Mongo: "mongodb://localhost:27017/",

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

	if path, exist := os.LookupEnv("DB_PATH"); exist {
		cfg.Path = path
	}

	if mongoServer := os.Getenv("MONGO_URL"); len(mongoServer) > 0 {
		cfg.Mongo = mongoServer
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

	if appId := os.Getenv("PUSHER_APP_ID"); len(appId) > 0 {
		cfg.PusherAppId = appId
	}

	if key := os.Getenv("PUSHER_Key"); len(key) > 0 {
		cfg.PusherKey = key
	}

	if secret := os.Getenv("PUSHER_SECRET"); len(secret) > 0 {
		cfg.PusherSecret = secret
	}

	if cluster := os.Getenv("PUSHER_CLUSTER"); len(cluster) > 0 {
		cfg.PusherCluster = cluster
	}

	if s := os.Getenv("PUSHER_SECURE"); len(s) > 0 {
		secure, err := strconv.ParseBool(s)
		if err == nil {
			cfg.PusherSecure = secure
		}
	}

	if instance := os.Getenv("BEANS_INSTANCE"); len(instance) > 0 {
		cfg.BeansInterest = instance
	}

	if secret := os.Getenv("BEANS_SECRET"); len(secret) > 0 {
		cfg.BeansSecret = secret
	}

	if amqp := os.Getenv("AMQP_CONNECTION"); len(amqp) > 0 {
		cfg.Amqp.Connection = amqp
	}

	if queue, exist := os.LookupEnv("AMQP_QUEUE"); exist {
		cfg.Amqp.Queue = queue
	}

	if consumer, exist := os.LookupEnv("AMQP_CONSUMER"); exist {
		cfg.Amqp.Consumer = consumer
	}

	if apiKey, exist := os.LookupEnv("ABLY_API_KEY"); exist {
		cfg.Ably.ApiKey = apiKey
	}

	if subscriberApiKey, exist := os.LookupEnv("ABLY_SUBSCRIBER_API_KEY"); exist {
		cfg.Ably.ApiSubscriperKey = subscriberApiKey
	}

	if privateKey, exist := os.LookupEnv("JWT_SECRET_KEY"); exist {
		cfg.JWTPrivateKey = privateKey
	}
	if cfg.JWTPrivateKey == "" {
		panic("JWT_SECRET_KEY is not set")
	}

	if walletApi, exist := os.LookupEnv("WALLET_API"); exist {
		cfg.WalletApi = walletApi
	}
	if cfg.WalletApi == "" {
		panic("WALLET_API is not set")
	}
	return cfg
}
