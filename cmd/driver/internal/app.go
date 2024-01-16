package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ably/ably-go/ably"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"

	"cubawheeler.io/cmd/driver/graph"
	abl "cubawheeler.io/pkg/ably"
	"cubawheeler.io/pkg/mailer"
	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/realtime"
	rdb "cubawheeler.io/pkg/redis"
	"cubawheeler.io/pkg/seed"
)

var tokenAuth *jwtauth.JWTAuth
var privateKey string

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)
}

type App struct {
	router             http.Handler
	rdb                *rdb.Redis
	mongo              *mongo.DB
	config             Config
	seed               seed.Seeder
	dialer             *gomail.Dialer
	done               chan struct{}
	realTime           *realtime.RealTimeService
	ablyRealTimeClient *ably.Realtime
	ablClient          *abl.Client
}

func New(cfg Config) *App {
	opt, _ := redis.ParseURL(cfg.Redis)
	client := redis.NewClient(opt)
	redisDB := rdb.NewRedis(client)
	mongoDB := mongo.NewDB(cfg.Mongo)

	user := mongo.NewUserService(
		mongoDB,
		cfg.ServiceDiscovery.WalletService,
		make(chan struct{}),
	)
	ablyRealTimeClient, err := ably.NewRealtime(ably.WithKey(cfg.Ably.ApiKey))
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	transport := abl.AuthTransport{
		Token: cfg.Ably.ApiKey,
	}
	ablClient := abl.NewClient(
		cfg.Amqp.Connection,
		done,
		cfg.Ably.ApiKey,
		ablyRealTimeClient,
		transport.Client(),
	)
	app := &App{
		rdb:    redisDB,
		config: cfg,
		mongo:  mongo.NewDB(cfg.Mongo),
		dialer: gomail.NewDialer(
			cfg.SMTPServer,
			int(cfg.SMTPPort),
			cfg.SMTPUSer,
			cfg.SMTPPassword,
		),

		done: done,

		ablClient: ablClient,

		ablyRealTimeClient: ablyRealTimeClient,

		realTime: realtime.NewRealTimeService(
			rdb.NewRealTimeService(redisDB),
			ablClient.Notifier,
			user,
			redisDB,
			mongo.NewOrderService(mongoDB, nil, redisDB, ablClient.Notifier),
			ablyRealTimeClient,
		),
	}

	app.loader()
	if s := os.Getenv("SEED"); len(s) > 0 {
		seed.RegisterSeeder("vehicle_categories", func() seed.Seeder { return seed.NewVehicleCategoryRate(app.mongo) })
		if err := seed.Up(); err != nil {
			fmt.Printf("unable to upload seeds: %s\n", err.Error())
		}
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server on", addr)

	ch := make(chan error, 1)
	go func() {
		err = httpSrv.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}

		close(ch)
	}()

	err = <-ch

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		a.done <- struct{}{}
		defer cancel()
		return httpSrv.Shutdown(timeout)
	}
}

func (a *App) loader() {

	mailer.NewMailer(
		a.config.SMTPServer,
		a.config.SMTPUSer,
		a.config.SMTPPassword,
		int(a.config.SMTPPort),
	)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST"},
		AllowedHeaders: []string{
			"User-Agent",
			"Content-Type",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
			"Connection",
			"Authorization",
			"DNT",
			"Host",
			"Origin",
			"Pragma",
			"Referer",
			"X-API-KEY",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		Debug:            true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	go a.ablClient.Consumer.Consume(
		a.config.Amqp.Queue,
		a.config.Amqp.Consumer,
		a.config.Amqp.AutoAsk,
		a.config.Amqp.Exclusive,
		a.config.Amqp.NoLocal,
		a.config.Amqp.NoWait,
		a.config.Amqp.Arg,
	)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(CanonicalLog)
	router.Use(TokenMiddleware)
	router.Use(ClientMiddleware)
	router.Use(AuthMiddleware)

	router.Group(func(r chi.Router) {

		srv := graph.NewHandler(a.config.ServiceDiscovery.OrderService, a.config.ServiceDiscovery.AuthService)

		r.Handle("/", playground.Handler("GraphQL playground", "/query"))
		r.Handle("/query", srv)
	})

	a.router = router
}
