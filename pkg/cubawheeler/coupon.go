package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Coupon struct {
	ID         string       `json:"id" bson:"_id"`
	Code       string       `json:"code" bson:"code"`
	Percent    *float64     `json:"percent,omitempty" bson:"percent,omitempty"`
	Amount     *int         `json:"amount,omitempty" bson:"amount,omitempty"`
	Status     CouponStatus `json:"status" bson:"status"`
	ValidFrom  int64        `json:"valid_from,omitempty" bson:"valid_from,omitempty"`
	ValidUntil int64        `json:"valid_until,omitempty" bson:"valid_until,omitempty"`
	CreatedAt  int64        `json:"-" bson:"created_at"`
	UpdatedAt  int64        `json:"updated_at" bson:"updated_at"`
}

type CouponRequest struct {
	ID         string
	Limit      int
	Token      string
	Ids        []string
	Code       string
	Percent    *float64
	Amount     *int
	Status     CouponStatus
	ValidFrom  *int64
	ValidUntil *int64
}

type CouponService interface {
	Create(context.Context, *CouponRequest) (*Coupon, error)
	FindByID(context.Context, string) (*Coupon, error)
	FindAll(context.Context, *CouponRequest) ([]*Coupon, string, error)
	FindByCode(context.Context, string) (*Coupon, error)
	Redeem(context.Context, string) (*Coupon, error)
}

type CouponStatus string

const (
	CouponStatusNew      CouponStatus = "NEW"
	CouponStatusActive   CouponStatus = "ACTIVE"
	CouponStatusInactive CouponStatus = "INACTIVE"
	CouponStatusRedeemed CouponStatus = "REDEEMED"
)

var AllCouponStatus = []CouponStatus{
	CouponStatusNew,
	CouponStatusActive,
	CouponStatusInactive,
	CouponStatusRedeemed,
}

func (e CouponStatus) IsValid() bool {
	switch e {
	case CouponStatusNew, CouponStatusActive, CouponStatusInactive, CouponStatusRedeemed:
		return true
	}
	return false
}

func (e CouponStatus) String() string {
	return string(e)
}

func (e *CouponStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CouponStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CouponStatus", str)
	}
	return nil
}

func (e CouponStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
