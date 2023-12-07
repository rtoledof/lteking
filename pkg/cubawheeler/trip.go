package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Trip struct {
	ID            string               `json:"id" bson:"_id"`
	PickUp        *Location            `json:"pick_up" bson:"pick_up"`
	DropOff       *Location            `json:"drop_off" bson:"drop_off"`
	Route         []Location           `json:"route" bson:"route"`
	History       []Location           `json:"history,omitempty" bson:"history,omitempty"`
	Driver        string               `json:"driver,omitempty" bson:"driver,omitempty"`
	Rider         string               `json:"rider" bson:"rider"`
	Status        TripStatus           `json:"status" bson:"status"`
	StatusHistory []*TripStatusHistory `json:"status_history,omitempty" bson:"status_history,omitempty"`
	Rate          int                  `json:"rate" bson:"rate"`
	Price         int                  `json:"price" bson:"price"`
	Coupon        string               `json:"coupon,omitempty" bson:"coupon,omitempty"`
	StartAt       int                  `json:"start_at" bson:"start_at"`
	EndAt         int                  `json:"end_at" bson:"end_at"`
	Review        string               `json:"review,omitempty" bson:"review"`
	CreatedAt     int64                `json:"created_at" bson:"created_at"`
	UpdatedAt     int64                `json:"updated_at" bson:"updated_at"`
}

type TripList struct {
	Token string  `json:"token"`
	Data  []*Trip `json:"data"`
}

type UpdateTrip struct {
	Trip   string
	Driver string
	Status TripStatus
	Price  int
}

type TripFilter struct {
	Limit  *int      `json:"limit,omitempty"`
	Token  *string   `json:"token,omitempty"`
	Ids    []*string `json:"ids,omitempty"`
	Rider  *string   `json:"rider,omitempty"`
	Driver *string   `json:"driver,omitempty"`
	Status *string   `json:"status,omitempty"`
}

type RequestTrip struct {
	PickUp  *LocationInput   `json:"pick_up"`
	DropOff *LocationInput   `json:"drop_off"`
	Route   []*LocationInput `json:"route"`
	Hours   int              `json:"hours"`
	Min     int              `json:"min"`
	Sec     int              `json:"se"`
	Kms     float64          `json:"kms"`
}

type TripService interface {
	Create(context.Context, *RequestTrip) (*Trip, error)
	Update(context.Context, *UpdateTrip) (*Trip, error)
	FindByID(context.Context, string) (*Trip, error)
	FindAll(context.Context, *TripFilter) (*TripList, error)
}

type AddPlace struct {
	Name     string         `json:"name"`
	Location *LocationInput `json:"location"`
}

type TripStatusHistory struct {
	Status    TripStatus `json:"status" bson:"status"`
	ChangedAt string     `json:"changed_at" bson:"changed_at"`
}

type TripStatus string

const (
	TripStatusNew      TripStatus = "NEW"
	TripStatusPickUp   TripStatus = "PICK_UP"
	TripStatusOnTheWay TripStatus = "ON_THE_WAY"
	TripStatusDropOff  TripStatus = "DROP_OFF"
)

var AllTripStatus = []TripStatus{
	TripStatusNew,
	TripStatusPickUp,
	TripStatusOnTheWay,
	TripStatusDropOff,
}

func (e TripStatus) IsValid() bool {
	switch e {
	case TripStatusNew, TripStatusPickUp, TripStatusOnTheWay, TripStatusDropOff:
		return true
	}
	return false
}

func (e TripStatus) String() string {
	return string(e)
}

func (e *TripStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TripStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TripStatus", str)
	}
	return nil
}

func (e TripStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
