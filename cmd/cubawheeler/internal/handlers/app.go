package handlers

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
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"

	"cubawheeler.io/pkg/database"
	"cubawheeler.io/pkg/graph"
	"cubawheeler.io/pkg/oauth"
)

var tokenAuth *jwtauth.JWTAuth
var privateKey string

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)
}

type App struct {
	router http.Handler
	rdb    *redis.Client
	config Config
}

func New(cfg Config) *App {
	opt, _ := redis.ParseURL(cfg.Redis)
	client := redis.NewClient(opt)
	app := &App{
		rdb:    client,
		config: cfg,
	}

	app.loader()

	return app
}

func (a *App) Start(ctx context.Context) error {
	//fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
	addr := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
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
		defer cancel()
		return httpSrv.Shutdown(timeout)
	}

	return nil
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

	router.Group(func(r chi.Router) {
		//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
		tokenStorer := oauth.NewTokenStore(a.rdb)
		userSrv := database.NewUserService(database.Db)
		r.Use(AuthMiddleware(userSrv))
		grapgqlSrv := graph.NewHandler(tokenStorer, userSrv)
		r.Handle("/", playground.Handler("GraphQL playground", "/query"))
		r.Handle("/query", grapgqlSrv)
	})

	a.router = router
}
