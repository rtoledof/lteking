package mongo

// import (
// 	"context"
// 	"encoding/base64"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/ably/ably-go/ably"
// 	"github.com/google/go-cmp/cmp"
// 	"github.com/redis/go-redis/v9"

// 	"order.io/pkg/currency"
// 	"order.io/pkg/mock"
// 	"order.io/pkg/order"
// 	"order.io/pkg/realtime"
// 	rdb "order.io/pkg/redis"
// )

// func TestOrderServiceCreate(t *testing.T) {
// 	user := &order.User{ID: "test", Role: order.RoleRider}
// 	ctx := order.NewContextWithUser(context.Background(), user)

// 	database = "test"
// 	db := NewTestDB()
// 	defer func() {
// 		db.client.Database(database).Drop(ctx)
// 		db.client.Disconnect(ctx)
// 	}()

// 	setupVethiclesCategoriesRate(t, ctx, db)
// 	rate := setupRate(t)
// 	if err := storeRate(ctx, db, rate); err != nil {
// 		t.Fatal(err)
// 	}

// 	var tests = []struct {
// 		name    string
// 		request func() *order.DirectionRequest
// 		want    *order.Order
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid",
// 			request: func() *order.DirectionRequest {
// 				return &order.DirectionRequest{
// 					Points: []*order.Point{
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 					},
// 				}
// 			},
// 			want: &order.Order{
// 				Items: order.OrderItem{
// 					Points: []*order.Point{
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 					},
// 					Riders: 1,
// 				},
// 				Status: order.OrderStatusNew,
// 				Rider:  user.ID,
// 				Route: &order.DirectionResponse{
// 					Distance: 100,
// 					Duration: 100,
// 					Routes: []*order.Route{
// 						{
// 							Distance: 100,
// 							Duration: 100,
// 							Geometry: "test",
// 						},
// 					},
// 				},
// 				Currency:    "CUP",
// 				RouteString: base64.StdEncoding.EncodeToString([]byte("test")),
// 				Distance:    100,
// 				Duration:    100,
// 				CategoryPrice: []*order.CategoryPrice{
// 					{
// 						Category: order.VehicleCategoryX,
// 						Price:    int(price(100, 100, *rate, 1)),
// 						Currency: "CUP",
// 					},
// 					{
// 						Category: order.VehicleCategoryXl,
// 						Price:    int(price(100, 100, *rate, 1) * 1.2),
// 						Currency: "CUP",
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "invalid",
// 			request: func() *order.DirectionRequest {
// 				return &order.DirectionRequest{
// 					Points: []*order.Point{
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 					},
// 				}
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	notifier := mock.Notifier{
// 		NotifyRiderOrderAcceptedFn: func(ctx context.Context, devices []string, notification realtime.OrderNotification) error {
// 			return nil
// 		},
// 		NotifyToDevicesFn: func(ctx context.Context, devices []string, notification realtime.OrderNotification, realTime *ably.Realtime, rest *ably.REST) error {
// 			return nil
// 		},
// 	}
// 	charger := mock.Charger{
// 		ChargeFn: func(ctx context.Context, method order.ChargeMethod, amount currency.Amount) (*order.Charge, error) {
// 			return &order.Charge{
// 				Status: order.ChargeStatusSucceeded,
// 			}, nil
// 		},
// 		RefundFn: func(ctx context.Context, payment string, amount currency.Amount) (*order.Charge, error) {
// 			return &order.Charge{}, nil
// 		},
// 	}

// 	client := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})

// 	cache := rdb.NewRedis(client)

// 	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"id": "test", "role": "admin"}`))
// 	}))

// 	direction := mock.Direction{
// 		GetRouteFn: func(ctx context.Context, request order.DirectionRequest) (_ *order.DirectionResponse, _ string, err error) {
// 			return &order.DirectionResponse{
// 				Distance: 100,
// 				Duration: 100,
// 				Routes: []*order.Route{
// 					{
// 						Distance: 100,
// 						Duration: 100,
// 						Geometry: "test",
// 					},
// 				},
// 			}, "test", nil
// 		},
// 	}

// 	server := NewOrderService(
// 		db,
// 		charger,
// 		cache,
// 		notifier,
// 		authServer.URL,
// 	)
// 	server.direction = direction

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := server.Create(ctx, tt.request())
// 			if err != nil && !tt.wantErr {
// 				t.Errorf("OrderService.Create() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.want != nil {
// 				got.ID = tt.want.ID
// 				got.CreatedAt = tt.want.CreatedAt
// 				if diff := cmp.Diff(tt.want, got); diff != "" {
// 					t.Errorf("OrderService.Create() mismatch (-want +got):\n%s", diff)
// 				}
// 			}
// 		})
// 	}
// }

// func TestOrderServiceUpdate(t *testing.T) {
// 	user := &order.User{ID: "test", Role: order.RoleRider}
// 	ctx := order.NewContextWithUser(context.Background(), user)

// 	database = "test"
// 	db := NewTestDB()
// 	defer func() {
// 		db.client.Database(database).Drop(ctx)
// 		db.client.Disconnect(ctx)
// 	}()

// 	categoriesRates := setupVethiclesCategoriesRate(t, ctx, db)
// 	rate := setupRate(t)
// 	if err := storeRate(ctx, db, rate); err != nil {
// 		t.Fatal(err)
// 	}

// 	order := setupOrder(t, db)

// 	var tests = []struct {
// 		name    string
// 		request func() *order.DirectionRequest
// 		want    *order.Order
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid",
// 			request: func() *order.DirectionRequest {
// 				return &order.DirectionRequest{
// 					ID: order.ID,
// 					Points: []*order.Point{
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 					},
// 				}
// 			},
// 			want: &order.Order{
// 				ID:     order.ID,
// 				Items:  order.Items,
// 				Status: order.OrderStatusNew,
// 				Rider:  user.ID,
// 				Route: &order.DirectionResponse{
// 					Distance: 101,
// 					Duration: 101,
// 					Routes: []*order.Route{
// 						{
// 							Distance: 101,
// 							Duration: 101,
// 							Geometry: "test",
// 						},
// 					},
// 				},
// 				Currency:    "CUP",
// 				RouteString: base64.StdEncoding.EncodeToString([]byte("test")),
// 				Distance:    101,
// 				Duration:    101,
// 				CategoryPrice: []*order.CategoryPrice{
// 					{
// 						Category: order.VehicleCategoryX,
// 						Price:    int(price(101, 101, *rate, 1) * categoriesRates[0].Factor),
// 						Currency: "CUP",
// 					},
// 					{
// 						Category: order.VehicleCategoryXl,
// 						Price:    int(price(101, 101, *rate, 1) * categoriesRates[1].Factor),
// 						Currency: "CUP",
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "invalid",
// 			request: func() *order.DirectionRequest {
// 				return &order.DirectionRequest{
// 					Points: []*order.Point{
// 						{
// 							Lat: 23.123,
// 							Lng: 23.123,
// 						},
// 					},
// 				}
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	fmt.Println(int(price(101, 101, *rate, 1) * categoriesRates[0].Factor))
// 	fmt.Println(int(price(101, 101, *rate, 1) * categoriesRates[1].Factor))

