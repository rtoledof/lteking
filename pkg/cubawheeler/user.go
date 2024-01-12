package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Device struct {
	Token  string `json:"token" bson:"id"`
	Active bool   `json:"active" bson:"active"`
}

type User struct {
	ID                  string         `json:"id" faker:"-" bson:"_id"`
	Name                string         `json:"name,omitempty" faker:"name" bson:"name"`
	Password            []byte         `json:"-" bson:"password,omitempty"`
	Email               string         `json:"email" faker:"email" bson:"email"`
	Pin                 []byte         `json:"-" faker:"number" bson:"pin,omitempty"`
	Otp                 string         `json:"-" bson:"otp,omitempty"`
	Rate                float64        `json:"rate,omitempty" bson:"rate,omitempty"`
	Available           bool           `json:"-" bson:"available,omitempty"`
	Status              UserStatus     `json:"status" bson:"status"`
	ActiveVehicle       string         `json:"active_vehicle,omitempty" bson:"active_vehicle"`
	Referer             string         `json:"refer" bson:"referer,omitempty"`
	Referal             string         `json:"referal,omitempty" bson:"referal,omitempty"`
	Role                Role           `json:"role" bson:"role"`
	Plan                string         `json:"plan,omitempty" bson:"plan"`
	Locations           []*Location    `json:"locations,omitempty" bson:"locations,omitempty"`
	LastLocations       []*Location    `json:"last_locations,omitempty" bson:"last_locations,omitempty"`
	Vehicles            []*Vehicle     `json:"vehicles,omitempty" bson:"vehicles"`
	FavoriteVehicles    []*Vehicle     `json:"favorite_vehicles,omitempty" bson:"favorite_vehicles,omitempty"`
	Orders              []*Order       `json:"orders,omitempty" bson:"orders,omitempty"`
	Profile             Profile        `json:"profile,omitempty" faker:"-" bson:"profile"`
	BeansToken          map[string]any `json:"beans_token,omitempty" bson:"beans_token,omitempty"`
	BeansTokenCreatedAt int64          `json:"beans_token_created_at,omitempty" bson:"beans_token_created_at,omitempty"`
	Devices             []Device       `json:"devices,omitempty" bson:"devices,omitempty"`
	jwt.StandardClaims
}

func (u *User) HasVehicle(id string) bool {
	for _, v := range u.Vehicles {
		if v.ID == id || v.Plate == id {
			return true
		}
	}
	return false
}

type UserList struct {
	Token string  `json:"token"`
	Data  []*User `json:"data"`
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive || u.Status == UserStatusOnReview
}

func (u *User) EncryptPassword(password string) error {
	var err error
	u.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EncryptPin(pin string) error {
	var err error
	u.Pin, err = bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}

func (u *User) ComparePin(pin string) error {
	return bcrypt.CompareHashAndPassword(u.Pin, []byte(pin))
}

type UserService interface {
	FindByID(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
	FindAll(context.Context, *UserFilter) (*UserList, error)
	CreateUser(context.Context, *User) error
	Me(context.Context) (*Profile, error)
	AddFavoritePlace(context.Context, AddPlace) (*Location, error)
	FavoritePlaces(context.Context) ([]*Location, error)
	Orders(context.Context, *OrderFilter) (*OrderList, error)
	LastNAddress(context.Context, int) ([]*Location, error)
	Login(context.Context, LoginRequest) (*User, error)
	AddFavoriteVehicle(context.Context, *string) (*Vehicle, error)
	FavoriteVehicles(context.Context) ([]*Vehicle, error)
	UpdatePlace(context.Context, *UpdatePlace) (*Location, error)
	UpdateProfile(context.Context, *UpdateProfile) error
	AddDevice(context.Context, string) error
	GetUserDevices(context.Context, []string) ([]string, error)
	SetAvailability(ctx context.Context, user string, available bool) error
	Update(context.Context, *User) error
}

type OtpService interface {
	Create(context.Context, string) (string, error)
	Otp(context.Context, string, string) error
}

type OTPServer interface {
	New(context.Context, string) (int, error)
}

type UserFilter struct {
	Limit  int
	Token  string
	Ids    []string
	Name   string
	Email  string
	Otp    string
	Pin    string
	User   string
	Role   Role
	Status []UserStatus
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
	UserStatusOnReview  UserStatus = "ON_REVIEW"
)

var AllUserStatus = []UserStatus{
	UserStatusActive,
	UserStatusInactive,
	UserStatusOff,
	UserStatusSuspended,
	UserStatusOnReview,
}

func (e UserStatus) IsValid() bool {
	switch e {
	case UserStatusActive, UserStatusInactive, UserStatusOff, UserStatusSuspended, UserStatusOnReview:
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
