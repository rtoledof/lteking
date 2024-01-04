package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mock"
)

func TestProfileHandlerUpdate(t *testing.T) {
	ctx := cannon.NewContextWithLogger(context.Background(), slog.Default())
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID: "123456789",
		Profile: cubawheeler.Profile{
			Phone:    "123456789",
			LastName: "Doe",
			DOB:      "01/01/2000",
			Photo:    "https://example.com/photo.jpg",
		},
	})
	type fields struct {
		User cubawheeler.UserService
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
			name: "update profile",
			fields: fields{
				User: &mock.UserService{
					UpdateFn: func(ctx context.Context, user *cubawheeler.User) error {
						return nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodPost, "/profile", nil)
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					return r
				},
			},
			wantErr:        false,
			wantStatusCode: 204,
		},
		{
			name: "update unexisting profile",
			fields: fields{
				User: &mock.UserService{
					UpdateFn: func(ctx context.Context, user *cubawheeler.User) error {
						return cubawheeler.NewNotFound("user")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodPost, "/profile", nil)
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					return r
				},
			},
			wantErr:        true,
			wantStatusCode: 404,
		},
		{
			name: "update unexisting profile",
			fields: fields{
				User: &mock.UserService{
					UpdateFn: func(ctx context.Context, user *cubawheeler.User) error {
						return cubawheeler.ErrNotExist
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodPost, "/profile", nil)
					r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					return r
				},
			},
			wantErr:        true,
			wantStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewProfileHandler(tt.fields.User)
			req := tt.args.r().WithContext(ctx)
			err := h.Update(tt.args.w, req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileHandler.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			e := &cubawheeler.Error{}
			if err != nil {
				if errors.As(err, &e) {
					if e.StatusCode != tt.wantStatusCode {
						t.Errorf("ProfileHandler.Update() status code = %v, wantStatusCode %v", e.StatusCode, tt.wantStatusCode)
					}
				} else if tt.wantStatusCode != tt.args.w.Code {
					t.Errorf("ProfileHandler.Update() status code = %v, wantStatusCode %v", tt.args.w.Code, tt.wantStatusCode)
				}
			}
		})
	}
}
