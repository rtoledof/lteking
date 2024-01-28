package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
	_ "time"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

func ContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type fn func(w http.ResponseWriter, r *http.Request) error

func handler(f fn) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		err := f(writer, r)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
	}
}

// / AuthMiddleware decodes the share session cookie and packs the session into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := cubawheeler.GetClaimsFromContext(r.Context())
		if claims == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		if len(claims) > 0 {
			userData, ok := claims["user"]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			var user cubawheeler.User
			if err := json.Unmarshal([]byte(userData), &user); err != nil {
				slog.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx = cubawheeler.NewContextWithUser(ctx, &user)
		}

		token := requestToken(r)
		if token != "" {
			ctx = cubawheeler.NewContextWithJWT(ctx, token)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// / AuthMiddleware decodes the share session cookie and packs the session into context
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := requestToken(r)
		if token != "" {
			ctx = cubawheeler.NewContextWithJWT(ctx, token)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ClientMiddleware decodes the share session cookie and packs the session into context
func ClientMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims := cubawheeler.GetClaimsFromContext(r.Context())
		if claims == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		if len(claims) > 0 {
			clientData, ok := claims["client"]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			var client cubawheeler.Application
			if err := json.Unmarshal([]byte(clientData), &client); err != nil {
				slog.Error(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx = cubawheeler.NewContextWithClient(ctx, &client)
		}

		token := requestToken(r)
		if token != "" {
			ctx = cubawheeler.NewContextWithJWT(ctx, token)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestToken(r *http.Request) string {
	const prefix = "Bearer "

	token, _, ok := r.BasicAuth()
	if !ok {
		h := r.Header.Get("Authorization")
		if strings.HasPrefix(h, prefix) {
			token = strings.TrimPrefix(h, prefix)
		}
	}
	return token
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (rw *statusResponseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	if !rw.headerWritten {
		rw.statusCode = status
		rw.headerWritten = true
	}
}

func (rw *statusResponseWriter) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func newStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{w, http.StatusOK, false}
}

func CanonicalLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const requestKey = "request_id"
		const headerRequestID = "X-Request-ID"
		requestID := r.Header.Get(headerRequestID)
		logger := slog.Default()
		if len(requestID) > 0 {
			logger = logger.With(requestKey, requestID)
		}
		l, reset := cannon.NewLogger(logger)
		defer reset()
		logger = l.Logger()
		r = r.WithContext(cannon.NewContextWithLogger(r.Context(), logger))
		start := time.Now()
		rw := newStatusResponseWriter(w)
		logger.Info("Request started",
			slog.Group("http",
				slog.String("method", r.Method),
				slog.String("client_ip", r.RemoteAddr),
				slog.String("path", r.URL.Path),
				slog.String("user_agent", r.Header.Get("User-Agent")),
			),
		)
		next.ServeHTTP(rw, r)
		l.Emit(
			slog.Group("http",
				slog.Int("status", rw.statusCode),
			),
			slog.String("duration", time.Since(start).String()),
		)
	})
}