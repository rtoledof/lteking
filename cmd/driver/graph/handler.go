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
) *handler.Server {
	resolver := &Resolver{
		orderService: orderService,
	}
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
}

func makeRequest(ctx context.Context, method string, url string, body url.Values) (*http.Response, error) {
	jwtToken := cubawheeler.JWTFromContext(ctx)
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

func assambleUser(user *cubawheeler.User) *model.User {
	return &model.User{
		ID:      user.ID,
		Email:   user.Email,
		Role:    model.Role(user.Role.String()),
		Profile: assambleProfile(&user.Profile),
	}
}

func assambleProfile(profile *cubawheeler.Profile) *model.Profile {
	gender := model.Gender(profile.Gender.String())
	return &model.Profile{
		ID:               profile.ID,
		Name:             &profile.Name,
		LastName:         &profile.LastName,
		Phone:            profile.Phone,
		Dob:              &profile.DOB,
		Photo:            profile.Photo,
		Dni:              &profile.Dni,
		Licence:          &profile.Licence,
		Gender:           &gender,
		PreferedCurrency: &profile.PreferedCurrency,
	}
}