// 	notifier := mock.Notifier{
// 		NotifyRiderOrderAcceptedFn: func(ctx context.Context, devices []string, notification realtime.OrderNotification) error {
// 			return nil
// 		},
// 		NotifyToDevicesFn: func(ctx context.Context, devices []string, notification realtime.OrderNotification, realTime *ably.Realtime, rest *ably.REST) error {
// 			return nil
// 		},
// 	}
// 	charger := mock.Charger{
// 		ChargeFn: func(ctx context.Context, method order.ChargeMethod, amount currency.Amount) (*order.Charge, error) {
// 			return &order.Charge{
// 				Status: order.ChargeStatusSucceeded,
// 			}, nil
// 		},
// 		RefundFn: func(ctx context.Context, payment string, amount currency.Amount) (*order.Charge, error) {
// 			return &order.Charge{}, nil
// 		},
// 	}

// 	client := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})

// 	cache := rdb.NewRedis(client)

// 	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"id": "test", "role": "admin"}`))
// 	}))

// 	direction := mock.Direction{
// 		GetRouteFn: func(ctx context.Context, request order.DirectionRequest) (_ *order.DirectionResponse, _ string, err error) {
// 			return &order.DirectionResponse{
// 				Distance: 101,
// 				Duration: 101,
// 				Routes: []*order.Route{
// 					{
// 						Distance: 101,
// 						Duration: 101,
// 						Geometry: "test",
// 					},
// 				},
// 			}, "test", nil
// 		},
// 	}

// 	server := NewOrderService(
// 		db,
// 		charger,
// 		cache,
// 		notifier,
// 		authServer.URL,
// 	)
// 	server.direction = direction

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := server.Update(ctx, tt.request())
// 			if err != nil && !tt.wantErr {
// 				t.Errorf("OrderService.Update() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.want != nil {
// 				got.ID = tt.want.ID
// 				got.CreatedAt = tt.want.CreatedAt
// 				got.UpdatedAt = tt.want.UpdatedAt
// 				if diff := cmp.Diff(tt.want, got); diff != "" {
// 					t.Errorf("OrderService.Update() mismatch (-want +got):\n%s", diff)
// 				}
// 			}
// 		})
// 	}
// }

// func setupOrder(t *testing.T, db *DB) *order.Order {
// 	t.Helper()
// 	order := &order.Order{
// 		Items: order.OrderItem{
// 			Points: []*order.Point{
// 				{
// 					Lat: 23.123,
// 					Lng: 23.123,
// 				},
// 				{
// 					Lat: 23.123,
// 					Lng: 23.123,
// 				},
// 			},
// 			Riders: 1,
// 		},
// 		Status: order.OrderStatusNew,
// 		Rider:  "test",
// 		Route: &order.DirectionResponse{
// 			Distance: 100,
// 			Duration: 100,
// 			Routes: []*order.Route{
// 				{
// 					Distance: 100,
// 					Duration: 100,
// 					Geometry: "test",
// 				},
// 			},
// 		},
// 		Currency:    "CUP",
// 		RouteString: base64.StdEncoding.EncodeToString([]byte("test")),
// 		Distance:    100,
// 		Duration:    100,
// 		CategoryPrice: []*order.CategoryPrice{
// 			{
// 				Category: order.VehicleCategoryX,
// 				Price:    100,
// 				Currency: "CUP",
// 			},
// 			{
// 				Category: order.VehicleCategoryXl,
// 				Price:    120,
// 				Currency: "CUP",
// 			},
// 		},
// 	}

// 	if err := storeOrder(context.Background(), db, order); err != nil {
// 		t.Fatal(err)
// 	}

// 	return order
// }
