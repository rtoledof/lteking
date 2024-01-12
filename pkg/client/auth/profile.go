package auth

import (
	"context"
	"net/http"
	"net/url"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

type UpdateProfile struct {
	Name     string
	Dob      string
	LastName string
	Gender   string
	Phone    string
	Photo    string
	License  string
	Dni      string
}

type ProfileService service

func (s *ProfileService) UpdateProfile(ctx context.Context, req UpdateProfile) (*cubawheeler.Profile, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("ProfileService.Profile")
	value := url.Values{
		"name":      []string{req.Name},
		"last_name": []string{req.LastName},
		"phone":     []string{req.Phone},
		"dob":       []string{req.Dob},
		"gender":    []string{req.Gender},
		"photo":     []string{req.Photo},
		"dni":       []string{req.Dni},
		"licence":   []string{req.License},
	}

	profileRequest, err := s.client.NewRequest(http.MethodPut, "/profile", value)
	if err != nil {
		return nil, err
	}
	var profile cubawheeler.Profile
	if _, err := s.client.Do(profileRequest, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *ProfileService) GetProfile(ctx context.Context) (*cubawheeler.User, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("ProfileService.Profile")
	profileRequest, err := s.client.NewRequest(http.MethodGet, "/me", nil)
	if err != nil {
		return nil, err
	}
	var user cubawheeler.User
	if _, err := s.client.Do(profileRequest, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *ProfileService) AddDevice(ctx context.Context, device string) (*cubawheeler.Profile, error) {
	logger := cannon.LoggerFromContext(ctx)
	logger.Info("ProfileService.AddDevice")
	value := url.Values{
		"device": []string{device},
	}
	profileRequest, err := s.client.NewRequest(http.MethodPost, "/profile/devices", value)
	if err != nil {
		return nil, err
	}
	var profile cubawheeler.Profile
	if _, err := s.client.Do(profileRequest, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
