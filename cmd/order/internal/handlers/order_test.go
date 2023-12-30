package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cubawheeler.io/cmd/driver/graph/model"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mock"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
)

type fn func(w http.ResponseWriter, r *http.Request) error

func handler(f fn) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		err := f(writer, r)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
		}
	}
}

func TestOrderHandlerCreate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{
		Role: cubawheeler.RoleAdmin,
	})
	type fields struct {
		service cubawheeler.OrderService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		StatusCode int
		want       *model.Order
	}{
		{
			name: "TestOrderHandlerCreate",
			fields: fields{
				service: &mock.OrderService{
					CreateFunc: func(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
						return &cubawheeler.Order{
							ID: "1",
							Items: cubawheeler.OrderItem{
								Points: []*cubawheeler.Point{
									{Lat: 37.772, Lng: -122.214},
									{Lat: 21.291, Lng: -157.821},
									{Lat: -18.142, Lng: 178.431},
									{Lat: -27.467, Lng: 153.027},
								},
								Baggages: true,
								Riders:   1,
								Currency: "CUP",
							},
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					data, err := os.ReadFile("testdata/order_ok.json")
					if err != nil {
						t.Fatal(err)
					}
					req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(data))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
			},
			wantErr:    false,
			StatusCode: http.StatusOK,
			want: &model.Order{
				ID: "1",
				Item: &model.Item{
					Points: []*model.Point{
						{Lat: 37.772, Lng: -122.214},
						{Lat: 21.291, Lng: -157.821},
						{Lat: -18.142, Lng: 178.431},
						{Lat: -27.467, Lng: 153.027},
					},
					Baggages: func() *bool { result := true; return &result }(),
					Riders:   func() *int { result := 1; return &result }(),
					Currency: func() *string { result := "CUP"; return &result }(),
				},
			},
		},
		{
			name: "TestOrderHandlerCreateInvalidInput",
			fields: fields{
				service: &mock.OrderService{
					CreateFunc: func(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
						return nil, cubawheeler.ErrInvalidInput
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					data, err := os.ReadFile("testdata/order_error.json")
					if err != nil {
						t.Fatal(err)
					}
					req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(data))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
			},
			wantErr:    true,
			StatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderHandler{
				Service: tt.fields.service,
			}
			tt.args.r = tt.args.r.WithContext(ctx)
			if err := o.Create(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Fatalf("OrderHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.StatusCode != tt.args.w.(*httptest.ResponseRecorder).Code {
				t.Fatalf("OrderHandler.Create() StatusCode = %v, want %v", tt.args.w.(*httptest.ResponseRecorder).Code, tt.StatusCode)
			}
			if tt.want != nil {
				var got model.Order
				if err := json.NewDecoder(tt.args.w.(*httptest.ResponseRecorder).Body).Decode(&got); err != nil {
					t.Fatalf("OrderHandler.Create() error = %v", err)
				}
				if diff := cmp.Diff(tt.want, &got); diff != "" {
					t.Fatalf("OrderHandler.Create() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestOrderHandlerUpdate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{
		Role: cubawheeler.RoleAdmin,
	})
	type fields struct {
		service cubawheeler.OrderService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		StatusCode int
		want       *model.Order
	}{
		{
			name: "TestOrderHandlerUpdate",
			fields: fields{
				service: &mock.OrderService{
					UpdateFunc: func(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
						return &cubawheeler.Order{
							ID: "1",
							Items: cubawheeler.OrderItem{
								Points: []*cubawheeler.Point{
									{Lat: 37.772, Lng: -122.214},
									{Lat: 21.291, Lng: -157.821},
									{Lat: -18.142, Lng: 178.431},
									{Lat: -27.467, Lng: 153.027},
								},
								Baggages: true,
								Riders:   req.Riders,
								Currency: "CUP",
							},
						}, nil
					},
					FindByIDFunc: func(ctx context.Context, id string) (*cubawheeler.Order, error) {
						return &cubawheeler.Order{
							ID: "1",
							Items: cubawheeler.OrderItem{
								Points: []*cubawheeler.Point{
									{Lat: 37.772, Lng: -122.214},
									{Lat: 21.291, Lng: -157.821},
									{Lat: -18.142, Lng: 178.431},
									{Lat: -27.467, Lng: 153.027},
								},
								Baggages: true,
								Riders:   1,
								Currency: "CUP",
							},
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					data, err := os.ReadFile("testdata/order_update_ok.json")
					if err != nil {
						t.Fatal(err)
					}
					req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewReader(data))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
			},
			wantErr:    false,
			StatusCode: http.StatusOK,
			want: &model.Order{
				ID: "1",
				Item: &model.Item{
					Points: []*model.Point{
						{Lat: 37.772, Lng: -122.214},
						{Lat: 21.291, Lng: -157.821},
						{Lat: -18.142, Lng: 178.431},
						{Lat: -27.467, Lng: 153.027},
					},
					Baggages: func() *bool { result := true; return &result }(),
					Riders:   func() *int { result := 2; return &result }(),
					Currency: func() *string { result := "CUP"; return &result }(),
				},
			},
		},
		{
			name: "TestOrderHandlerUpdateInvalidInput",
			fields: fields{
				service: &mock.OrderService{
					UpdateFunc: func(ctx context.Context, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
						return nil, cubawheeler.ErrInvalidInput
					},
					FindByIDFunc: func(ctx context.Context, id string) (*cubawheeler.Order, error) {
						return &cubawheeler.Order{
							ID: "1",
							Items: cubawheeler.OrderItem{
								Points: []*cubawheeler.Point{
									{Lat: 37.772, Lng: -122.214},
									{Lat: 21.291, Lng: -157.821},
									{Lat: -18.142, Lng: 178.431},
									{Lat: -27.467, Lng: 153.027},
								},
								Baggages: true,
								Riders:   1,
								Currency: "CUP",
							},
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					data, err := os.ReadFile("testdata/order_error.json")
					if err != nil {
						t.Fatal(err)
					}
					req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewReader(data))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
			},
			wantErr:    true,
			StatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderHandler{
				Service: tt.fields.service,
			}
			router := chi.NewRouter()
			router.Put("/orders/{id}", handler(o.Update))
			tt.args.r = tt.args.r.WithContext(ctx)
			if err := o.Update(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Fatalf("OrderHandler.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.StatusCode != tt.args.w.(*httptest.ResponseRecorder).Code {
				t.Fatalf("OrderHandler.Update() StatusCode = %v, want %v", tt.args.w.(*httptest.ResponseRecorder).Code, tt.StatusCode)
			}
			if tt.want != nil {
				var got model.Order
				if err := json.NewDecoder(tt.args.w.(*httptest.ResponseRecorder).Body).Decode(&got); err != nil {
					t.Fatalf("OrderHandler.Update() error = %v", err)
				}
				if diff := cmp.Diff(tt.want, &got); diff != "" {
					t.Fatalf("OrderHandler.Update() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
