package internal

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
	_ "time"

	"order.io/pkg/cannon"
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
