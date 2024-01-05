package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mock"
)

func TestLoginHandlerLogin(t *testing.T) {
	ctx := cannon.NewContextWithLogger(context.Background(), slog.Default())
	type fields struct {
		User        cubawheeler.UserService
		OTP         cubawheeler.OtpService
		Application cubawheeler.ApplicationService
		Token       cubawheeler.TokenVerifier
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
			name: "test login handler",
			fields: fields{
				User: &mock.UserService{
					LoginFn: func(ctx context.Context, req cubawheeler.LoginRequest) (*cubawheeler.User, error) {
						return &cubawheeler.User{
							ID:   "123",
							Role: cubawheeler.RoleDriver,
						}, nil
					},
				},
				OTP: &mock.OtpService{
					OtpFn: func(ctx context.Context, otp, email string) error {
						return nil
					},
				},
				Token: &mock.TokenVerifier{
					RemoveByAccessFn: func(ctx context.Context, token string) error {
						return nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					value := url.Values{
						"username": []string{"test@email.com"},
						"password": []string{"123456"},
					}
					req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(value.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "test login handler with missing username",
			fields: fields{
				User: &mock.UserService{
					LoginFn: func(ctx context.Context, req cubawheeler.LoginRequest) (*cubawheeler.User, error) {
						return &cubawheeler.User{
							ID:   "123",
							Role: cubawheeler.RoleDriver,
						}, nil
					},
				},
				OTP: &mock.OtpService{
					OtpFn: func(ctx context.Context, otp, email string) error {
						return nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					value := url.Values{
						"otp": []string{},
					}
					req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(value.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LoginHandler{
				User:        tt.fields.User,
				OTP:         tt.fields.OTP,
				Application: tt.fields.Application,
			}
			req := tt.args.r().WithContext(ctx)
			err := h.Login(tt.args.w, req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoginHandler.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if errors.Is(err, &cubawheeler.Error{}) {
				e := &cubawheeler.Error{}
				errors.As(err, &e)
				if e.StatusCode != tt.wantStatusCode {
					t.Fatalf("LoginHandler.Login() = %v, want %v", tt.args.w.Code, http.StatusUnauthorized)
				}
			} else if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("LoginHandler.Login() = %v, want %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}
