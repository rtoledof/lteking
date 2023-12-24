package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Ads struct {
	ID          string      `json:"id" bson:"_id"`
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	Photo       string      `json:"photo" bson:"photo"`
	Owner       string      `json:"-" bson:"owner"`
	Inpression  *Impression `json:"inpression,omitempty" bson:"inpression"`
	Status      AdsStatus   `json:"status" bson:"status"`
	Priority    int         `json:"priority" bson:"priority"`
	ValidFrom   int         `json:"valid_from" bson:"valid_from"`
	ValidUntil  int         `json:"valid_until" bson:"valid_until"`
	Client      string      `json:"-" bson:"client"`
}

type AdsRequest struct {
	ID          string
	Limit       int
	Token       string
	Ids         []string
	Name        string
	Description string
	Photo       string
	Owner       string
	Status      AdsStatus
	Priority    int
	ValidFrom   int
	ValidUntil  int
}

type AdsService interface {
	Create(context.Context, *AdsRequest) (*Ads, error)
	Update(context.Context, *AdsRequest) (*Ads, error)
	FindById(context.Context, string) (*Ads, error)
	FindAll(context.Context, *AdsRequest) ([]*Ads, string, error)
}

type Impression string

const (
	ImpressionClick Impression = "CLICK"
)

var AllImpression = []Impression{
	ImpressionClick,
}

func (e Impression) IsValid() bool {
	switch e {
	case ImpressionClick:
		return true
	}
	return false
}

func (e Impression) String() string {
	return string(e)
}

func (e *Impression) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Impression(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Impression", str)
	}
	return nil
}

func (e Impression) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type AdsStatus string

const (
	AdsStatusNew       AdsStatus = "NEW"
	AdsStatusActive    AdsStatus = "ACTIVE"
	AdsStatusInactive  AdsStatus = "INACTIVE"
	AdsStatusSuspended AdsStatus = "SUSPENDED"
)

var AllAdsStatus = []AdsStatus{
	AdsStatusNew,
	AdsStatusActive,
	AdsStatusInactive,
	AdsStatusSuspended,
}

func (e AdsStatus) IsValid() bool {
	switch e {
	case AdsStatusNew, AdsStatusActive, AdsStatusInactive, AdsStatusSuspended:
		return true
	}
	return false
}

func (e AdsStatus) String() string {
	return string(e)
}

func (e *AdsStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AdsStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AdsStatus", str)
	}
	return nil
}

func (e AdsStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
