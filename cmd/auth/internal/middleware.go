package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"cubawheeler.io/cmd/auth/internal/handlers"
	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

type fn func(w http.ResponseWriter, r *http.Request) error

func handler(f fn) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		err := f(writer, r)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			content := err.Error()
			var internalError = &cubawheeler.Error{}
			if errors.As(err, &internalError) {
				writer.WriteHeader(internalError.StatusCode)
				data, err := json.Marshal(internalError)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
				content = string(data)
			}
			writer.Write([]byte(content))
		}
	}
}

// AuthMiddleware decodes the share session cookie and packs the session into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/authorize" {
			next.ServeHTTP(w, r)
			return
		}

		claims := cubawheeler.GetClaimsFromContext(r.Context())
		if claims == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		userData, ok := claims["user"]
		if ok {
			var user cubawheeler.User
			if err := json.Unmarshal([]byte(userData), &user); err != nil {
				next.ServeHTTP(w, r)
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

// ClientMiddleware decodes the share session cookie and packs the session into context
func ClientMiddleware(srv cubawheeler.ApplicationService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := cubawheeler.GetClaimsFromContext(r.Context())
			if claims == nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()

			claimClient, ok := claims["client"]
			if ok {
				var client cubawheeler.Application
				if err := json.Unmarshal([]byte(claimClient), &client); err != nil {
					next.ServeHTTP(w, r)
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
}

// / TokenMiddleware decodes the share session cookie and packs the session into context
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := requestToken(r)
		if token != "" {
			ctx = handlers.NewContextWithToken(ctx, token)
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

var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromToken,
}

var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func parseToken(r *http.Request) (*jwt.Token, error) {
	jwtToken, err := request.ParseFromRequest(r, authHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(os.Getenv("JWT_SECRET"))
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("inavalid token provided: %w", err)
	}
	return jwtToken, nil
}

var replacer = strings.NewReplacer("sk_", "", "pk_", "", "test_", "")

func stripBearerPrefixFromToken(token string) (string, error) {
	const prefix = "Bearer "
	if strings.HasPrefix(token, prefix) {
		token = strings.TrimPrefix(token, prefix)
	}
	return token, nil
}
