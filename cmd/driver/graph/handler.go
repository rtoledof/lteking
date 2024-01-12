package graph

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cubawheeler.io/cmd/driver/graph/model"
	"cubawheeler.io/pkg/cubawheeler"
	"github.com/99designs/gqlgen/graphql/handler"
)

func NewHandler(
	orderService string,
	authService string,
) *handler.Server {
	resolver := &Resolver{
		orderService: orderService,
		authService:  authService,
	}
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}

func makeRequest(ctx context.Context, method string, url string, body url.Values) (*http.Response, error) {
	jwtToken := cubawheeler.JWTFromContext(ctx)
	if jwtToken == "" && !strings.Contains(url, "authorize") {
		return nil, fmt.Errorf("no token found: %w", cubawheeler.ErrAccessDenied)
	}
	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(body.Encode())
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v: %w", err, cubawheeler.ErrInternal)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v: %w", err, cubawheeler.ErrInternal)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, cubawheeler.ErrAccessDenied
		}
	}
	return resp, nil
}

func assambleUser(user *cubawheeler.User) *model.ProfileOutput {
	return &model.ProfileOutput{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Profile.Name,
		LastName: user.Profile.LastName,
		Phone:    user.Profile.Phone,
		Dob:      user.Profile.DOB,
		Photo:    user.Profile.Photo,
		Rate:     user.Rate,
		Status:   user.Profile.Status.String(),
		Gender:   user.Profile.Gender.String(),
	}
}
