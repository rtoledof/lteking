package cubawheeler

import "context"

type Rate struct {
	ID                string `json:"id" bson:"_id"`
	Code              string `json:"code" bson:"code"`
	BasePrice         int    `json:"base_price" bson:"base_price"`
	PricePerMin       int    `json:"price_per_min" bson:"price_per_min"`
	PricePerKm        int    `json:"price_per_km" bson:"price_per_km"`
	PricePerPassenger *int   `json:"price_per_passenger,omitempty" bson:"price_per_passenger,omitempty"`
	PricePerBaggage   int    `json:"price_per_baggage" bson:"price_per_baggage"`
	StartTime         *int   `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime           *int   `json:"end_time,omitempty" bson:"end_time,omitempty"`
	StartDate         *int   `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate           *int   `json:"end_date,omitempty" bson:"end_date,omitempty"`
	MinKm             *int   `json:"min_km,omitempty" bson:"min_km,omitempty"`
	MaxKm             *int   `json:"max_km,omitempty" bson:"max_km,omitempty"`
}

type RateRequest struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	BasePrice         int    `json:"base_price"`
	PricePerMin       *int   `json:"price_per_min,omitempty"`
	PricePerKm        *int   `json:"price_per_km,omitempty"`
	PricePerPassenger *int   `json:"price_per_passenger,omitempty"`
	PricePerBaggage   *int   `json:"price_per_baggage,omitempty"`
	StartTime         *int   `json:"start_time,omitempty"`
	EndTime           *int   `json:"end_time,omitempty"`
	StartDate         *int   `json:"start_date,omitempty"`
	EndDate           *int   `json:"end_date,omitempty"`
	MinKm             *int   `json:"min_km,omitempty"`
	MaxKm             *int   `json:"max_km,omitempty"`
}

type RateService interface {
	Create(context.Context, *RateRequest) (*Rate, error)
	Update(context.Context, *RateRequest) (*Rate, error)
	FindID(context.Context, string) (*Rate, error)
	FindAll(context.Context, *RateRequest) ([]Rate, string, error)
}
