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

func TestCarHandlerAdd(t *testing.T) {
	ctx := cannon.NewContextWithLogger(context.Background(), slog.Default())
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID:   "123",
		Role: cubawheeler.RoleDriver,
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
			name: "test car handler",
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
					value := url.Values{
						"plate":    []string{"123456"},
						"name":     []string{"test"},
						"brand":    []string{string(cubawheeler.BrandBmw.String())},
						"model":    []string{"test"},
						"year":     []string{"2020"},
						"color":    []string{"red"},
						"type":     []string{string(cubawheeler.TypeAuto)},
						"seats":    []string{"4"},
						"category": []string{string(cubawheeler.VehicleCategoryX)},
					}
					r := httptest.NewRequest(http.MethodPost, "/car", strings.NewReader(value.Encode()))
					r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					r = r.WithContext(ctx)
					return r
				},
			},
			wantErr:        false,
			wantStatusCode: http.StatusNoContent,
		},
		{
			name: "test car handler add with bad brand",
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
					value := url.Values{
						"plate": []string{"123456"},
						"name":  []string{"test"},
						"brand": []string{"test"},
						"model": []string{"test"},
						"year":  []string{"2020"},
						"color": []string{"red"},
						"type":  []string{string(cubawheeler.TypeAuto)},
						"seats": []string{"4"},
					}
					r := httptest.NewRequest(http.MethodPost, "/car", strings.NewReader(value.Encode()))
					r = r.WithContext(ctx)
					r = r.WithContext(ctx)
					return r
				},
			},
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCarHandler(tt.fields.User)
			req := tt.args.r().WithContext(ctx)
			err := h.Add(tt.args.w, req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CarHandler.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if errors.Is(err, &cubawheeler.Error{}) {
				if err.(*cubawheeler.Error).StatusCode != tt.wantStatusCode {
					t.Fatalf("CarHandler.Add() error = %v, wantErr %v", err, tt.wantStatusCode)
				}
			} else if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("CarHandler.Add() error = %v, wantErr %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestCardhandlerSetActiveVehicle(t *testing.T) {
	ctx := cannon.NewContextWithLogger(context.Background(), slog.Default())
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID:   "123",
		Role: cubawheeler.RoleDriver,
		Vehicles: []*cubawheeler.Vehicle{
			{
				ID: "123",
			},
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
			name: "test car handler",
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
					value := url.Values{
						"car": []string{"123"},
					}
					r := httptest.NewRequest(http.MethodPost, "/car", strings.NewReader(value.Encode()))
					r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					r = r.WithContext(ctx)
					return r
				},
			},
			wantErr:        false,
			wantStatusCode: http.StatusNoContent,
		},
		{
			name: "test car handler add with not found car",
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
					value := url.Values{
						"car": []string{"1234"},
					}
					r := httptest.NewRequest(http.MethodPost, "/car", strings.NewReader(value.Encode()))
					r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					r = r.WithContext(ctx)
					return r
				},
			},
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCarHandler(tt.fields.User)
			req := tt.args.r().WithContext(ctx)
			err := h.SetActiveVehicle(tt.args.w, req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CarHandler.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if errors.Is(err, &cubawheeler.Error{}) {
				if err.(*cubawheeler.Error).StatusCode != tt.wantStatusCode {
					t.Fatalf("CarHandler.Add() error = %v, wantErr %v", err, tt.wantStatusCode)
				}
			} else if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("CarHandler.Add() error = %v, wantErr %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}
