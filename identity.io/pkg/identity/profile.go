package identity

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Profile struct {
	ID                string        `json:"-" faker:"-" bson:"_id"`
	Name              string        `json:"name,omitempty" faker:"name" bson:"name"`
	LastName          string        `json:"last_name,omitempty" faker:"last_name" bson:"last_name"`
	DOB               string        `json:"dob,omitempty" bson:"dob"`
	Phone             string        `json:"phone,omitempty" faker:"phone_number" bson:"phone"`
	Photo             string        `json:"photo,omitempty" faker:"url" bson:"photo"`
	Gender            Gender        `json:"gender,omitempty" bson:"gender"`
	Licence           string        `json:"-" bson:"licence"`
	Dni               string        `json:"-" bson:"dni"`
	UserId            string        `faker:"-" json:"-" bson:"user_id"`
	Status            ProfileStatus `json:"-" bson:"status"`
	Circulation       string        `json:"-" bson:"circulation"`
	TechnicInspection string        `json:"-" bson:"technic_inspection"`
	PreferedCurrency  string        `json:"prefered_currency,omitempty" bson:"prefered_currency"`
}

func (p *Profile) IsCompleted(role Role) bool {
	switch role {
	case RoleRider:
		return len(p.Name) > 0 && len(p.LastName) > 0 && len(p.DOB) > 0
	case RoleDriver:
		return len(p.Name) > 0 && len(p.LastName) > 0 &&
			len(p.DOB) > 0 && len(p.Licence) > 0 && len(p.Circulation) > 0 &&
			len(p.TechnicInspection) > 0
	default:
		return true
	}
}

type UpdateProfile struct {
	Name              string `json:"name"`
	LastName          string `json:"last_name"`
	Dob               string `json:"dob"`
	Phone             string `json:"phone"`
	Photo             string `json:"photo"`
	Gender            Gender `json:"gender"`
	Licence           string `json:"licence"`
	Dni               string `json:"dni"`
	Circulation       string `json:"circulation,omitempty"`
	TechnicInspection string `json:"technic_inspection,omitempty"`
	PreferedCurrency  string `json:"prefered_currency,omitempty"`
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

type ProfileStatus string

const (
	ProfileStatusIncompleted ProfileStatus = "INCOMPLETED"
	ProfileStatusOnReview    ProfileStatus = "ON_REVIEW"
	ProfileStatusCompleted   ProfileStatus = "COMPLETED"
)

var AllProfileStatus = []ProfileStatus{
	ProfileStatusIncompleted,
	ProfileStatusOnReview,
	ProfileStatusCompleted,
}

func (e ProfileStatus) IsValid() bool {
	switch e {
	case ProfileStatusIncompleted, ProfileStatusOnReview, ProfileStatusCompleted:
		return true
	}
	return false
}

func (e ProfileStatus) String() string {
	return string(e)
}

func (e *ProfileStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProfileStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProfileStatus", str)
	}
	return nil
}

func (e ProfileStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProfileService interface {
	Create(context.Context, *UpdateProfile) (*Profile, error)
	Update(context.Context, *UpdateProfile) (*Profile, error)
	FindByUser(context.Context) (*Profile, error)
	FindAll(context.Context, *ProfileFilter) ([]*Profile, string, error)
	ChangePin(ctx context.Context, old *string, pin string) (*Profile, error)
}