package cubawheeler

import (
	"context"
)

type Profile struct {
	ID       string `json:"id" faker:"-" bson:"_id"`
	Name     string `json:"name,omitempty" faker:"name" bson:"name"`
	LastName string `json:"last_name,omitempty" faker:"last_name" bson:"last_name"`
	DOB      string `json:"dob,omitempty" bson:"dob"`
	Phone    string `json:"phone,omitempty" faker:"phone_number" bson:"phone"`
	Photo    string `json:"photo,omitempty" faker:"url" bson:"photo"`
	Gender   Gender `json:"gender" bson:"gender"`
	Licence  string `json:"licence,omitempty" bson:"licence"`
	Dni      string `json:"dni,omitempty" bson:"dni"`
	UserId   string `faker:"-" bson:"user_id" bson:"user_id"`
}

type UpdateProfile struct {
	Name     *string `json:"name"`
	LastName *string `json:"last_name"`
	DOB      *string `json:"dob"`
	Phone    *string `json:"phone"`
	Photo    *string `json:"photo"`
	Gender   *Gender `json:"gender"`
	Licence  *string `json:"licence"`
	Dni      *string `json:"dni"`
}

type ProfileFilter struct {
	Limit   int
	Token   string
	IDs     []string
	Dni     string
	Licence string
	Gender  Gender
	User    string
}

type ProfileRequest struct {
	Name     *string
	LastName *string
	DOB      *string
	Pin      *string
	Phone    *string
	Photo    *string
	Gender   *Gender
	Licence  *string
	Dni      *string
}

type ProfileService interface {
	Create(context.Context, *ProfileRequest) (*Profile, error)
	Update(context.Context, *ProfileRequest) (*Profile, error)
	FindByUser(context.Context) (*Profile, error)
	FindAll(context.Context, *ProfileFilter) ([]*Profile, string, error)
	ChangePin(ctx context.Context, old *string, pin string) (*Profile, error)
}
