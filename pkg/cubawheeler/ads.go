package cubawheeler

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type Ads struct {
	gorm.Model
	ID          string      `json:"id" gorm:"primaryKey;varchar(36);not null"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Photo       string      `json:"photo"`
	ClientID    string      `json:"owner"`
	Owner       string      `json:"-" gorm:"foreignKey:ClientID"`
	Inpression  *Impression `json:"inpression,omitempty"`
	Status      AdsStatus   `json:"status"`
	Priority    int         `json:"priority"`
	ValidFrom   int         `json:"valid_from"`
	ValidUntil  int         `json:"valid_until"`
}

func (a *Ads) BeforeSave(*gorm.DB) error {
	if a.ID == "" {
		a.ID = NewID().String()
	}
	return nil
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
