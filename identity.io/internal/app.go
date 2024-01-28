package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"

	"identity.io/graph"
	"identity.io/pkg/identity"
	"identity.io/pkg/mailer"
	"identity.io/pkg/mongo"
	rdb "identity.io/pkg/redis"
	"identity.io/pkg/seed"
)

type App struct {
	router      http.Handler
	rdb         *rdb.Redis
	redisClient *redis.Client
	mongo       *mongo.DB
	config      Config
	dialer      *gomail.Dialer
	client      identity.ClientService
	done        chan struct{}
	tokenAuth   *jwtauth.JWTAuth
}

func New(cfg Config) *App {
	opt, _ := redis.ParseURL(cfg.Redis)
	client := redis.NewClient(opt)
	if cfg.RedisDB > 0 {
		client.Options().DB = cfg.RedisDB
	}

	redisDB := rdb.NewRedis(client)
	app := &App{
		rdb:         redisDB,
		redisClient: client,
		config:      cfg,
		mongo:       mongo.NewDB(cfg.Mongo.ConnectionString(), cfg.Mongo.Database),
		dialer: gomail.NewDialer(
			cfg.SMTPServer,
			int(cfg.SMTPPort),
			cfg.SMTPUSer,
			cfg.SMTPPassword,
		),
		done:      make(chan struct{}),
		tokenAuth: jwtauth.New("HS256", []byte(cfg.JWTPrivateKey), nil),
	}

	app.client = mongo.NewClientService(app.mongo)
	app.loader()
	// TODO: check how to load the seeds
	if s := os.Getenv("SEED"); len(s) > 0 {
		seed.RegisterSeeder("client", func() seed.Seeder { return seed.NewClient(app.mongo) })
		if err := seed.Up(); err != nil {
			slog.Info(err.Error())
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

	if err := a.mongo.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to mongo: %w", err)
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

	userSrv := mongo.NewUserService(a.mongo, a.config.WalletApi, a.done, a.rdb, a.tokenAuth)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(CanonicalLog)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(ClientAuthenticate(a.client))
	router.Use(TokenAuthMiddleware(a.tokenAuth))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := a.mongo.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	router.Group(func(r chi.Router) {
		grapgqlSrv := graph.NewHandler(userSrv, rdb.NewOtpService(a.rdb))

		r.Handle("/", playground.Handler("Identity playground", "/query"))
		r.Handle("/query", grapgqlSrv)
	})

	a.router = router
}
