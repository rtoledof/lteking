package internal

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"identity.io/pkg/cannon"
	"identity.io/pkg/identity"
)

type fn func(w http.ResponseWriter, r *http.Request) error

func handler(f fn) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		err := f(writer, r)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			content := err.Error()
			var internalError = &identity.Error{}
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

func ContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
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

func ClientAuthenticate(service identity.ClientService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				next.ServeHTTP(w, r)
				return
			}
			token := requestToken(r)
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			if strings.HasPrefix(token, "sk_") || strings.HasPrefix(token, "pk_") {
				token, err := stripBearerPrefixFromToken(token)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				client, err := service.FindByKey(r.Context(), token)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				ctx = identity.NewContextWithClient(r.Context(), client)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserMidleware(service identity.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/v1/login" || r.URL.Path == "/v1/otp" {
				next.ServeHTTP(w, r)
				return
			}
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			id, ok := claims["id"].(string)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := service.FindByID(r.Context(), id)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := identity.NewContextWithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TokenAuthMiddleware(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := requestToken(r)
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !strings.HasPrefix(token, "sk_") &&
				!strings.HasPrefix(token, "pk_") {
				jwtauth.Verifier(tokenAuth)(next).ServeHTTP(w, r)
				jwtauth.Authenticator(next).ServeHTTP(w, r)
			}

			token, err := stripBearerPrefixFromToken(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := identity.NewContextWithToken(r.Context(), token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var replacer = strings.NewReplacer("sk_", "", "pk_", "", "test_", "")

func stripBearerPrefixFromToken(token string) (string, error) {
	const prefix = "Bearer "
	token = strings.TrimPrefix(token, prefix)
	return replacer.Replace(token), nil
}
