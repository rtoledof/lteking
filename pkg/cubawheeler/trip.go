package cubawheeler

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type Trip struct {
	gorm.Model
	ID              string               `json:"id" gotm:"privateKey;varchar(36);not null"`
	CurrentPosition *Location            `json:"current_position gorm:"-"`
	History         []Location           `json:"history,omitempty" gorm:"many2many:trip_history"`
	Driver          string               `json:"driver,omitempty"`
	Rider           string               `json:"rider"`
	Status          TripStatus           `json:"status"`
	StatusHistory   []*TripStatusHistory `json:"status_history,omitempty"`
	Rate            int                  `json:"rate"`
	Price           int                  `json:"price"`
	Coupon          string               `json:"coupon,omitempty"`
	StartAt         int                  `json:"start_at"`
	EndAt           int                  `json:"end_at"`
	Review          string               `json:"review,omitempty"`
	User            string               `json:"user"`
}

func (t *Trip) BeforeSave(*gorm.DB) error {
	if t.ID == "" {
		t.ID = NewID().String()
	}
	return nil
}

type TripStatusHistory struct {
	Status    TripStatus `json:"status"`
	ChangedAt string     `json:"changed_at"`
}

type TripStatus string

const (
	TripStatusPickUp   TripStatus = "PICK_UP"
	TripStatusOnTheWay TripStatus = "ON_THE_WAY"
	TripStatusDropOff  TripStatus = "DROP_OFF"
)

var AllTripStatus = []TripStatus{
	TripStatusPickUp,
	TripStatusOnTheWay,
	TripStatusDropOff,
}

func (e TripStatus) IsValid() bool {
	switch e {
	case TripStatusPickUp, TripStatusOnTheWay, TripStatusDropOff:
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
