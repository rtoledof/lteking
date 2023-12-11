package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Point struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lon float64 `json:"lon" bson:"lon"`
}

type Address struct {
	Street1 string `json:"street_1" bson:"street_1"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	ZipCode string `json:"zip_code" bson:"zip_code"`
}

type GeoLocation struct {
	Type        ShapeType `json:"type"` // One of Point, Line
	Coordinates []float64 `json:"coordinates"`
}

type Location struct {
	ID          string      `json:"id" bson:"_id"`
	Name        string      `json:"name,omitempty" bson:"name,omitempty"`
	User        string      `json:"-" bson:"user_id"`
	CreatedAt   uint        `json:"created_at" bson:"created_at"`
	Address     Address     `json:"address,omitempty" bson:"address,omitempty"`
	Geolocation GeoLocation `json:"geo" bson:"geo"`
}

type LocationRequest struct {
	Limit int
	Token string
	Ids   []string
	Name  string
	Lat   float64
	Long  float64
	User  *string
}

type LocationService interface {
	Create(context.Context, *LocationRequest) (*Location, error)
	Update(context.Context, *LocationRequest) (*Location, error)
	FindByID(context.Context, string) (*Location, error)
	FindAll(context.Context, *LocationRequest) ([]*Location, string, error)
}

type LastLocations interface {
	Locations(context.Context, int) ([]*Location, error)
}

type UpdatePlace struct {
	Name     string         `json:"name"`
	Location *LocationInput `json:"location,omitempty"`
}

type LocationInput struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type ShapeType string

const (
	ShapeTypePoint     ShapeType = "POINT"
	ShapeTypeLine      ShapeType = "LINE"
	ShapeTypePoligon   ShapeType = "POLIGON"
	ShapeTypeMultiline ShapeType = "MULTILINE"
)

var AllShapeType = []ShapeType{
	ShapeTypePoint,
	ShapeTypeLine,
	ShapeTypePoligon,
	ShapeTypeMultiline,
}

func (e ShapeType) IsValid() bool {
	switch e {
	case ShapeTypePoint, ShapeTypeLine, ShapeTypePoligon, ShapeTypeMultiline:
		return true
	}
	return false
}

func (e ShapeType) String() string {
	return string(e)
}

func (e *ShapeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ShapeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ShapeType", str)
	}
	return nil
}

func (e ShapeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
