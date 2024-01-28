package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/go-chi/jwtauth"

	"wallet.io/pkg/cannon"
)

type fn func(w http.ResponseWriter, r *http.Request) error

func handler(f fn) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		err := f(writer, r)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
		}
	}
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

func TokenAuthMiddleware(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path != "/health" && r.URL.Path != "/query" {
				jwtauth.Verifier(tokenAuth)(next).ServeHTTP(w, r)
				jwtauth.Authenticator(next).ServeHTTP(w, r)
			}

			next.ServeHTTP(w, r)
		})
	}
}
