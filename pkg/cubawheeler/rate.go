package cubawheeler

type Rate struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	BasePrice         int    `json:"base_price"`
	PricePerMin       int    `json:"price_per_min"`
	PricePerKm        int    `json:"price_per_km"`
	PricePerPassenger *int   `json:"price_per_passenger,omitempty"`
	PricePerBaggage   int    `json:"price_per_baggage"`
	StartTime         *int   `json:"start_time,omitempty"`
	EndTime           *int   `json:"end_time,omitempty"`
	StartDate         *int   `json:"start_date,omitempty"`
	EndDate           *int   `json:"end_date,omitempty"`
	MinKm             *int   `json:"min_km,omitempty"`
	MaxKm             *int   `json:"max_km,omitempty"`
}
