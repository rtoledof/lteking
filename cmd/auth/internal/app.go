package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"

	"cubawheeler.io/cmd/auth/internal/handlers"
	"cubawheeler.io/pkg/mailer"
	"cubawheeler.io/pkg/mongo"
	rdb "cubawheeler.io/pkg/redis"
	"cubawheeler.io/pkg/seed"
)

type App struct {
	router      http.Handler
	rdb         *rdb.Redis
	redisClient *redis.Client
	mongo       *mongo.DB
	config      Config
	dialer      *gomail.Dialer
	done        chan struct{}
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
		mongo:       mongo.NewDB(cfg.Mongo),
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
		seed.RegisterSeeder("application", func() seed.Seeder { return seed.NewApplication(app.mongo) })
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

	userSrv := mongo.NewUserService(a.mongo, a.done)
	appSrv := mongo.NewApplicationService(a.mongo)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(CanonicalLog)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(TokenMiddleware)

	tokenVerifier := rdb.NewTokenVerifier(a.rdb, userSrv, appSrv)

	s := oauth.NewBearerServer(
		a.config.JWTPrivateKey,
		time.Hour*24*30,
		tokenVerifier,
		nil,
	)

	{
		h := handlers.NewAuthorizeHandler(appSrv)

		router.Post("/authorize", func(w http.ResponseWriter, r *http.Request) {
			err := h.Authorize(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			s.ClientCredentials(w, r)
		})
	}

	router.Group(func(r chi.Router) {

		r.Use(oauth.Authorize(a.config.JWTPrivateKey, nil))
		r.Use(AuthMiddleware)
		r.Use(ClientMiddleware(appSrv))
		r.Use(TokenMiddleware)

		{
			h := &handlers.OtpHandler{
				OTP:  rdb.NewOtpService(a.rdb),
				User: userSrv,
			}

			router.Post("/otp", handler(h.Otp))
		}

		{
			h := &handlers.LoginHandler{
				User:        userSrv,
				OTP:         rdb.NewOtpService(a.rdb),
				Application: appSrv,
			}

			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				err := h.Login(w, r)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				s.UserCredentials(w, r)
			})
		}

		{
			h := handlers.NewStatusHandler(userSrv)

			r.Post("/status", handler(h.Availability))
		}

		{
			h := handlers.NewCarHandler(userSrv)

			r.Post("/car", handler(h.SetActiveVehicle))
		}

		{
			h := handlers.NewProfileHandler(userSrv)

			r.Put("/profile", handler(h.Update))
			r.Get("/me", handler(h.Get))
			r.Post("/profile/devices", handler(h.AddDevice))
		}

		{
			h := handlers.NewVehicleHandler(userSrv)

			r.Route("/v1/vehicles", func(r chi.Router) {
				r.Post("/", handler(h.Add))
				r.Put("/{id}", handler(h.Update))
				r.Delete("/{id}", handler(h.Remove))
				r.Get("/", handler(h.List))
				r.Post("/{id}", handler(h.SetActiveVehicle))
			})
		}

	})

	a.router = router
}
