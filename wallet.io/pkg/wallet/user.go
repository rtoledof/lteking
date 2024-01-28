package wallet

import (
	"fmt"
	"io"
	"strconv"
)

type Role string

const (
	RoleRider  Role = "RIDER"
	RoleDriver Role = "DRIVER"
	RoleAdmin  Role = "ADMIN"
)

var AllRole = []Role{
	RoleRider,
	RoleDriver,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleRider, RoleDriver:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type User struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	LastName string `json:"last_name" bson:"last_name"`
	Email    string `json:"email" bson:"email"`
	Role     Role   `json:"role" bson:"role"`
}
