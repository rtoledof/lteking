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

func TestOtpHandlerOtp(t *testing.T) {
	ctx := context.Background()
	ctx = cannon.NewContextWithLogger(ctx, slog.Default())
	type fields struct {
		OTP  cubawheeler.OtpService
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
			name: "test otp handler",
			fields: fields{
				OTP: &mock.OtpService{
					CreateFn: func(ctx context.Context, email string) (string, error) {
						return "123456", nil
					},
				},
				User: &mock.UserService{
					FindByEmailFn: func(ctx context.Context, email string) (*cubawheeler.User, error) {
						return &cubawheeler.User{
							ID:    "1",
							Email: "",
							Role:  cubawheeler.RoleRider,
							Profile: cubawheeler.Profile{
								Status: cubawheeler.ProfileStatusIncompleted,
							},
							Status: cubawheeler.UserStatusOnReview,
						}, nil
					},
					CreateUserFn: func(ctx context.Context, user *cubawheeler.User) error {
						return nil
					},
					UpdateFn: func(ctx context.Context, user *cubawheeler.User) error {
						return nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/otp", strings.NewReader("email=client&client_secret=secret"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantStatusCode: 200,
		},
		{
			name: "test otp handler with invalid email",
			fields: fields{
				OTP: &mock.OtpService{
					CreateFn: func(ctx context.Context, email string) (string, error) {
						return "123456", nil
					},
				},
				User: &mock.UserService{
					FindByEmailFn: func(ctx context.Context, email string) (*cubawheeler.User, error) {
						return nil, fmt.Errorf("error finding user")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/otp", strings.NewReader("email=client&client_secret=secret"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:        true,
			wantStatusCode: 404,
		},
		{
			name: "test otp handler with not found client",
			fields: fields{
				OTP: &mock.OtpService{
					CreateFn: func(ctx context.Context, email string) (string, error) {
						return "123456", nil
					},
				},
				User: &mock.UserService{
					FindByEmailFn: func(ctx context.Context, email string) (*cubawheeler.User, error) {
						return nil, cubawheeler.ErrNotFound
					},
					CreateUserFn: func(ctx context.Context, user *cubawheeler.User) error {
						return nil
					},
					UpdateFn: func(ctx context.Context, user *cubawheeler.User) error {
						return nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/otp", strings.NewReader("email=client&client_secret=secret"))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:        false,
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &OtpHandler{
				OTP:  tt.fields.OTP,
				User: tt.fields.User,
			}
			req := tt.args.r().WithContext(ctx)
			if err := h.Otp(tt.args.w, req); (err != nil) != tt.wantErr {
				t.Errorf("Otp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantStatusCode != tt.args.w.Code {
				t.Errorf("Otp() wantStatusCode = %v, got %v", tt.wantStatusCode, tt.args.w.Code)
			}
		})
	}
}
