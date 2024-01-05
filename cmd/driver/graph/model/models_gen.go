// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type AddVehicle struct {
	Plate             string            `json:"plate"`
	Category          *VehicleCategory  `json:"category,omitempty"`
	Brand             string            `json:"brand"`
	Year              int               `json:"year"`
	Model             string            `json:"model"`
	Seats             int               `json:"seats"`
	Status            *VehicleStatus    `json:"status,omitempty"`
	Color             []string          `json:"color,omitempty"`
	Type              *VehicleType      `json:"type,omitempty"`
	Photos            []*graphql.Upload `json:"photos,omitempty"`
	Facilities        []*Facilities     `json:"facilities,omitempty"`
	OperativeLicense  *graphql.Upload   `json:"operative_license,omitempty"`
	TechnicInspection *graphql.Upload   `json:"technic_inspection,omitempty"`
}

type Amount struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryPrice struct {
	Category VehicleCategory `json:"category"`
	Price    *Amount         `json:"price"`
}

type Device struct {
	Token  string `json:"token"`
	Active bool   `json:"active"`
}

type Item struct {
	Points   []*Point `json:"points"`
	Riders   *int     `json:"riders,omitempty"`
	Baggages *bool    `json:"baggages,omitempty"`
	Currency *string  `json:"currency,omitempty"`
}

type Model struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Brand *Brand `json:"brand"`
}

type Order struct {
	ID            string           `json:"id"`
	Rate          int              `json:"rate"`
	Price         *Amount          `json:"price"`
	Rider         string           `json:"rider"`
	Driver        *string          `json:"driver,omitempty"`
	Status        string           `json:"status"`
	StatusHistory []*string        `json:"status_history,omitempty"`
	History       []*Point         `json:"history,omitempty"`
	Coupon        *string          `json:"coupon,omitempty"`
	StartAt       int              `json:"start_at"`
	EndAt         int              `json:"end_at"`
	Item          *Item            `json:"item"`
	Cost          []*CategoryPrice `json:"cost,omitempty"`
	SelectedCost  *CategoryPrice   `json:"selected_cost,omitempty"`
	RouteString   *string          `json:"route_string,omitempty"`
}

type Payment struct {
	ID        string `json:"id"`
	Order     *Order `json:"order"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Plan struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Recurrintg bool     `json:"recurrintg"`
	Orders     int      `json:"orders"`
	Price      int      `json:"price"`
	Interval   Interval `json:"interval"`
	Code       string   `json:"code"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Profile struct {
	ID               string         `json:"id"`
	Name             *string        `json:"name,omitempty"`
	LastName         *string        `json:"last_name,omitempty"`
	Dob              *string        `json:"dob,omitempty"`
	Phone            string         `json:"phone"`
	Photo            string         `json:"photo"`
	Gender           *Gender        `json:"gender,omitempty"`
	Licence          *string        `json:"licence,omitempty"`
	Circulation      *string        `json:"circulation,omitempty"`
	Dni              *string        `json:"dni,omitempty"`
	User             *User          `json:"user"`
	Status           *ProfileStatus `json:"status,omitempty"`
	PreferedCurrency *string        `json:"prefered_currency,omitempty"`
}

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type User struct {
	ID            string      `json:"id"`
	Email         string      `json:"email"`
	Password      *string     `json:"password,omitempty"`
	Pin           string      `json:"pin"`
	Otp           string      `json:"otp"`
	Rate          *float64    `json:"rate,omitempty"`
	Available     *bool       `json:"available,omitempty"`
	Status        *UserStatus `json:"status,omitempty"`
	ActiveVehicle *Vehicle    `json:"active_vehicle,omitempty"`
	Code          string      `json:"code"`
	Referer       string      `json:"referer"`
	Role          Role        `json:"role"`
	Profile       *Profile    `json:"profile"`
	Plan          *Plan       `json:"plan,omitempty"`
	Vehicles      []*Vehicle  `json:"vehicles,omitempty"`
	Orders        []*Order    `json:"orders,omitempty"`
	Devices       []*Device   `json:"devices,omitempty"`
}

type Vehicle struct {
	ID                string          `json:"id"`
	Plate             *string         `json:"plate,omitempty"`
	Category          VehicleCategory `json:"category"`
	Brand             *Brand          `json:"brand"`
	Year              int             `json:"year"`
	Model             *Model          `json:"model"`
	Seats             int             `json:"seats"`
	Status            VehicleStatus   `json:"status"`
	Color             []string        `json:"color,omitempty"`
	Type              VehicleType     `json:"type"`
	Photos            []*string       `json:"photos,omitempty"`
	Facilities        []*Facilities   `json:"facilities,omitempty"`
	TechnicInspection *string         `json:"technic_inspection,omitempty"`
}

type Facilities string

const (
	FacilitiesAirConditioner Facilities = "AIR_CONDITIONER"
	FacilitiesPetsAllowed    Facilities = "PETS_ALLOWED"
	FacilitiesSmokeAllowed   Facilities = "SMOKE_ALLOWED"
)

var AllFacilities = []Facilities{
	FacilitiesAirConditioner,
	FacilitiesPetsAllowed,
	FacilitiesSmokeAllowed,
}

func (e Facilities) IsValid() bool {
	switch e {
	case FacilitiesAirConditioner, FacilitiesPetsAllowed, FacilitiesSmokeAllowed:
		return true
	}
	return false
}

func (e Facilities) String() string {
	return string(e)
}

func (e *Facilities) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Facilities(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Facilities", str)
	}
	return nil
}

