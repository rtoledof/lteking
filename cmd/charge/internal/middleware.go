package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	_ "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"cubawheeler.io/pkg/cubawheeler"
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

// AuthMiddleware decodes the share session cookie and packs the session into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := cubawheeler.GetClaimsFromContext(r.Context())
		if claims == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
		next.ServeHTTP(w, r.WithContext(cubawheeler.NewContextWithUser(r.Context(), &user)))
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
