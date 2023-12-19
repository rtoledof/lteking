package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"cubawheeler.io/pkg/cubawheeler"
)

func TestRateServiceCreate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(
		context.Background(),
		&cubawheeler.User{
			ID:   "test",
			Role: cubawheeler.RoleAdmin,
		},
	)

	db := NewTestDB()
	s := NewRateService(db)
	defer func() {
		db.client.Database(database).Collection("rates").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()

	currenctTime := time.Now().UTC()
	var tests = []struct {
		name    string
		request func() *cubawheeler.RateRequest
		want    *cubawheeler.Rate
		wantErr bool
	}{
		{
			name: "create a rate",
			request: func() *cubawheeler.RateRequest {
				code := "test"
				basePrice := 100
				pricePerMin := 10
				pricePerKm := 10
				pricePerPassenger := 10
				pricePerBaggage := 10
				startDate := int(currenctTime.Unix())
				endDate := int(currenctTime.AddDate(0, 1, 0).Unix())
				minKm := 10
				maxKm := 100
				return &cubawheeler.RateRequest{
					ID:                "test",
					Code:              code,
					BasePrice:         basePrice,
					PricePerMin:       &pricePerMin,
					PricePerKm:        &pricePerKm,
					PricePerPassenger: &pricePerPassenger,
					PricePerBaggage:   &pricePerBaggage,
					StartDate:         &startDate,
					EndDate:           &endDate,
					MinKm:             &minKm,
					MaxKm:             &maxKm,
				}
			},
			want: &cubawheeler.Rate{
				Code:              "test",
				BasePrice:         100,
				PricePerMin:       10,
				PricePerKm:        10,
				PricePerBaggage:   10,
				PricePerPassenger: 10,
				StartDate:         int(time.Now().UTC().Unix()),
				EndDate:           int(currenctTime.AddDate(0, 1, 0).Unix()),
				MinKm:             10,
				MaxKm:             100,
			},
		},
		{
			name: "create a rate with empty code",
			request: func() *cubawheeler.RateRequest {
				code := ""
				basePrice := 100
				pricePerMin := 10
				pricePerKm := 10
				pricePerPassenger := 10
				pricePerBaggage := 10
				startDate := int(currenctTime.Unix())
				endDate := int(currenctTime.AddDate(0, 1, 0).Unix())
				minKm := 10
				maxKm := 100
				return &cubawheeler.RateRequest{
					ID:                "test",
					Code:              code,
					BasePrice:         basePrice,
					PricePerMin:       &pricePerMin,
					PricePerKm:        &pricePerKm,
					PricePerPassenger: &pricePerPassenger,
					PricePerBaggage:   &pricePerBaggage,
					StartDate:         &startDate,
					EndDate:           &endDate,
					MinKm:             &minKm,
					MaxKm:             &maxKm,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create(ctx, *tt.request())
			if (err != nil) != tt.wantErr {
				t.Errorf("RateService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRateServiceFindAll(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(
		context.Background(),
		&cubawheeler.User{
			ID:   "test",
			Role: cubawheeler.RoleAdmin,
		},
	)
	database = "test"
	db := NewTestDB()
	s := NewRateService(db)
	defer func() {
		db.client.Database(database).Collection("rates").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()

	rates := []*cubawheeler.Rate{
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test1"
			rate.Code = "test_code"
			return rate
		}(),
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test2"
			return rate
		}(),
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test3"
			return rate
		}(),
	}

	for _, rate := range rates {
		if err := insertRate(ctx, db, rate); err != nil {
			t.Fatalf("unable to insert rate: %v", err)
		}
	}

	currenctTime := time.Now().UTC()
	var tests = []struct {
		name    string
		request cubawheeler.RateFilter
		want    []*cubawheeler.Rate
		cursor  string
		wantErr bool
	}{
		{
			name: "find all rates",
			request: cubawheeler.RateFilter{
				Code: []string{"test_code"},
			},
			want: []*cubawheeler.Rate{
				{
					ID:                "test1",
					Code:              "test_code",
					BasePrice:         100,
					PricePerMin:       10,
					PricePerKm:        10,
					PricePerBaggage:   10,
					PricePerPassenger: 10,
					StartDate:         int(currenctTime.Unix()),
					EndDate:           int(currenctTime.AddDate(0, 1, 0).Unix()),
					MinKm:             10,
					MaxKm:             100,
				},
			},
		},
		{
			name:    "find all rates",
			request: cubawheeler.RateFilter{},
			want: []*cubawheeler.Rate{
				func() *cubawheeler.Rate {
					rate := testRate(t)
					rate.ID = "test1"
					rate.Code = "test_code"
					return rate
				}(),
				func() *cubawheeler.Rate {
					rate := testRate(t)
					rate.ID = "test2"
					return rate
				}(),
				func() *cubawheeler.Rate {
					rate := testRate(t)
					rate.ID = "test3"
					return rate
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cursor, err := s.FindAll(ctx, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RateService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
			if cursor != tt.cursor {
				t.Errorf("RateService.FindAll() cursor = %v, want %v", cursor, tt.cursor)
			}
		})
	}
}

func TestRateServiceFindByID(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(
		context.Background(),
		&cubawheeler.User{
			ID:   "test",
			Role: cubawheeler.RoleAdmin,
		},
	)

	database = "test"
	db := NewTestDB()
	s := NewRateService(db)
	defer func() {
		db.client.Database(database).Collection("rates").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()

	rates := []*cubawheeler.Rate{
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test1"
			rate.Code = "test_code"
			return rate
		}(),
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test2"
			return rate
		}(),
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test3"
			return rate
		}(),
	}

	for _, rate := range rates {
		if err := insertRate(ctx, db, rate); err != nil {
			t.Fatalf("unable to insert rate: %v", err)
		}
	}

	var tests = []struct {
		name    string
		request string
		want    *cubawheeler.Rate
		wantErr bool
	}{
		{
			name:    "find a rate",
			request: "test1",
			want: func() *cubawheeler.Rate {
				rate := testRate(t)
				rate.ID = "test1"
				rate.Code = "test_code"
				return rate
			}(),
		},
		{
			name:    "find a rate with invalid id",
			request: "test4",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.FindByID(ctx, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RateService.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRateServiceUpdate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(
		context.Background(),
		&cubawheeler.User{
			ID:   "test",
			Role: cubawheeler.RoleAdmin,
		},
	)

	db := NewTestDB()
	s := NewRateService(db)
	defer func() {
		db.client.Database(database).Collection("rates").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()

	rates := []*cubawheeler.Rate{
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test1"
			rate.Code = "test update"
			return rate
		}(),
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test2"
			rate.Code = "test update"
			return rate
		}(),
		func() *cubawheeler.Rate {
			rate := testRate(t)
			rate.ID = "test3"
			rate.Code = "test update"
			return rate
		}(),
	}

	for _, rate := range rates {
		if err := insertRate(ctx, db, rate); err != nil {
			t.Fatalf("unable to insert rate: %v", err)
		}
	}

	currenctTime := time.Now().UTC()
	var tests = []struct {
		name    string
		request func() *cubawheeler.RateRequest
		want    *cubawheeler.Rate
		wantErr bool
	}{
		{
			name: "update a rate",
			request: func() *cubawheeler.RateRequest {
				code := "test update"
				basePrice := 100
				pricePerMin := 10
				pricePerKm := 10
				pricePerPassenger := 10
				pricePerBaggage := 10
				startDate := int(currenctTime.Unix())
				endDate := int(currenctTime.AddDate(0, 1, 0).Unix())
				minKm := 10
				maxKm := 100
				return &cubawheeler.RateRequest{
					ID:                "test1",
					Code:              code,
					BasePrice:         basePrice,
					PricePerMin:       &pricePerMin,
					PricePerKm:        &pricePerKm,
					PricePerPassenger: &pricePerPassenger,
					PricePerBaggage:   &pricePerBaggage,
					StartDate:         &startDate,
					EndDate:           &endDate,
					MinKm:             &minKm,
					MaxKm:             &maxKm,
				}
			},
			want: &cubawheeler.Rate{
				ID:                "test1",
				Code:              "test update",
				BasePrice:         100,
				PricePerMin:       10,
				PricePerKm:        10,
				PricePerBaggage:   10,
				PricePerPassenger: 10,
				StartDate:         int(currenctTime.Unix()),
				EndDate:           int(currenctTime.AddDate(0, 1, 0).Unix()),
				MinKm:             10,
				MaxKm:             100,
			},
		},
		{
			name: "update a rate with empty code",
			request: func() *cubawheeler.RateRequest {
				code := ""
				basePrice := 100
				pricePerMin := 10
				pricePerKm := 10
				pricePerPassenger := 10
				pricePerBaggage := 10
				startDate := int(currenctTime.Unix())
				endDate := int(currenctTime.AddDate(0, 1, 0).Unix())
				minKm := 10
				maxKm := 100
				return &cubawheeler.RateRequest{
					ID:                "test",
					Code:              code,
					BasePrice:         basePrice,
					PricePerMin:       &pricePerMin,
					PricePerKm:        &pricePerKm,
					PricePerPassenger: &pricePerPassenger,
					PricePerBaggage:   &pricePerBaggage,
					StartDate:         &startDate,
					EndDate:           &endDate,
					MinKm:             &minKm,
					MaxKm:             &maxKm,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Update(ctx, tt.request())
			if (err != nil) != tt.wantErr {
				t.Errorf("RateService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func testRate(t *testing.T) *cubawheeler.Rate {
	return &cubawheeler.Rate{
		Code:              "test",
		BasePrice:         100,
		PricePerMin:       10,
		PricePerKm:        10,
		PricePerBaggage:   10,
		PricePerPassenger: 10,
		StartDate:         int(time.Now().UTC().Unix()),
		EndDate:           int(time.Now().UTC().AddDate(0, 1, 0).Unix()),
		MinKm:             10,
		MaxKm:             100,
	}
}
