// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"cubawheeler.io/pkg/cubawheeler"
)

type LoginRequest struct {
	Email string  `json:"email"`
	Otp   *string `json:"otp,omitempty"`
	Pin   *string `json:"pin,omitempty"`
}

type TripStatusHistory struct {
	Status    cubawheeler.TripStatus `json:"status"`
	ChangedAt string                 `json:"changed_at"`
}

type UpdateProfile struct {
	Name *string `json:"name,omitempty"`
	Dob  *string `json:"dob,omitempty"`
}

type UpdateTrip struct {
	Driver *string                 `json:"driver,omitempty"`
	Status *cubawheeler.TripStatus `json:"status,omitempty"`
}

type UserFilter struct {
	Ids   []*string `json:"ids,omitempty"`
	Email *string   `json:"email,omitempty"`
	Token *string   `json:"token,omitempty"`
	Limit *int      `json:"limit,omitempty"`
	Name  *string   `json:"name,omitempty"`
}
