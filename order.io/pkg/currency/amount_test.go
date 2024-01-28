package currency

import (
	"testing"

	"golang.org/x/text/currency"
)

func TestAmountString(t *testing.T) {
	type fields struct {
		Amount   int64
		Currency Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test to string",
			fields: fields{
				Amount:   10000,
				Currency: Currency{Unit: currency.EUR},
			},
			want: "100.00",
		}, {
			name: "no decimals",
			fields: fields{
				Amount:   10000,
				Currency: Currency{Unit: currency.JPY},
			},
			want: "10000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Amount{
				Amount:   tt.fields.Amount,
				Currency: tt.fields.Currency,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("Amount.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAmountToMinAmount(t *testing.T) {
	var testCase = []struct {
		name string
		in   struct {
			amount float64
			cur    string
		}
		want int64
	}{
		{
			name: "convert USD to min amount",
			in: struct {
				amount float64
				cur    string
			}{amount: 10.00, cur: "USD"},
			want: 1000,
		}, {
			name: "convert JPY to min amount",
			in: struct {
				amount float64
				cur    string
			}{amount: 1000, cur: "JPY"},
			want: 1000,
		}, {
			name: "convert BTC to min amount",
			in: struct {
				amount float64
				cur    string
			}{amount: 1, cur: "BTC"},
			want: 100000000,
		},
	}

	for _, v := range testCase {
		t.Run(v.name, func(t *testing.T) {
			var amount Amount
			amount.ToMinAmount(v.in.amount, v.in.cur)
			if v.want != amount.Amount {
				t.Fatalf("got = %d, want %d", amount.Amount, v.want)
			}
		})
	}
}
