package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Plan struct {
	ID         string   `json:"id" bson:"_id"`
	Name       string   `json:"name" bson:"name"`
	Recurrintg bool     `json:"recurrintg" bson:"recurrintg"`
	Trips      int      `json:"trips" bson:"trips"`
	Price      int      `json:"price" bson:"price"`
	Interval   Interval `json:"interval" bson:"interval"`
	Code       string   `json:"code" bson:"code"`
}

type PlanRequest struct {
	ID         string
	Name       *string
	Recurring  *bool
	TotalTrips *int
	Price      *int
	Interval   *Interval
	Code       *string
}

type PlanFilter struct {
	Ids        []string
	Limit      int
	Token      string
	Name       string
	Recurring  bool
	TotalTrips int
	Price      int
	Interval   Interval
	Code       string
}

type PlanService interface {
	Create(context.Context, *PlanRequest) (*Plan, error)
	Update(context.Context, *PlanRequest) (*Plan, error)
	FindByID(context.Context, string) (*Plan, error)
	FindAll(context.Context, *PlanFilter) ([]*Plan, string, error)
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
