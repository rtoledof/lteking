package mapbox

import (
	"context"
	"fmt"
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

type DirectionService service

func (s *DirectionService) GetRoute(ctx context.Context, request cubawheeler.DirectionRequest) (_ *cubawheeler.DirectionResponse, _ string, err error) {
	defer derrors.Wrap(&err, "mapbox.DirectionService.GetRoute")
	if request.Valid() {
		return nil, "", cubawheeler.ErrInvalidInput
	}
	path := fmt.Sprintf("%s/directions/v5/mapbox/driving/%s", s.client.BaseURL, request.String())
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}
	q := req.URL.Query()
	q.Add("access_token", s.client.AccessToken)
	q.Add("steps", "true")
	q.Add("language", "es")
	req.URL.RawQuery = q.Encode()
	var response cubawheeler.DirectionResponse
	_, strBody, err := s.client.Do(req, &response)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get response: %w", err)
	}
	return &response, strBody, nil
}
