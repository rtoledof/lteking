package mongo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
)

func TestChargeServiceCreate(t *testing.T) {
	ctx := cubawheeler.NewContextWithUser(context.Background(), &cubawheeler.User{ID: "test", Role: cubawheeler.RoleAdmin})

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("plans").Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	service := NewChargeService(db)

	var tests = []struct {
		name    string
		request func() *cubawheeler.ChargeRequest
		want    *cubawheeler.Charge
		wantErr bool
	}{
		{
			name: "create a charge",
			request: func() *cubawheeler.ChargeRequest {
				method := cubawheeler.ChargeMethodCard
				reference := "test"
				return &cubawheeler.ChargeRequest{
					Amount:      100,
					Currency:    "usd",
					Order:       "test",
					Method:      method,
					Reference:   &reference,
					Description: "test",
				}
			},
			want: &cubawheeler.Charge{
				Amount: currency.Amount{
					Currency: currency.MustParse("usd"),
					Amount:   100,
				},
				Description:       "test",
				Order:             "test",
				Disputed:          false,
				ReceiptEmail:      "",
				Status:            cubawheeler.ChargeStatusPending,
				Method:            "card",
				ExternalReference: "test",
			},
		},
		{
			name: "create a charge with invalid method",
			request: func() *cubawheeler.ChargeRequest {
				method := cubawheeler.ChargeMethod("INVALID")
				reference := "test"
				return &cubawheeler.ChargeRequest{
					Amount:      100,
					Currency:    "usd",
					Order:       "test",
					Method:      method,
					Reference:   &reference,
					Description: "test",
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := service.Create(ctx, tt.request())
			if (err != nil) != tt.wantErr {
				t.Errorf("ChargeService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ChargeService.Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
