package mongo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"cubawheeler.io/pkg/cubawheeler"
)

func TestPlanService_Create(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	var tests = []struct {
		name    string
		request func() *cubawheeler.PlanRequest
		want    *cubawheeler.Plan
		wantErr bool
	}{
		{
			name: "create a plan",
			request: func() *cubawheeler.PlanRequest {
				name := "test"
				trips := 10
				price := 100
				interval := cubawheeler.IntervalDay
				code := "test"
				return &cubawheeler.PlanRequest{
					Name:       &name,
					TotalTrips: &trips,
					Price:      &price,
					Interval:   &interval,
					Code:       &code,
				}
			},
			want: &cubawheeler.Plan{
				Name:     "test",
				Trips:    10,
				Price:    100,
				Interval: cubawheeler.IntervalDay,
				Code:     "test",
			},
		},
		{
			name: "create a plan with invalid interval",
			request: func() *cubawheeler.PlanRequest {
				trips := 10
				price := 100
				interval := cubawheeler.Interval("INVALID")
				code := "test"
				return &cubawheeler.PlanRequest{
					TotalTrips: &trips,
					Price:      &price,
					Interval:   &interval,
					Code:       &code,
				}
			},
			wantErr: true,
		},
	}

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("plans").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	service := NewPlanService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Create(ctx, tt.request())
			if err != nil && !tt.wantErr {
				t.Errorf("PlanService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want.ID = got.ID
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("response mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestPlanServiceUpdate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	var tests = []struct {
		name    string
		create  *cubawheeler.PlanRequest
		request func() *cubawheeler.PlanRequest
		want    *cubawheeler.Plan
		wantErr bool
	}{
		{
			name:   "update a plan",
			create: testPlan(t),
			request: func() *cubawheeler.PlanRequest {
				name := "updated-test"
				trips := 20
				price := 200
				interval := cubawheeler.IntervalWeek
				code := "updated-test"
				return &cubawheeler.PlanRequest{
					Name:       &name,
					TotalTrips: &trips,
					Price:      &price,
					Interval:   &interval,
					Code:       &code,
				}
			},
			want: &cubawheeler.Plan{
				ID:       "test",
				Name:     "updated-test",
				Trips:    20,
				Price:    200,
				Interval: cubawheeler.IntervalWeek,
				Code:     "updated-test",
			},
		},
		{
			name:   "update a plan with invalid interval",
			create: testPlan(t),
			request: func() *cubawheeler.PlanRequest {
				trips := 10
				price := 100
				interval := cubawheeler.Interval("INVALID")
				code := "test"
				return &cubawheeler.PlanRequest{
					TotalTrips: &trips,
					Price:      &price,
					Interval:   &interval,
					Code:       &code,
				}
			},
			wantErr: true,
		},
	}

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("plans").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	service := NewPlanService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := service.Create(ctx, testPlan(t))
			if err != nil && !tt.wantErr {
				t.Errorf("PlanService.Update() error = %v, wantErr %v", err, tt.wantErr)

			}
			req := tt.request()
			req.ID = plan.ID
			if tt.want != nil {
				tt.want.ID = plan.ID
				got, err := service.Update(ctx, req)
				if err != nil && !tt.wantErr {
					t.Errorf("PlanService.Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("response mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestPlanServiceFindByID(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	var tests = []struct {
		name    string
		create  *cubawheeler.PlanRequest
		request string
		want    *cubawheeler.Plan
		wantErr bool
	}{
		{
			name:    "find a plan by id",
			create:  testPlan(t),
			request: "test",
			want: &cubawheeler.Plan{
				ID:       "test",
				Name:     "test",
				Trips:    10,
				Price:    100,
				Interval: cubawheeler.IntervalDay,
				Code:     "test",
			},
		},
		{
			name:    "find a plan by invalid id",
			request: "INVALID",
			wantErr: true,
		},
	}

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("plans").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	service := NewPlanService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				p := testPlan(t)
				p.ID = tt.request
				_, err := service.Create(ctx, p)
				if err != nil && !tt.wantErr {
					t.Errorf("PlanService.FindByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			got, err := service.FindByID(ctx, tt.request)
			if err != nil && !tt.wantErr {
				t.Errorf("PlanService.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("response mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestPlanServiceFindAll(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	var tests = []struct {
		name    string
		create  []*cubawheeler.PlanRequest
		request *cubawheeler.PlanFilter
		want    []*cubawheeler.Plan
		token   string
		wantErr bool
	}{
		{
			name:    "find all plans",
			request: &cubawheeler.PlanFilter{},
			want: []*cubawheeler.Plan{
				{
					ID:       "test1",
					Name:     "test",
					Trips:    10,
					Price:    100,
					Interval: cubawheeler.IntervalDay,
					Code:     "test",
				},
				{
					ID:       "test2",
					Name:     "test",
					Trips:    10,
					Price:    100,
					Interval: cubawheeler.IntervalDay,
					Code:     "test",
				},
				{
					ID:       "test3",
					Name:     "test",
					Trips:    10,
					Price:    100,
					Interval: cubawheeler.IntervalDay,
					Code:     "test",
				},
			},
		},
		{
			name: "find all plans with limit",
			request: &cubawheeler.PlanFilter{
				Limit: 2,
			},
			token: "test3",
			want: []*cubawheeler.Plan{
				{
					ID:       "test1",
					Name:     "test",
					Trips:    10,
					Price:    100,
					Interval: cubawheeler.IntervalDay,
					Code:     "test",
				},
				{
					ID:       "test2",
					Name:     "test",
					Trips:    10,
					Price:    100,
					Interval: cubawheeler.IntervalDay,
					Code:     "test",
				},
			},
		},
	}

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("plans").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	service := NewPlanService(db)

	create := []*cubawheeler.PlanRequest{
		func() *cubawheeler.PlanRequest {
			t := testPlan(t)
			t.ID = "test1"
			return t
		}(),
		func() *cubawheeler.PlanRequest {
			t := testPlan(t)
			t.ID = "test2"
			return t
		}(),
		func() *cubawheeler.PlanRequest {
			t := testPlan(t)
			t.ID = "test3"
			return t
		}(),
	}
	for _, p := range create {
		_, err := service.Create(ctx, p)
		if err != nil {
			t.Fatalf("PlanService.FindAll() error = %v", err)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, token, err := service.FindAll(ctx, tt.request)
			if err != nil && !tt.wantErr {
				t.Fatalf("PlanService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if token != tt.token {
				t.Fatalf("PlanService.FindAll() token = %v, want %v", token, tt.token)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func testPlan(t *testing.T) *cubawheeler.PlanRequest {
	t.Helper()
	name := "test"
	trips := 10
	price := 100
	interval := cubawheeler.IntervalDay
	code := "test"
	return &cubawheeler.PlanRequest{
		Name:       &name,
		TotalTrips: &trips,
		Price:      &price,
		Interval:   &interval,
		Code:       &code,
	}
}
