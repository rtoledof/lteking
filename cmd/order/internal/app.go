package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/oauth"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"

	"cubawheeler.io/cmd/order/internal/handlers"
	"cubawheeler.io/pkg/mailer"
	"cubawheeler.io/pkg/mongo"
	rdb "cubawheeler.io/pkg/redis"
	"cubawheeler.io/pkg/seed"
)

var tokenAuth *jwtauth.JWTAuth
var privateKey string

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)
}

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
	if s := os.Getenv("SEED"); len(s) > 0 {
		seed.RegisterSeeder("rate", func() seed.Seeder { return seed.NewRate(app.mongo) })
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

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(CanonicalLog)
	router.Mount("/debug", middleware.Profiler())

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to driver api"))
	})

	router.Group(func(r chi.Router) {
		{
			h := handlers.OrderHandler{
				Service: mongo.NewOrderService(a.mongo, nil, a.rdb),
			}
			r.Use(TokenMiddleware)
			r.Use(oauth.Authorize(a.config.JWTPrivateKey, nil))
			r.Use(AuthMiddleware)

			r.Route("/v1/orders", func(r chi.Router) {
				r.Post("/", handler(h.Create))
				r.Get("/", handler(h.List))

				r.Get("/{id}", handler(h.FindByID))
				r.Put("/{id}", handler(h.Update))
				r.Post("/{id}/start", handler(h.Start))
				r.Post("/{id}/accept", handler(h.Accept))
				r.Post("/{id}/cancel", handler(h.Cancel))
				r.Post("/{id}/complete", handler(h.Complete))
				r.Post("/{id}/confirm", handler(h.Confirm))
			})

		}
	})

	a.router = router
}