func (e Facilities) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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

type Interval string

const (
	IntervalDay   Interval = "DAY"
	IntervalWeek  Interval = "WEEK"
	IntervalMonth Interval = "MONTH"
	IntervalYear  Interval = "YEAR"
)

var AllInterval = []Interval{
	IntervalDay,
	IntervalWeek,
	IntervalMonth,
	IntervalYear,
}

func (e Interval) IsValid() bool {
	switch e {
	case IntervalDay, IntervalWeek, IntervalMonth, IntervalYear:
		return true
	}
	return false
}

func (e Interval) String() string {
	return string(e)
}

func (e *Interval) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Interval(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Interval", str)
	}
	return nil
}

func (e Interval) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProfileStatus string

const (
	ProfileStatusIncompleted ProfileStatus = "INCOMPLETED"
	ProfileStatusOnReview    ProfileStatus = "ON_REVIEW"
	ProfileStatusCompleted   ProfileStatus = "COMPLETED"
)

var AllProfileStatus = []ProfileStatus{
	ProfileStatusIncompleted,
	ProfileStatusOnReview,
	ProfileStatusCompleted,
}

func (e ProfileStatus) IsValid() bool {
	switch e {
	case ProfileStatusIncompleted, ProfileStatusOnReview, ProfileStatusCompleted:
		return true
	}
	return false
}

func (e ProfileStatus) String() string {
	return string(e)
}

func (e *ProfileStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProfileStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProfileStatus", str)
	}
	return nil
}

func (e ProfileStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleRider Role = "RIDER"
)

var AllRole = []Role{
	RoleRider,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleRider:
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
	UserStatusOnReview  UserStatus = "ON_REVIEW"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

var AllUserStatus = []UserStatus{
	UserStatusActive,
	UserStatusInactive,
	UserStatusOff,
	UserStatusOnReview,
	UserStatusSuspended,
}

func (e UserStatus) IsValid() bool {
	switch e {
	case UserStatusActive, UserStatusInactive, UserStatusOff, UserStatusOnReview, UserStatusSuspended:
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

type VehicleCategory string

const (
	VehicleCategoryX        VehicleCategory = "X"
	VehicleCategoryXl       VehicleCategory = "XL"
	VehicleCategoryConfort  VehicleCategory = "CONFORT"
	VehicleCategoryGreen    VehicleCategory = "GREEN"
	VehicleCategoryPets     VehicleCategory = "PETS"
	VehicleCategoryPackage  VehicleCategory = "PACKAGE"
	VehicleCategoryPriority VehicleCategory = "PRIORITY"
)

var AllVehicleCategory = []VehicleCategory{
	VehicleCategoryX,
	VehicleCategoryXl,
	VehicleCategoryConfort,
	VehicleCategoryGreen,
	VehicleCategoryPets,
	VehicleCategoryPackage,
	VehicleCategoryPriority,
}

func (e VehicleCategory) IsValid() bool {
	switch e {
	case VehicleCategoryX, VehicleCategoryXl, VehicleCategoryConfort, VehicleCategoryGreen, VehicleCategoryPets, VehicleCategoryPackage, VehicleCategoryPriority:
		return true
	}
	return false
}

func (e VehicleCategory) String() string {
	return string(e)
}

func (e *VehicleCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VehicleCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VehicleCategory", str)
	}
	return nil
}

func (e VehicleCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type VehicleStatus string

const (
	VehicleStatusNew       VehicleStatus = "NEW"
	VehicleStatusActive    VehicleStatus = "ACTIVE"
	VehicleStatusInactive  VehicleStatus = "INACTIVE"
	VehicleStatusSuspended VehicleStatus = "SUSPENDED"
)

var AllVehicleStatus = []VehicleStatus{
	VehicleStatusNew,
	VehicleStatusActive,
	VehicleStatusInactive,
	VehicleStatusSuspended,
}

func (e VehicleStatus) IsValid() bool {
	switch e {
	case VehicleStatusNew, VehicleStatusActive, VehicleStatusInactive, VehicleStatusSuspended:
		return true
	}
	return false
}

func (e VehicleStatus) String() string {
	return string(e)
}

func (e *VehicleStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VehicleStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VehicleStatus", str)
	}
	return nil
}

func (e VehicleStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type VehicleType string

const (
	VehicleTypeAuto VehicleType = "AUTO"
	VehicleTypeMoto VehicleType = "MOTO"
	VehicleTypeBike VehicleType = "BIKE"
)

var AllVehicleType = []VehicleType{
	VehicleTypeAuto,
	VehicleTypeMoto,
	VehicleTypeBike,
}

func (e VehicleType) IsValid() bool {
	switch e {
	case VehicleTypeAuto, VehicleTypeMoto, VehicleTypeBike:
		return true
	}
	return false
}

func (e VehicleType) String() string {
	return string(e)
}

func (e *VehicleType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VehicleType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VehicleType", str)
	}
	return nil
}

func (e VehicleType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
