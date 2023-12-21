package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"cubawheeler.io/pkg/cubawheeler"
)

func TestVehicleCategoryRateServiceCreate(t *testing.T) {
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{Role: cubawheeler.RoleAdmin})

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
		request func() *cubawheeler.VehicleCategoryRateRequest
		want    *cubawheeler.VehicleCategoryRate
		wantErr bool
	}{
		{
			name: "create a vehicle category rate",
			request: func() *cubawheeler.VehicleCategoryRateRequest {
				return &cubawheeler.VehicleCategoryRateRequest{
					ID:       "test",
					Category: "test",
					Factor:   1.0,
				}
			},
			want: &cubawheeler.VehicleCategoryRate{
				ID:       "test",
				Category: "test",
				Factor:   1.0,
			},
		},
		{
			name: "create a vehicle category rate with invalid factor",
			request: func() *cubawheeler.VehicleCategoryRateRequest {
				return &cubawheeler.VehicleCategoryRateRequest{
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
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{Role: cubawheeler.RoleAdmin})

	// Create a VehicleCategoryService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(VehicleCategoryRateCollection.String()).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	s.Create(ctx, &cubawheeler.VehicleCategoryRateRequest{
		ID:       "test",
		Category: "test",
		Factor:   1.0,
	})

	var tests = []struct {
		name    string
		request func() *cubawheeler.VehicleCategoryRateRequest
		want    *cubawheeler.VehicleCategoryRate
		wantErr bool
	}{
		{
			name: "update a vehicle category rate",
			request: func() *cubawheeler.VehicleCategoryRateRequest {
				return &cubawheeler.VehicleCategoryRateRequest{
					ID:       "test",
					Category: "test",
					Factor:   1.5,
				}
			},
			want: &cubawheeler.VehicleCategoryRate{
				ID:       "test",
				Category: "test",
				Factor:   1.5,
			},
		},
		{
			name: "update a vehicle category rate with invalid factor",
			request: func() *cubawheeler.VehicleCategoryRateRequest {
				return &cubawheeler.VehicleCategoryRateRequest{
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
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{Role: cubawheeler.RoleAdmin})

	// Create a VehicleCategoryService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(VehicleCategoryRateCollection.String()).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	s.Create(ctx, &cubawheeler.VehicleCategoryRateRequest{
		ID:       "test",
		Category: "test",
		Factor:   1.0,
	})

	var tests = []struct {
		name    string
		request string
		want    *cubawheeler.VehicleCategoryRate
		wantErr bool
	}{
		{
			name:    "find a vehicle category rate by id",
			request: "test",
			want: &cubawheeler.VehicleCategoryRate{
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
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{Role: cubawheeler.RoleAdmin})

	// Create a VehicleCategoryService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(VehicleCategoryRateCollection.String()).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	s.Create(ctx, &cubawheeler.VehicleCategoryRateRequest{
		ID:       "test",
		Category: "test",
		Factor:   1.0,
	})

	var tests = []struct {
		name    string
		request cubawheeler.VehicleCategory
		want    *cubawheeler.VehicleCategoryRate
		wantErr bool
	}{
		{
			name:    "find a vehicle category rate by category",
			request: "test",
			want: &cubawheeler.VehicleCategoryRate{
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
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{Role: cubawheeler.RoleAdmin})

	// Create a VehicleCategoryService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(VehicleCategoryRateCollection.String()).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewVehicleCategoryRateService(db)
	for i := 0; i < 10; i++ {
		s.Create(ctx, &cubawheeler.VehicleCategoryRateRequest{
			ID:       fmt.Sprintf("test-%d", i),
			Category: cubawheeler.VehicleCategory(fmt.Sprintf("test-%d", i)),
			Factor:   1.0,
		})
	}

	var tests = []struct {
		name    string
		request cubawheeler.VehicleCategoryRateFilter
		want    []*cubawheeler.VehicleCategoryRate
		token   string
		wantErr bool
	}{
		{
			name:    "find all vehicle category rates",
			request: cubawheeler.VehicleCategoryRateFilter{},
			want: func() []*cubawheeler.VehicleCategoryRate {
				var rates []*cubawheeler.VehicleCategoryRate
				for i := 0; i < 10; i++ {
					rates = append(rates, &cubawheeler.VehicleCategoryRate{
						ID:       fmt.Sprintf("test-%d", i),
						Category: cubawheeler.VehicleCategory(fmt.Sprintf("test-%d", i)),
						Factor:   1.0,
					})
				}
				return rates
			}(),
		},
		{
			name:    "find all vehicle category rates by category",
			request: cubawheeler.VehicleCategoryRateFilter{Category: []cubawheeler.VehicleCategory{"test-1"}},
			want: []*cubawheeler.VehicleCategoryRate{
				{
					ID:       "test-1",
					Category: "test-1",
					Factor:   1.0,
				},
			},
		},
		{
			name:    "find all vehicle category rates with limit",
			request: cubawheeler.VehicleCategoryRateFilter{Limit: 1},
			token:   "test-1",
			want: []*cubawheeler.VehicleCategoryRate{
				{
					ID:       "test-0",
					Category: "test-0",
					Factor:   1.0,
				},
			},
		},
		{
			name:    "find all vehicle category rates by invalid category",
			request: cubawheeler.VehicleCategoryRateFilter{Category: []cubawheeler.VehicleCategory{"INVALID"}},
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
