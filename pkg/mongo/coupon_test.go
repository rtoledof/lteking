package mongo

import (
	"context"
	"testing"
	"time"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
	"github.com/google/go-cmp/cmp"
)

func TestCouponServiceCreate(t *testing.T) {
	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	ctx := context.Background()
	defer func() {
		db.Collection(CouponCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewCouponService(db)

	currentTime := time.Now().UTC()
	var test = []struct {
		name    string
		request func() *cubawheeler.CouponRequest
		want    *cubawheeler.Coupon
		wantErr bool
	}{
		{
			name: "create a coupon",
			request: func() *cubawheeler.CouponRequest {
				code := "test"
				percent := float64(10)
				status := cubawheeler.CouponStatusActive
				validFrom := currentTime.Unix()
				validUntil := currentTime.Add(time.Hour * 24).UTC().Unix()
				return &cubawheeler.CouponRequest{
					ID:         "test",
					Code:       code,
					Percent:    &percent,
					Amount:     100,
					Currency:   "CUP",
					Status:     status,
					ValidFrom:  &validFrom,
					ValidUntil: &validUntil,
				}
			},
			want: &cubawheeler.Coupon{
				ID:      "test",
				Code:    "test",
				Percent: func() *float64 { f := float64(10); return &f }(),
				Amount: currency.Amount{
					Amount:   100,
					Currency: currency.MustParse("CUP"),
				},
				Status:     cubawheeler.CouponStatusActive,
				ValidFrom:  currentTime.Unix(),
				ValidUntil: currentTime.Add(time.Hour * 24).Unix(),
			},
		},
		{
			name: "create a coupon with invalid status",
			request: func() *cubawheeler.CouponRequest {
				code := "test"
				percent := 10.0
				status := cubawheeler.CouponStatus("INVALID")
				validFrom := currentTime.Unix()
				validUntil := currentTime.Add(time.Hour * 24).Unix()
				return &cubawheeler.CouponRequest{
					Code:       code,
					Percent:    &percent,
					Amount:     100,
					Currency:   "CUP",
					Status:     status,
					ValidFrom:  &validFrom,
					ValidUntil: &validUntil,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request()
			got, err := service.Create(ctx, request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CouponService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("CouponService.Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCouponServiceFindByCode(t *testing.T) {
	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test"})
	defer func() {
		db.Collection(CouponCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewCouponService(db)

	currentTime := time.Now().UTC()
	coupon := &cubawheeler.Coupon{
		ID:      "test",
		Code:    "test",
		Percent: func() *float64 { f := float64(10); return &f }(),
		Amount: currency.Amount{
			Amount:   100,
			Currency: currency.MustParse("CUP"),
		},
		Status:     cubawheeler.CouponStatusActive,
		ValidFrom:  currentTime.Unix(),
		ValidUntil: currentTime.Add(time.Hour * 24).Unix(),
	}
	if _, err := db.Collection(CouponCollection).InsertOne(ctx, coupon); err != nil {
		t.Fatal(err)
	}

	var test = []struct {
		name    string
		code    string
		want    *cubawheeler.Coupon
		wantErr bool
	}{
		{
			name: "find a coupon by code",
			code: "test",
			want: coupon,
		},
		{
			name:    "find a coupon by code that does not exist",
			code:    "test2",
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.FindByCode(ctx, tt.code)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CouponService.FindByCode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("CouponService.FindByCode() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCouponServiceRedeem(t *testing.T) {
	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	user := &cubawheeler.User{ID: cubawheeler.NewID().String()}
	ctx := cubawheeler.NewContextWithUser(context.Background(), user)
	us := NewUserService(db, nil, nil, nil)
	us.CreateUser(ctx, user)
	defer func() {
		db.Collection(CouponCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewCouponService(db)

	currentTime := time.Now().UTC()
	coupon := &cubawheeler.Coupon{
		ID:      "test",
		Code:    "test",
		Percent: func() *float64 { f := float64(10); return &f }(),
		Amount: currency.Amount{
			Amount:   100,
			Currency: currency.MustParse("CUP"),
		},
		Status:     cubawheeler.CouponStatusActive,
		ValidFrom:  currentTime.Unix(),
		ValidUntil: currentTime.Add(time.Hour * 24).Unix(),
	}
	if _, err := db.Collection(CouponCollection).InsertOne(ctx, coupon); err != nil {
		t.Fatal(err)
	}

	var test = []struct {
		name    string
		code    string
		want    *cubawheeler.Coupon
		wantErr bool
	}{
		{
			name: "redeem a coupon",
			code: "test",
			want: &cubawheeler.Coupon{
				ID:      "test",
				Code:    "test",
				Percent: func() *float64 { f := float64(10); return &f }(),
				Amount: currency.Amount{
					Amount:   100,
					Currency: currency.MustParse("CUP"),
				},
				Status:     cubawheeler.CouponStatusRedeemed,
				ValidFrom:  currentTime.Unix(),
				ValidUntil: currentTime.Add(time.Hour * 24).Unix(),
			},
		},
		{
			name:    "redeem a coupon that does not exist",
			code:    "test2",
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Redeem(ctx, tt.code)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CouponService.Redeem() error = %v, wantErr %v", err, tt.want)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("CouponService.Redeem() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCouponServiceFindAll(t *testing.T) {
	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test"})
	defer func() {
		db.Collection(CouponCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewCouponService(db)

	currentTime := time.Now().UTC()
	coupon := &cubawheeler.Coupon{
		ID:      cubawheeler.NewID().String(),
		Code:    "test",
		Percent: func() *float64 { f := float64(10); return &f }(),
		Amount: currency.Amount{
			Amount:   100,
			Currency: currency.MustParse("CUP"),
		},
		Status:     cubawheeler.CouponStatusActive,
		ValidFrom:  currentTime.Unix(),
		ValidUntil: currentTime.Add(time.Hour * 24).Unix(),
	}
	if _, err := db.Collection(CouponCollection).InsertOne(ctx, coupon); err != nil {
		t.Fatal(err)
	}

	var test = []struct {
		name    string
		request func() *cubawheeler.CouponRequest
		want    []*cubawheeler.Coupon
		wantErr bool
	}{
		{
			name: "find all coupons",
			request: func() *cubawheeler.CouponRequest {
				return &cubawheeler.CouponRequest{}
			},
			want: []*cubawheeler.Coupon{coupon},
		},
		{
			name: "find all coupons with code",
			request: func() *cubawheeler.CouponRequest {
				return &cubawheeler.CouponRequest{
					Code: "test",
				}
			},
			want: []*cubawheeler.Coupon{coupon},
		},
		{
			name: "find all coupons with status",
			request: func() *cubawheeler.CouponRequest {
				return &cubawheeler.CouponRequest{
					Status: cubawheeler.CouponStatusActive,
				}
			},
			want: []*cubawheeler.Coupon{coupon},
		},
		{
			name: "find all coupons with valid from",
			request: func() *cubawheeler.CouponRequest {
				from := currentTime.Unix()
				return &cubawheeler.CouponRequest{
					ValidFrom: &from,
				}
			},
			want: []*cubawheeler.Coupon{coupon},
		},
		{
			name: "find all coupons with valid until",
			request: func() *cubawheeler.CouponRequest {
				until := currentTime.Add(time.Hour * 24).Unix()
				return &cubawheeler.CouponRequest{
					ValidUntil: &until,
				}
			},
			want: []*cubawheeler.Coupon{coupon},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request()
			got, _, err := service.FindAll(ctx, request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CouponService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("CouponService.FindAll() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCouponServiceFindByID(t *testing.T) {
	// Create a CouponService instance with the mock collection
	database = "test"
	db := NewTestDB()
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test"})
	defer func() {
		db.Collection(CouponCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewCouponService(db)

	currentTime := time.Now().UTC()
	coupon := &cubawheeler.Coupon{
		ID:      cubawheeler.NewID().String(),
		Code:    "test",
		Percent: func() *float64 { f := float64(10); return &f }(),
		Amount: currency.Amount{
			Amount:   100,
			Currency: currency.MustParse("CUP"),
		},
		Status:     cubawheeler.CouponStatusActive,
		ValidFrom:  currentTime.Unix(),
		ValidUntil: currentTime.Add(time.Hour * 24).Unix(),
	}
	service.Create(ctx, &cubawheeler.CouponRequest{
		ID:         coupon.ID,
		Code:       coupon.Code,
		Percent:    coupon.Percent,
		Amount:     coupon.Amount.Amount,
		Currency:   coupon.Amount.Currency.String(),
		Status:     coupon.Status,
		ValidFrom:  &coupon.ValidFrom,
		ValidUntil: &coupon.ValidUntil,
	})

	var test = []struct {
		name    string
		id      string
		want    *cubawheeler.Coupon
		wantErr bool
	}{
		{
			name: "find all coupons",
			id:   coupon.ID,
			want: coupon,
		},
		{
			name:    "finding a coupon that does not exist",
			id:      "test",
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.FindByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CouponService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil {
				tt.want.CreatedAt = got.CreatedAt
				tt.want.UpdatedAt = got.UpdatedAt
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Fatalf("CouponService.FindAll() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
