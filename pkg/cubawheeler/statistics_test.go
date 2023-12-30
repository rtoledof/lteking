package cubawheeler

import (
	"testing"

	"cubawheeler.io/pkg/currency"
	"github.com/google/go-cmp/cmp"
)

func TestOrderStatisticsAddOrder(t *testing.T) {
	var tests = []struct {
		name  string
		order Order
		date  Time
		want  *OrderStatistics
	}{
		{
			name: "add order",
			order: Order{
				ID: "1",
				Price: &currency.Amount{
					Amount:   100,
					Currency: currency.MustParse("CUP"),
				},
			},
			date: Time{Now().Time},
			want: &OrderStatistics{
				Orders: map[string]*Statistics{
					Time{Now().Time}.Format("2006-01-02"): {
						Total: 1,
						Amount: currency.Amount{
							Amount:   100,
							Currency: currency.MustParse("CUP"),
						},
						Orders: []string{"1"},
					},
					Time{Now().Time}.Format("2006-01"): {
						Total: 1,
						Amount: currency.Amount{
							Amount:   100,
							Currency: currency.MustParse("CUP"),
						},
						Orders: []string{"1"},
					},
					Time{Now().Time}.Format("2006"): {
						Total: 1,
						Amount: currency.Amount{
							Amount:   100,
							Currency: currency.MustParse("CUP"),
						},
						Orders: []string{"1"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderStatistics{
				User: tt.order.Rider,
			}
			o.AddOrder(tt.order, tt.date)
			if diff := cmp.Diff(tt.want, o); diff != "" {
				t.Errorf("OrderStatistics.AddOrder() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
