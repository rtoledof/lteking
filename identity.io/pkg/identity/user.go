package identity

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

const ExpireIn = time.Hour * 24 * 60 * 60

type Device struct {
	Token  string `json:"token" bson:"id"`
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
	Active bool   `json:"active" bson:"active"`
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
	Role                Role              `json:"-" bson:"role"`
	Locations           []*Location       `json:"locations,omitempty" bson:"locations,omitempty"`
	LastLocations       []*Location       `json:"last_locations,omitempty" bson:"last_locations,omitempty"`
	Vehicles            []*Vehicle        `json:"vehicles,omitempty" bson:"vehicles,omitempty"`
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

func (u *User) AddVehicle(vehicle *Vehicle) error {
	if vehicle == nil {
		return NewInvalidParameter("vehicle", vehicle)
	}
	if !vehicle.IsValid() {
		return NewInvalidParameter("vehicle", vehicle)
	}
	if u.HasVehicle(vehicle.ID) {
		return NewInvalidParameter("vehicle", vehicle)
	}
	u.Vehicles = append(u.Vehicles, vehicle)
	return nil
}

func (u *User) UpdateVehicle(vehicle *Vehicle) error {
	if vehicle == nil {
		return NewInvalidParameter("vehicle", vehicle)
	}
	if !vehicle.IsValid() {
		return NewInvalidParameter("vehicle", vehicle)
	}
	for i, v := range u.Vehicles {
		if v.ID == vehicle.ID {
			u.Vehicles[i] = vehicle
			return nil
		}
	}
	return NewInvalidParameter("vehicle", vehicle)
}

func (u *User) DeleteVehicle(id string) error {
	for i, v := range u.Vehicles {
		if v.ID == id {
			u.Vehicles = append(u.Vehicles[:i], u.Vehicles[i+1:]...)
			return nil
		}
	}
	return NewInvalidParameter("vehicle", id)
}

func (u *User) GetActiveVehicle() *Vehicle {
	for _, v := range u.Vehicles {
		if v.ID == u.ActiveVehicle {
			return v
		}
	}
	return nil
}

func (u *User) GetVehicles() []*Vehicle {
	return u.Vehicles
}

func (u *User) HasVehicle(id string) bool {
	for _, v := range u.Vehicles {
		if v.ID == id || v.Plate == id {
			return true
		}
	}
	return false
}

func (u *User) HasFavoriteVehicle(plate string) bool {
	for _, v := range u.FavoriteVehicles {
		if v.Plate == plate {
			return true
		}
	}
	return false
}

func (u *User) AddFavoriteVehicle(plate, name string) {
	u.FavoriteVehicles = append(u.FavoriteVehicles, FavoriteVehicle{
		Plate: plate,
		Name:  name,
	})
}

func (u *User) DeleteFavoriteVehicle(plate string) {
	for i, v := range u.FavoriteVehicles {
		if v.Plate == plate {
			u.FavoriteVehicles = append(u.FavoriteVehicles[:i], u.FavoriteVehicles[i+1:]...)
			return
		}
	}
}

func (u *User) GetFavoriteVehicles() []string {
	vehicles := make([]string, len(u.FavoriteVehicles))
	for i, v := range u.FavoriteVehicles {
		vehicles[i] = v.Plate
	}
	return vehicles
}

func (u *User) SetActiveVehicle(id string) bool {
	for _, v := range u.Vehicles {
		if v.ID == id || v.Plate == id {
			u.ActiveVehicle = v.ID
			return true
		}
	}
	return false
}

func (u *User) GetFavoritePlace(id string) *Location {
	for _, l := range u.Locations {
		if l.ID == id {
			return l
		}
	}
	return nil
}

func (u *User) LastNAddress(n int) []*Location {
	if n > len(u.LastLocations) {
		n = len(u.LastLocations)
	}
	return u.LastLocations[:n]
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
	Me(context.Context) (*User, error)
	SetPreferedCurrency(context.Context, string) error

	AddDeviceToken(context.Context, string, string) error
	RemoveDeviceToken(context.Context, string) error
	DeviceTokens(context.Context) ([]string, error)

	//TODO: move this to another service
	AddFavoritePlace(context.Context, string, Point) (*Location, error)
	FavoritePlaces(context.Context) ([]*Location, error)
	DeleteFavoritePlace(context.Context, string) error
	UpdateFavoritePlace(context.Context, string, UpdatePlace) error
	FavoritePlace(context.Context, string) (*Location, error)

	LastNAddress(context.Context, int) ([]*Location, error)
	Login(context.Context, string, string, ...string) (*User, error)

	AddFavoriteVehicle(context.Context, string, *string) error
	FavoriteVehicles(context.Context) ([]string, error)
	DeleteFavoriteVehicle(context.Context, string) error

	UpdatePlace(context.Context, *UpdatePlace) (*Location, error)
	UpdateProfile(context.Context, *UpdateProfile) error
	AddDevice(context.Context, string) error
	GetUserDevices(context.Context, UserFilter) ([]string, error)
	SetAvailability(ctx context.Context, available bool) error
	Update(context.Context, *User) error
	Token(context.Context, *User) (string, error)
	Logout(context.Context) error
	// Driver: car handler
	AddVehicle(context.Context, *Vehicle) error
	UpdateVehicle(context.Context, *Vehicle) error
	DeleteVehicle(context.Context, string) error
	SetActiveVehicle(context.Context, string) error
	Vehicles(context.Context) ([]*Vehicle, error)
	Vehicle(context.Context, string) (*Vehicle, error)
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
