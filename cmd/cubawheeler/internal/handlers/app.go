package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
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

	abl "cubawheeler.io/pkg/ably"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/graph"
	"cubawheeler.io/pkg/mailer"
	"cubawheeler.io/pkg/mongo"
	"cubawheeler.io/pkg/pusher"
	rdb "cubawheeler.io/pkg/redis"
	"cubawheeler.io/pkg/seed"
)

var tokenAuth *jwtauth.JWTAuth
var privateKey string

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)
}

type App struct {
	router       http.Handler
	rdb          *rdb.Redis
	mongo        *mongo.DB
	pmConfig     cubawheeler.PaymentmethodConfig
	config       Config
	pusher       *pusher.Pusher
	notification *pusher.PushNotification
	seed         seed.Seeder
	dialer       *gomail.Dialer
	done         chan struct{}
	orderChan    chan *cubawheeler.Order
	realTime     cubawheeler.RealTimeService
	rest         *ably.REST
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
		pusher: pusher.NewPusher(
			cfg.PusherAppId,
			cfg.PusherKey,
			cfg.PusherSecret,
			cfg.PusherCluster,
			cfg.PusherSecure,
		),
		notification: pusher.NewPushNotification(cfg.BeansInterest, cfg.BeansSecret),
		done:         make(chan struct{}),
		orderChan:    make(chan *cubawheeler.Order, 10000),
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

	pushNotification := pusher.NewPushNotification(a.config.BeansInterest, a.config.BeansSecret)
	userSrv := mongo.NewUserService(
		a.mongo,
		rdb.NewBeansToken(a.rdb, pushNotification),
		pushNotification,
		a.done,
	)
	appSrv := mongo.NewApplicationService(a.mongo)
	client := abl.NewClient(a.config.Amqp.Connection, a.done, a.config.Ably.ApiKey)

	go client.Consumer.Consume(
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
	router.Use(AuthMiddleware(userSrv))
	router.Use(ClientMiddleware(appSrv))

	router.Post("/pusher/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params, _ := ioutil.ReadAll(r.Body)
		user := cubawheeler.UserFromContext(r.Context())
		if user == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		auth, err := a.pusher.Authenticate(params, user, bytes.Contains(params, []byte("presence")))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Info("an error ocurred authenticating the user in a channel: %v", err)
			return
		}
		fmt.Printf("Request: %s\n", params)
		fmt.Printf("Response: %s\n", auth)
		w.WriteHeader(http.StatusOK)
		w.Write(auth)
	})

	router.Get("/pusher/beans-auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		usr := cubawheeler.UserFromContext(r.Context())
		if usr == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		beansToken, err := pushNotification.GenerateToken(usr.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(beansToken)
	})

	router.Post("/pusher/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		user := chi.URLParam(r, "id")
		err := pushNotification.PublishToUser([]string{user}, pusher.Notification{
			Title:    "Test",
			Body:     "Test message",
			Metadata: map[string]string{"aditional": "no se puede"},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	router.Group(func(r chi.Router) {
		grapgqlSrv := graph.NewHandler(
			a.rdb,
			a.mongo,
			userSrv,
			a.done,
			a.config.Amqp.Connection,
			a.config.Ably.ApiKey,
			a.pmConfig,
		)
		r.Handle("/", playground.Handler("GraphQL playground", "/query"))
		r.Handle("/query", grapgqlSrv)
	})

	a.router = router
}
