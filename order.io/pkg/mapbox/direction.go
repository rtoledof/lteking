package mapbox

import (
	"context"
	"fmt"
	"net/http"

	"order.io/pkg/derrors"
	"order.io/pkg/order"
)

type DirectionService service

func (s *DirectionService) GetRoute(ctx context.Context, request order.DirectionRequest) (_ *order.DirectionResponse, _ string, err error) {
	defer derrors.Wrap(&err, "mapbox.DirectionService.GetRoute")
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
	var response order.DirectionResponse
	_, strBody, err := s.client.Do(req, &response)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get response: %w", err)
	}
	return &response, strBody, nil
}
