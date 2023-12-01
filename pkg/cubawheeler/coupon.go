package cubawheeler

import (
	"fmt"
	"io"
	"strconv"
)

type Coupon struct {
	ID         string       `json:"id"`
	Code       string       `json:"code"`
	Percent    *float64     `json:"percent,omitempty"`
	Amount     *int         `json:"amount,omitempty"`
	Status     CouponStatus `json:"status"`
	ValidFrom  *int         `json:"valid_from,omitempty"`
	ValidUntil *int         `json:"valid_until,omitempty"`
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
