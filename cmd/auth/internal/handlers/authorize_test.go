package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mock"
)

func TestAuthorizeHandlerAuthorize(t *testing.T) {
	ctx := cannon.NewContextWithLogger(context.Background(), slog.Default())
	type fields struct {
		Service cubawheeler.ApplicationService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r func() *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "test authorize handler",
			fields: fields{
				Service: &mock.ApplicationService{
					FindByClientFn: func(ctx context.Context, client string) (*cubawheeler.Application, error) {
						return &cubawheeler.Application{
							ID:     "1",
							Client: "client",
							Secret: "secret",
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/authorize", strings.NewReader("client_id=client&client_secret=secret"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantStatusCode: 200,
		},
		{
			name: "test authorize handler with invalid client",
			fields: fields{
				Service: &mock.ApplicationService{
					FindByClientFn: func(ctx context.Context, client string) (*cubawheeler.Application, error) {
						return nil, fmt.Errorf("error finding application")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/authorize", strings.NewReader("client_id=client&client_secret=secret"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:        true,
			wantStatusCode: 401,
		},
		{
			name: "test authorize handler with invalid client secret",
			fields: fields{
				Service: &mock.ApplicationService{
					FindByClientFn: func(ctx context.Context, client string) (*cubawheeler.Application, error) {
						return &cubawheeler.Application{
							ID:     "1",
							Client: "client",
							Secret: "secret",
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/authorize", strings.NewReader("client_id=client&client_secret=invalid"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:        true,
			wantStatusCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAuthorizeHandler(tt.fields.Service)
			req := tt.args.r().WithContext(ctx)
			err := h.Authorize(tt.args.w, req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AuthorizeHandler.Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantStatusCode != 0 && tt.wantStatusCode != tt.args.w.Code {
				t.Fatalf("AuthorizeHandler.Authorize() statusCode = %v, wantStatusCode %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}
