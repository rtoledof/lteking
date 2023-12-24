package mongo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

func TestOrderStatisticsAddOrder(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("plans").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	service := NewOrderStatistics(db)

	var tests = []struct {
		name    string
		request *cubawheeler.OrderStatistics
		want    *cubawheeler.OrderStatistics
		wantErr bool
	}{
		{
			name: "create a plan",
			request: &cubawheeler.OrderStatistics{
				ID:   "test",
				User: "test",
				Orders: map[string]*cubawheeler.Statistics{
					"2006-01-02": {
						Total: 1,
						Amount: currency.Amount{
							Amount:   100,
							Currency: currency.MustParse("CUP"),
						},
						Orders: []string{"1"},
					},
					"2006-01": {
						Total: 1,
						Amount: currency.Amount{
							Amount:   100,
							Currency: currency.MustParse("CUP"),
						},
						Orders: []string{"1"},
					},
					"2006": {
						Total: 1,
						Amount: currency.Amount{
							Amount:   100,
							Currency: currency.MustParse("CUP"),
						},
						Orders: []string{"1"},
					},
				},
			},
			want: &cubawheeler.OrderStatistics{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.AddOrder(ctx, *tt.request); err != nil && !tt.wantErr {
				t.Fatalf("OrderStatistics.AddOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderStatisticsFindAllStatistics(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(StatisticsCollection.String()).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewOrderStatistics(db)

	var statisticst = []*cubawheeler.OrderStatistics{
		{
			ID:   cubawheeler.NewID().String(),
			User: "test",
			Orders: map[string]*cubawheeler.Statistics{
				"2006-01-02": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
				"2006-01-01": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
				"2006-01": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
				"2006": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
			},
		},
	}

	for _, tt := range statisticst {
		if err := service.AddOrder(ctx, *tt); err != nil {
			t.Fatalf("OrderStatistics.AddOrder() error = %v", err)
		}
	}

	var tests = []struct {
		name    string
		filter  cubawheeler.OrderStatisticsFilter
		want    []*cubawheeler.OrderStatistics
		token   string
		wantErr bool
	}{
		{
			name: "create a plan",
			want: statisticst,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statistics, token, err := service.FindAllStatistics(ctx, tt.filter)
			if err != nil && !tt.wantErr {
				t.Fatalf("OrderStatistics.AddOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if token != tt.token {
				t.Errorf("OrderStatistics.AddOrder() token = %v, want %v", token, tt.token)
			}
			if diff := cmp.Diff(tt.want, statistics); diff != "" {
				t.Errorf("OrderStatistics.AddOrder() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrderStatisticsFindStatistictsByUser(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection(StatisticsCollection.String()).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewOrderStatistics(db)

	var statisticst = []*cubawheeler.OrderStatistics{
		{
			ID:   cubawheeler.NewID().String(),
			User: "test",
			Orders: map[string]*cubawheeler.Statistics{
				"2006-01-02": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
				"2006-01-01": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
				"2006-01": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
				"2006": {
					Total: 1,
					Amount: currency.Amount{
						Amount:   100,
						Currency: currency.MustParse("CUP"),
					},
					Orders: []string{"1"},
				},
			},
		},
	}

	for _, tt := range statisticst {
		if err := service.AddOrder(ctx, *tt); err != nil {
			t.Fatalf("OrderStatistics.AddOrder() error = %v", err)
		}
	}

	var tests = []struct {
		name    string
		user    string
		want    *cubawheeler.OrderStatistics
		wantErr bool
	}{
		{
			name: "create a plan",
			user: "test",
			want: statisticst[0],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statistics, err := service.FindStatistictsByUser(ctx, tt.user)
			if err != nil && !tt.wantErr {
				t.Fatalf("OrderStatistics.AddOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, statistics); diff != "" {
				t.Errorf("OrderStatistics.AddOrder() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
