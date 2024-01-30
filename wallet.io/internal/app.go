package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"

	"wallet.io/graph"
	"wallet.io/pkg/mongo"
)

type App struct {
	router  http.Handler
	mongo   *mongo.DB
	config  Config
	jwtAuth *jwtauth.JWTAuth
	done    chan struct{}
}

func New(cfg Config) *App {
	app := &App{
		config:  cfg,
		mongo:   mongo.NewDB(cfg.DB.ConnectionString(), cfg.DB.Database),
		jwtAuth: jwtauth.New("HS256", []byte(cfg.JWTPrivateKey), nil),
		done:    make(chan struct{}),
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

	fmt.Println("Starting server on", addr)

	ch := make(chan error, 1)
	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}

		close(ch)
	}()

	select {
	case err := <-ch:
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
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
	router.Use(jwtauth.Verifier(a.jwtAuth))
	router.Use(TokenAuthMiddleware(a.jwtAuth))
	router.Use(middleware.Heartbeat("/ping"))

	router.Mount("/debug", middleware.Profiler())

	router.Group(func(r chi.Router) {
		grapgqlSrv := graph.NewHandler(
			mongo.NewWalletService(a.mongo),
		)

		r.Handle("/", playground.Handler("Wallet playground", "/query"))
		r.Handle("/query", grapgqlSrv)
	})

	a.router = router
}
