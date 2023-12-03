package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	_ "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

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

// Middleware decodes the share session cookie and packs the session into context
func AuthMiddleware(srv cubawheeler.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, err := parseToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			claim, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}

			user, err := srv.FindByEmail(r.Context(), claim["jti"].(string))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(cubawheeler.NewContextWithUser(r.Context(), user)))
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
