package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"

	"cubawheeler.io/cmd/wallet/internal/handlers"
	"cubawheeler.io/pkg/mailer"
	"cubawheeler.io/pkg/mongo"
	rdb "cubawheeler.io/pkg/redis"
)

type App struct {
	router http.Handler
	rdb    *rdb.Redis
	mongo  *mongo.DB
	config Config
	dialer *gomail.Dialer
	done   chan struct{}
}

func New(cfg Config) *App {
	opt, _ := redis.ParseURL(cfg.Redis)
	client := redis.NewClient(opt)
	redisDB := rdb.NewRedis(client)
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
		done: make(chan struct{}),
	}

	app.loader()

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

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(CanonicalLog)

	router.Mount("/debug", middleware.Profiler())

	router.Group(func(r chi.Router) {
		{
			h := handlers.NewWalletHandler(mongo.NewWalletService(a.mongo))

			r.Use(oauth.Authorize(a.config.JWTPrivateKey, nil))
			r.Use(TokenMiddleware)
			r.Use(AuthMiddleware)
			r.Use(ClientMiddleware)

			r.Route("/v1/wallet", func(r chi.Router) {
				r.Get("/", handler(h.Balance))
				r.Post("/", handler(h.Create))
				r.Get("/transactions", handler(h.Transactions))
				r.Post("/transfer", handler(h.Transfer))
				r.Post("/topup", handler(h.TopUp))
			})
		}
	})

	a.router = router
}
