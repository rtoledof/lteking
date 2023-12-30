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

	"cubawheeler.io/cmd/auth/internal/handlers"
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
	router      http.Handler
	rdb         *rdb.Redis
	redisClient *redis.Client
	mongo       *mongo.DB
	config      Config
	seed        seed.Seeder
	dialer      *gomail.Dialer
	done        chan struct{}
}

func New(cfg Config) *App {
	opt, _ := redis.ParseURL(cfg.Redis)
	client := redis.NewClient(opt)
	client.Options().DB = cfg.RedisDB
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
		app.seed = seed.NewSeed(app.mongo)
		if err := app.seed.Up(); err != nil {
			fmt.Println("unable to upload seeds")
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

	userSrv := mongo.NewUserService(a.mongo, nil, nil, a.done)
	appSrv := mongo.NewApplicationService(a.mongo)
	// tokenStore := oauth.NewTokenStore(a.redisClient)
	// client := abl.NewClient(a.config.Amqp.Connection, a.done, a.config.Ably.ApiKey)

	// manager := manage.NewDefaultManager()
	// manager.MustTokenStorage(token, nil)

	// clientStore := mongo.NewApplicationService(a.mongo)
	// manager.MapClientStorage(clientStore)

	// srv := server.NewDefaultServer(manager)
	// srv.SetAllowGetAccessRequest(true)
	// srv.SetClientInfoHandler(server.ClientFormHandler)

	// srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
	// 	fmt.Println("Internal Error:", err.Error())
	// 	return
	// })

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(AuthMiddleware(userSrv))
	router.Use(ClientMiddleware(appSrv))

	tokenVerifier := rdb.NewTokenVerifier(a.rdb, userSrv, appSrv)

	s := oauth.NewBearerServer(
		os.Getenv("JWT_PRIVATE_KEY"),
		time.Hour*24*30,
		tokenVerifier,
		nil,
	)

	{
		h := &handlers.LoginHandler{
			User:        userSrv,
			OTP:         rdb.NewOtpService(a.rdb),
			Application: appSrv,
		}

		router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			err := h.Login(w, r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			switch r.FormValue("grant_type") {
			case "password":
				s.UserCredentials(w, r)
			case "client_credentials":
				s.ClientCredentials(w, r)
			}
		})
	}

	{
		h := &handlers.OtpHandler{
			OTP:  rdb.NewOtpService(a.rdb),
			User: userSrv,
		}

		router.Post("/otp", handler(h.Otp))
	}

	{
		h := handlers.NewStatusHandler(userSrv)

		router.Post("/status", handler(h.Availability))
	}

	{
		h := handlers.NewCarHandler(userSrv)

		router.Post("/car", handler(h.Car))
	}

	{
		h := handlers.NewProfileHandler(userSrv)

		router.Put("/profile", handler(h.Update))
		router.Get("/me", handler(h.Get))
	}

	{
		h := handlers.NewVehicleHandler(userSrv)

		router.Route("/v1/vehicles", func(r chi.Router) {
			r.Post("/", handler(h.Add))
			r.Put("/{id}", handler(h.Update))
			r.Delete("/{id}", handler(h.Remove))
			r.Get("/", handler(h.List))
			r.Post("/{id}", handler(h.SetActiveVehicle))
		})
	}

	a.router = router
}
