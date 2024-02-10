package internal

import (
	"context"
	"fmt"
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

	"order.io/graph"
	"order.io/pkg/mongo"
	rdb "order.io/pkg/redis"
	"order.io/pkg/seed"
)

type App struct {
	router    http.Handler
	rdb       *rdb.Redis
	mongo     *mongo.DB
	config    Config
	dialer    *gomail.Dialer
	done      chan struct{}
	tokenAuth *jwtauth.JWTAuth
}

func New(cfg Config) *App {
	opt, _ := redis.ParseURL(cfg.Redis)
	client := redis.NewClient(opt)
	redisDB := rdb.NewRedis(client)

	app := &App{
		rdb:       redisDB,
		config:    cfg,
		mongo:     mongo.NewDB(cfg.DB.ConnectionString(), cfg.DB.Database),
		done:      make(chan struct{}),
		tokenAuth: jwtauth.New("HS256", []byte(cfg.JWTPrivateKey), nil),
	}

	app.loader()
	if v := os.Getenv("SEED"); len(v) > 0 {
		if v == "true" {
			seed.RegisterSeeder("rate", func() seed.Seeder { return seed.NewRate(app.mongo) })
			seed.RegisterSeeder("vehicle_category_rate", func() seed.Seeder { return seed.NewVehicleCategoryRate(app.mongo) })
			if err := seed.Up(); err != nil {
				fmt.Println("failed to seed rate", err)
			}
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
	router.Use(jwtauth.Verifier(a.tokenAuth))
	router.Use(TokenAuthMiddleware(a.tokenAuth))
	router.Mount("/debug", middleware.Profiler())

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to driver api"))
	})

	router.Group(func(r chi.Router) {
		grapgqlSrv := graph.NewHandler(
			mongo.NewOrderService(a.mongo, a.rdb),
		)

		r.Handle("/", playground.Handler("Order playground", "/query"))
		r.Handle("/query", grapgqlSrv)
	})

	a.router = router
}
