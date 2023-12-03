package cubawheeler

import "context"

type Review struct {
	ID      string  `json:"id" bson:"_id"`
	From    string  `json:"from" bson:"from"`
	To      string  `json:"to" bson:"to"`
	Comment *string `json:"comment,omitempty" bson:"comment,omitempty"`
	Rate    float64 `json:"rate" bson:"rate"`
}

type ReviewCreate struct {
	From    string
	To      string
	Comment string
	Rate    float64
}

type ReviewFilter struct {
	From string
	To   string
}

type ReviewService interface {
	Create(context.Context, *ReviewCreate) (*Review, error)
	Update(context.Context, *ReviewCreate) (*Review, error)
	FindById(context.Context, string) (*Review, error)
	FindAll(context.Context, ReviewFilter) ([]*Review, string, error)
}
