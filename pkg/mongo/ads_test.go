package mongo

import (
	"context"
	"testing"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/google/go-cmp/cmp"
)

func TestAdsServiceCreate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})
	// Create a AdsService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.Collection(AdsCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewAdsService(db)

	var tests = []struct {
		name    string
		request func() *cubawheeler.AdsRequest
		want    *cubawheeler.Ads
		wantErr bool
	}{
		{
			name: "create a ads",
			request: func() *cubawheeler.AdsRequest {
				name := "test"
				description := "test"
				photo := "test"
				owner := "test"
				status := cubawheeler.AdsStatusActive
				priority := 1
				validFrom := 1
				validUntil := 1
				return &cubawheeler.AdsRequest{
					ID:          "test",
					Name:        name,
					Description: description,
					Photo:       photo,
					Owner:       owner,
					Status:      status,
					Priority:    priority,
					ValidFrom:   validFrom,
					ValidUntil:  validUntil,
				}
			},
			want: &cubawheeler.Ads{
				ID:          "test",
				Name:        "test",
				Description: "test",
				Photo:       "test",
				Owner:       "test",
				Status:      cubawheeler.AdsStatusActive,
				Priority:    1,
				ValidFrom:   1,
				ValidUntil:  1,
			},
		},
		{
			name: "create a ads with invalid status",
			request: func() *cubawheeler.AdsRequest {
				name := "test"
				description := "test"
				photo := "test"
				owner := "test"
				status := cubawheeler.AdsStatus("INVALID")
				priority := 1
				validFrom := 1
				validUntil := 1
				return &cubawheeler.AdsRequest{
					ID:          "test1",
					Name:        name,
					Description: description,
					Photo:       photo,
					Owner:       owner,
					Status:      status,
					Priority:    priority,
					ValidFrom:   validFrom,
					ValidUntil:  validUntil,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request()
			got, err := service.Create(ctx, request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AdsService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("AdsService.Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAdsServiceUpdate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})
	// Create a AdsService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.Collection(AdsCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewAdsService(db)

	var tests = []struct {
		name    string
		request func() *cubawheeler.AdsRequest
		want    *cubawheeler.Ads
		wantErr bool
	}{
		{
			name: "update a ads",
			request: func() *cubawheeler.AdsRequest {
				name := "test"
				description := "test"
				photo := "test"
				owner := "test"
				status := cubawheeler.AdsStatusActive
				priority := 1
				validFrom := 1
				validUntil := 1
				return &cubawheeler.AdsRequest{
					ID:          "test",
					Name:        name,
					Description: description,
					Photo:       photo,
					Owner:       owner,
					Status:      status,
					Priority:    priority,
					ValidFrom:   validFrom,
					ValidUntil:  validUntil,
				}
			},
			want: &cubawheeler.Ads{
				ID:          "test",
				Name:        "test",
				Description: "test",
				Photo:       "test",
				Owner:       "test",
				Status:      cubawheeler.AdsStatusActive,
				Priority:    1,
				ValidFrom:   1,
				ValidUntil:  1,
			},
		},
		{
			name: "update a ads with invalid status",
			request: func() *cubawheeler.AdsRequest {
				name := "test"
				description := "test"
				photo := "test"
				owner := "test"
				status := cubawheeler.AdsStatus("INVALID")
				priority := 1
				validFrom := 1
				validUntil := 1
				return &cubawheeler.AdsRequest{
					ID:          "test1",
					Name:        name,
					Description: description,
					Photo:       photo,
					Owner:       owner,
					Status:      status,
					Priority:    priority,
					ValidFrom:   validFrom,
					ValidUntil:  validUntil,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request()
			got, err := service.Update(ctx, request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AdsService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("AdsService.Update() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAdsServiceFindByID(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	var tests = []struct {
		name    string
		create  *cubawheeler.AdsRequest
		request string
		want    *cubawheeler.Ads
		wantErr bool
	}{
		{
			name:    "find a ads by id",
			create:  testAds(t),
			request: "test",
			want: &cubawheeler.Ads{
				ID:          "test",
				Name:        "test",
				Description: "test",
				Photo:       "test",
				Owner:       "test",
				Status:      cubawheeler.AdsStatusActive,
				Priority:    1,
				ValidFrom:   1,
				ValidUntil:  1,
			},
		},
		{
			name:    "find a ads by invalid id",
			request: "INVALID",
			wantErr: true,
		},
	}

	// Create a AdsService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.Collection(AdsCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewAdsService(db)
	ad := testAds(t)
	ad.ID = "test"
	service.Create(ctx, ad)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.FindById(ctx, tt.request)
			if err != nil && !tt.wantErr {
				t.Errorf("AdsService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAdsServiceFindAll(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	var tests = []struct {
		name    string
		create  *cubawheeler.AdsRequest
		request *cubawheeler.AdsRequest
		want    []*cubawheeler.Ads
		wantErr bool
	}{
		{
			name:    "find all ads",
			create:  testAds(t),
			request: &cubawheeler.AdsRequest{},
			want: []*cubawheeler.Ads{
				{
					ID:          "test",
					Name:        "test",
					Description: "test",
					Photo:       "test",
					Owner:       "test",
					Status:      cubawheeler.AdsStatusActive,
					Priority:    1,
					ValidFrom:   1,
					ValidUntil:  1,
				},
			},
		},
		{
			name:    "find all ads by invalid status",
			request: &cubawheeler.AdsRequest{Status: cubawheeler.AdsStatus("INVALID")},
			wantErr: true,
		},
	}

	// Create a AdsService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.Collection(AdsCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewAdsService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad, err := service.Create(ctx, testAds(t))
			if err != nil && !tt.wantErr {
				t.Errorf("AdsService.Update() error = %v, wantErr %v", err, tt.wantErr)

			}
			if tt.want != nil {
				tt.want[0].ID = ad.ID
			}
			got, _, err := service.FindAll(ctx, tt.request)
			if err != nil && !tt.wantErr {
				t.Errorf("AdsService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func testAds(t *testing.T) *cubawheeler.AdsRequest {
	name := "test"
	description := "test"
	photo := "test"
	owner := "test"
	status := cubawheeler.AdsStatusActive
	priority := 1
	validFrom := 1
	validUntil := 1
	return &cubawheeler.AdsRequest{
		ID:          cubawheeler.NewID().String(),
		Name:        name,
		Description: description,
		Photo:       photo,
		Owner:       owner,
		Status:      status,
		Priority:    priority,
		ValidFrom:   validFrom,
		ValidUntil:  validUntil,
	}
}
