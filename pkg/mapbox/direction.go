package mapbox

import (
	"fmt"
	"net/http"
	"net/url"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

type DirectionService service

func (s *DirectionService) GetRoute(request cubawheeler.DirectionRequest) (_ *cubawheeler.DirectionResponse, err error) {
	defer derrors.Wrap(&err, "mapbox.DirectionService.GetRoute")
	if request.Valid() {
		return nil, cubawheeler.ErrInvalidInput
	}
	path := fmt.Sprintf("/directions/v5/mapbox/driving/%s", request.String())
	var url = url.Values{
		"access_token": []string{s.client.AccessToken},
		"steps":        []string{"true"},
		"language":     []string{"es"},
	}
	path = path + "?" + url.Encode()

	req, err := s.client.NewRequest(http.MethodGet, path, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	var response cubawheeler.DirectionResponse
	_, err = s.client.Do(req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	return &response, nil
}
