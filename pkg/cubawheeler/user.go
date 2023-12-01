package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID               string     `json:"id" gorm:"type:varchar(36);primaryKey" faker:"-"`
	Name             string     `json:"name" gorm:"type:varchar(100);not null" faker:"name"`
	Password         []byte     `json:"-"`
	Email            string     `json:"email" gorm:"type:varchar(100);not null;unique" faker:"email"`
	Pin              string     `json:"pin" gotm:"type:varchar(4)" faker:"number"`
	OTP              string     `json:"-"`
	Rate             float64    `json:"rate"`
	Available        bool       `json:"-"`
	Status           UserStatus `json:"status"`
	ActiveVehicle    string     `json:"active_vehicle,omitempty"`
	Code             string     `json:"referal_code"`
	Referer          string     `json:"-"`
	Role             Role       `json:"-"`
	Plan             string     `json:"plan,omitempty"`
	Locations        []Location `json:"locations,omitempty" gorm:"-"`
	Vehicles         []Vehicle  `json:"vehicles,omitempty" gorm:"-"`
	FavoriteVehicles []Vehicle  `json:"favorite_vehicles,omitempty" `
	Trips            []Trip     `json:"trips,omitempty" gorm:"-"`
	Profile          Profile    `json:"profile" faker:-`
}

func (u *User) BeforeSave(*gorm.DB) error {
	if u.ID == "" {
		u.ID = NewID().String()
	}
	return nil
}

func (u *User) EncryptPassword(password string) error {
	var err error
	u.Password, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}

func (u *User) GenToken() (*Token, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Id:        u.ID,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "cubawheeler",
	})

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: time.Duration(time.Hour * 24 * 7),
		AccessTokenCreatedAt: time.Now(),
	}, nil
}

type UserService interface {
	FindByID(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
	FindAll(context.Context, UserFilter) ([]*User, string, error)
	UpdateOTP(context.Context, string, uint64) error
	CreateUser(context.Context, *User) error
}

type OTPServer interface {
	New(context.Context, string) (int, error)
}

type UserFilter struct {
	Limit uint
	Token string
	Ids   []string
	Name  string
	Email string
	OTP   int
}

type Gender string

const (
	GenderFemale     Gender = "FEMALE"
	GenderMale       Gender = "MALE"
	GenderNotDefined Gender = "NOT_DEFINED"
)

var AllGender = []Gender{
	GenderFemale,
	GenderMale,
	GenderNotDefined,
}

func (e Gender) IsValid() bool {
	switch e {
	case GenderFemale, GenderMale, GenderNotDefined:
		return true
	}
	return false
}

func (e Gender) String() string {
	return string(e)
}

func (e *Gender) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Gender(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

func (e Gender) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleRider  Role = "RIDER"
	RoleDriver Role = "DRIVER"
	RoleSale   Role = "SALE"
	RoleAdmin  Role = "ADMIN"
	RoleClient Role = "CLIENT"
)

var AllRole = []Role{
	RoleRider,
	RoleDriver,
	RoleSale,
	RoleAdmin,
	RoleClient,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleRider, RoleDriver, RoleSale, RoleAdmin, RoleClient:
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

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusOff       UserStatus = "OFF"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

var AllUserStatus = []UserStatus{
	UserStatusActive,
	UserStatusInactive,
	UserStatusOff,
	UserStatusSuspended,
}

func (e UserStatus) IsValid() bool {
	switch e {
	case UserStatusActive, UserStatusInactive, UserStatusOff, UserStatusSuspended:
		return true
	}
	return false
}

func (e UserStatus) String() string {
	return string(e)
}

func (e *UserStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserStatus", str)
	}
	return nil
}

func (e UserStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
