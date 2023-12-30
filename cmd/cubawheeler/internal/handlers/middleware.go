package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	_ "time"

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
			w.WriteHeader(http.StatusUnauthorized)
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

// ClientMiddleware decodes the share session cookie and packs the session into context
func ClientMiddleware(srv cubawheeler.ApplicationService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, pass, ok := r.BasicAuth()
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			app, err := srv.FindByClient(r.Context(), user)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if app.Secret != pass {
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(cubawheeler.NewContextWithClient(r.Context(), app)))
		})
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
