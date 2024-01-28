package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/jwt"

	"order.io/pkg/order"
)

func TestVehicleCategoryRateServiceCreate(t *testing.T) {
	ctx := prepareContext(t, order.RoleAdmin)

	// Create a VehicleCategoryService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(VehicleCategoryRateCollection.String()).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)

	var tests = []struct {
		name    string
		request func() *order.VehicleCategoryRateRequest
		want    *order.VehicleCategoryRate
		wantErr bool
	}{
		{
			name: "create a vehicle category rate",
			request: func() *order.VehicleCategoryRateRequest {
				return &order.VehicleCategoryRateRequest{
					ID:       "test",
					Category: "test",
					Factor:   1.0,
				}
			},
			want: &order.VehicleCategoryRate{
				ID:       "test",
				Category: "test",
				Factor:   1.0,
			},
		},
		{
			name: "create a vehicle category rate with invalid factor",
			request: func() *order.VehicleCategoryRateRequest {
				return &order.VehicleCategoryRateRequest{
					ID:       "test",
					Category: "test",
					Factor:   -1.0,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(ctx, tt.request())
			if (err != nil) != tt.wantErr {
				t.Fatalf("VehicleCategoryRateService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(got.ID, tt.want.ID); diff != "" {
				t.Fatalf("VehicleCategoryRateService.Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestVehicleCategoryRateServiceUpdate(t *testing.T) {
	ctx := prepareContext(t, order.RoleAdmin)

	// Create a VehicleCategoryService instance with the mock collection
	db := NewTestDB()
	defer func() {
		db.Collection(VehicleCategoryRateCollection).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	s.Create(ctx, &order.VehicleCategoryRateRequest{
		ID:       "test",
		Category: "test",
		Factor:   1.0,
	})

	var tests = []struct {
		name    string
		request func() *order.VehicleCategoryRateRequest
		want    *order.VehicleCategoryRate
		wantErr bool
	}{
		{
			name: "update a vehicle category rate",
			request: func() *order.VehicleCategoryRateRequest {
				return &order.VehicleCategoryRateRequest{
					ID:       "test",
					Category: "test",
					Factor:   1.5,
				}
			},
			want: &order.VehicleCategoryRate{
				ID:       "test",
				Category: "test",
				Factor:   1.5,
			},
		},
		{
			name: "update a vehicle category rate with invalid factor",
			request: func() *order.VehicleCategoryRateRequest {
				return &order.VehicleCategoryRateRequest{
					ID:       "test",
					Category: "test",
					Factor:   -1.0,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Update(ctx, tt.request())
			if (err != nil) != tt.wantErr {
				t.Fatalf("VehicleCategoryRateService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if diff := cmp.Diff(got.ID, tt.want.ID); diff != "" {
					t.Fatalf("VehicleCategoryRateService.Update() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestVehicleCategoryRateServiceFindByID(t *testing.T) {
	ctx := prepareContext(t, order.RoleAdmin)

	// Create a VehicleCategoryService instance with the mock collection
	db := NewTestDB()
	defer func() {
		db.Collection(VehicleCategoryRateCollection).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	s.Create(ctx, &order.VehicleCategoryRateRequest{
		ID:       "test",
		Category: "test",
		Factor:   1.0,
	})

	var tests = []struct {
		name    string
		request string
		want    *order.VehicleCategoryRate
		wantErr bool
	}{
		{
			name:    "find a vehicle category rate by id",
			request: "test",
			want: &order.VehicleCategoryRate{
				ID:       "test",
				Category: "test",
				Factor:   1.0,
			},
		},
		{
			name:    "find a vehicle category rate by invalid id",
			request: "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.FindByID(ctx, tt.request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("VehicleCategoryRateService.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if diff := cmp.Diff(got.ID, tt.want.ID); diff != "" {
					t.Fatalf("VehicleCategoryRateService.FindByID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestVehicleCategoryRateServiceFindByCategory(t *testing.T) {
	ctx := prepareContext(t, order.RoleAdmin)

	// Create a VehicleCategoryService instance with the mock collection
	db := NewTestDB()
	defer func() {
		db.Collection(VehicleCategoryRateCollection).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	s.Create(ctx, &order.VehicleCategoryRateRequest{
		ID:       "test",
		Category: "test",
		Factor:   1.0,
	})

	var tests = []struct {
		name    string
		request order.VehicleCategory
		want    *order.VehicleCategoryRate
		wantErr bool
	}{
		{
			name:    "find a vehicle category rate by category",
			request: "test",
			want: &order.VehicleCategoryRate{
				ID:       "test",
				Category: "test",
				Factor:   1.0,
			},
		},
		{
			name:    "find a vehicle category rate by invalid category",
			request: "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.FindByCategory(ctx, tt.request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("VehicleCategoryRateService.FindByCategory() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if diff := cmp.Diff(got.ID, tt.want.ID); diff != "" {
					t.Fatalf("VehicleCategoryRateService.FindByCategory() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestVehicleCategoryRateServiceFindAll(t *testing.T) {
	ctx := prepareContext(t, order.RoleAdmin)

	// Create a VehicleCategoryService instance with the mock collection
	db := NewTestDB()
	defer func() {
		db.Collection(VehicleCategoryRateCollection).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	for i := 0; i < 10; i++ {
		s.Create(ctx, &order.VehicleCategoryRateRequest{
			ID:       fmt.Sprintf("test-%d", i),
			Category: order.VehicleCategory(fmt.Sprintf("test-%d", i)),
			Factor:   1.0,
		})
	}

	var tests = []struct {
		name    string
		request order.VehicleCategoryRateFilter
		want    []*order.VehicleCategoryRate
		token   string
		wantErr bool
	}{
		{
			name:    "find all vehicle category rates",
			request: order.VehicleCategoryRateFilter{},
			want: func() []*order.VehicleCategoryRate {
				var rates []*order.VehicleCategoryRate
				for i := 0; i < 10; i++ {
					rates = append(rates, &order.VehicleCategoryRate{
						ID:       fmt.Sprintf("test-%d", i),
						Category: order.VehicleCategory(fmt.Sprintf("test-%d", i)),
						Factor:   1.0,
					})
				}
				return rates
			}(),
		},
		{
			name:    "find all vehicle category rates by category",
			request: order.VehicleCategoryRateFilter{Category: []order.VehicleCategory{"test-1"}},
			want: []*order.VehicleCategoryRate{
				{
					ID:       "test-1",
					Category: "test-1",
					Factor:   1.0,
				},
			},
		},
		{
			name:    "find all vehicle category rates with limit",
			request: order.VehicleCategoryRateFilter{Limit: 1},
			token:   "test-1",
			want: []*order.VehicleCategoryRate{
				{
					ID:       "test-0",
					Category: "test-0",
					Factor:   1.0,
				},
			},
		},
		{
			name:    "find all vehicle category rates by invalid category",
			request: order.VehicleCategoryRateFilter{Category: []order.VehicleCategory{"INVALID"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, token, err := s.FindAll(ctx, tt.request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("VehicleCategoryRateService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}

			if token != tt.token {
				t.Fatalf("VehicleCategoryRateService.FindAll() token = %v, want %v", token, tt.token)
			}

			if tt.want != nil {
				if diff := cmp.Diff(len(got), len(tt.want)); diff != "" {
					t.Fatalf("VehicleCategoryRateService.FindAll() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func setupVethiclesCategoriesRate(t *testing.T, ctx context.Context, db *DB) []*order.VehicleCategoryRate {
	t.Helper()
	var rates = []*order.VehicleCategoryRate{
		{
			ID:       string(order.NewID()),
			Category: order.VehicleCategoryX,
			Factor:   1.0,
		},
		{
			ID:       order.NewID().String(),
			Category: order.VehicleCategoryXl,
			Factor:   1.2,
		},
	}

	for _, rate := range rates {
		if err := insertVehicleCategoryRate(ctx, db, rate); err != nil {
			t.Fatalf("unable to insert rate: %v", err)
		}
	}
	return rates
}

func prepateContext(t *testing.T, roles ...order.Role) context.Context {
	t.Helper()
	ctx := context.Background()

	token := jwt.New()
	token.Set("id", order.NewID().String())
	user := order.User{
		ID:    order.NewID().String(),
		Name:  "test",
		Email: "test",
		Role:  "rider",
	}
	if roles != nil {
		user.Role = roles[0]
	}
	userData, _ := json.Marshal(user)
	token.Set("user", userData)
	return jwtauth.NewContext(ctx, token, nil)
}
