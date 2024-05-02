package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

const ExpireIn = time.Hour * 24 * 60 * 60

type Device struct {
	ID     primitive.ObjectID `bson:"_id"`
	Token  string             `bson:"token"`
	Name   string             `bson:"name,omitempty"`
	Active bool               `bson:"active"`
}

type DeviceFilter struct {
	Name   []string
	Active *bool
	Limit  int
	Token  string
}

type FavoriteVehicle struct {
	Plate string `json:"plate" bson:"plate"`
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
}

type User struct {
	ID                  string            `json:"id" faker:"-" bson:"_id"`
	Name                string            `json:"name,omitempty" faker:"name" bson:"name"`
	Password            []byte            `json:"-" bson:"password,omitempty"`
	Email               string            `json:"email" faker:"email" bson:"email"`
	Pin                 []byte            `json:"-" faker:"number" bson:"pin,omitempty"`
	Otp                 string            `json:"-" bson:"otp,omitempty"`
	Rate                float64           `json:"rate,omitempty" bson:"rate,omitempty"`
	Available           bool              `json:"-" bson:"available,omitempty"`
	Status              UserStatus        `json:"status" bson:"status"`
	ActiveVehicle       string            `json:"active_vehicle,omitempty" bson:"active_vehicle"`
	Referer             string            `json:"refer,omitempty" bson:"referer,omitempty"`
	Referal             string            `json:"referal,omitempty" bson:"referal,omitempty"`
	Role                Role              `json:"role" bson:"role"`
	FavoriteVehicles    []FavoriteVehicle `json:"favorite_vehicles,omitempty" bson:"favorite_vehicles,omitempty"`
	Profile             *Profile          `json:"profile,omitempty" faker:"-" bson:"profile"`
	BeansToken          map[string]any    `json:"beans_token,omitempty" bson:"beans_token,omitempty"`
	BeansTokenCreatedAt int64             `json:"beans_token_created_at,omitempty" bson:"beans_token_created_at,omitempty"`
	Devices             []Device          `json:"devices,omitempty" bson:"devices,omitempty"`
}

func (u User) Claim() map[string]interface{} {
	claims := map[string]interface{}{
		"user": u,
	}
	jwtauth.SetExpiryIn(claims, ExpireIn)
	jwtauth.SetIssuedNow(claims)
	return claims
}

func (u *User) GetDevices() []string {
	devices := make([]string, len(u.Devices))
	addedDevices := make(map[string]bool)
	for i, device := range u.Devices {
		if _, ok := addedDevices[device.Token]; !ok {
			devices[i] = device.Token
			addedDevices[device.Token] = true
		}
	}
	return devices
}

func (u *User) AddDevice(token, name string) error {
	if u.HasDevice(token) {
		return NewInvalidParameter("token", token)
	}
	u.Devices = append(u.Devices, Device{
		Token:  token,
		Name:   name,
		Active: true,
	})
	return nil
}

func (u *User) HasDevice(token string) bool {
	for _, device := range u.Devices {
		if device.Token == token {
			return true
		}
	}
	return false
}

func (u *User) RemoveDevice(token string) error {
	for i, device := range u.Devices {
		if device.Token == token {
			u.Devices = append(u.Devices[:i], u.Devices[i+1:]...)
			return nil
		}
	}
	return NewInvalidParameter("token", token)
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

type UserManager interface {
	FindByID(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
	FindAll(context.Context, *UserFilter) (*UserList, error)
	Update(context.Context, *User) error
	Token(context.Context, *User) (string, error)
}

type DeviceManager interface {
	Remove(context.Context, string) error
	FindAll(context.Context) ([]Device, error)
	AddDevice(context.Context, Device) error
}

type Deleter interface {
	DeleteAccount(context.Context) error
}

type ProfileManager interface {
	Me(context.Context) (*User, error)
	Update(context.Context, *Profile) error
	SetPin(context.Context, string) error
	ChangePin(context.Context, string, string) error
	SetDefaultCurrency(context.Context, string) error
}

type Reset interface {
	Profile(context.Context) (*Profile, error)
	Update(context.Context, *Profile) error
}

type Register interface {
	Register(context.Context, string, string, string) error
}

type LoginManager interface {
	Login(context.Context, string, string, ...string) (Token, error)
	Logout(context.Context) error
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
