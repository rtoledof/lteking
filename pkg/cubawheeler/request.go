package cubawheeler

import (
	"fmt"
	"io"
	"strconv"
)

type LoginRequest struct {
	Email        string    `json:"email"`
	Otp          string    `json:"otp"`
	GrantType    GrantType `json:"grant_type"`
	Client       string    `json:"client"`
	Secret       string    `json:"secret"`
	Referer      string    `json:"referer"`
	RefreshToken string    `json:"refresh_token"`
}

type GrantType string

const (
	GrantTypePassword GrantType = "password"
	GrantTypeRefresh  GrantType = "refresh_token"
	GrantTypeClient   GrantType = "client_credentials"
	GrantTypeAuthCode GrantType = "authorization_code"
)

func (e GrantType) IsValid() bool {
	switch e {
	case GrantTypePassword, GrantTypeAuthCode, GrantTypeClient, GrantTypeRefresh:
		return true
	}
	return false
}

func (e GrantType) String() string {
	return string(e)
}

func (e *GrantType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GrantType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

func (e GrantType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
