package cubawheeler

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type Plan struct {
	gorm.Model
	ID         string   `json:"id" gorm:"primaryKey;varchar(36);not null"`
	Name       string   `json:"name"`
	Recurrintg bool     `json:"recurrintg"`
	Trips      int      `json:"trips"`
	Price      int      `json:"price"`
	Interval   Interval `json:"interval"`
	Code       string   `json:"code"`
}

func (p *Plan) BeforeSave(*gorm.DB) error {
	if p.ID == "" {
		p.ID = NewID().String()
	}
	return nil
}

type Interval string

const (
	IntervalDay   Interval = "DAY"
	IntervalWeek  Interval = "WEEK"
	IntervalMonth Interval = "MONTH"
	IntervalYear  Interval = "YEAR"
)

var AllInterval = []Interval{
	IntervalDay,
	IntervalWeek,
	IntervalMonth,
	IntervalYear,
}

func (e Interval) IsValid() bool {
	switch e {
	case IntervalDay, IntervalWeek, IntervalMonth, IntervalYear:
		return true
	}
	return false
}

func (e Interval) String() string {
	return string(e)
}

func (e *Interval) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Interval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Interval", str)
	}
	return nil
}

func (e Interval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
