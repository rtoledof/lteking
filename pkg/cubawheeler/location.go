package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Point struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lon" bson:"lon"`
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
	Bearing     float64   `json:"bearing"`
	Lat         float64   `json:"lat"`
	Long        float64   `json:"long"`
}

//{
//	"_id" : ObjectId("59a47286cfa9a3a73e51e75c"),
//	"theaterId" : 104,
//	"location" : {
//		"address" : {
//			"street1" : "5000 W 147th St",
//			"city" : "Hawthorne",
//			"state" : "CA",
//			"zipcode" : "90250"
//		},
//		"geo" : {
//			"type" : "Point",
//			"coordinates" : [
//				-118.36559,
//				33.897167
//			]
//		}
//	}
//}

type Location struct {
	ID          string      `json:"id" bson:"_id"`
	Name        string      `json:"name,omitempty" bson:"name,omitempty"`
	User        string      `json:"-" bson:"user_id"`
	CreatedAt   int64       `json:"created_at" bson:"created_at"`
	UpdatedAt   int64       `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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

type Intersection struct {
	Out      int       `json:"out" bson:"out"`
	In       int       `json:"in" bson:"in"`
	Entry    []bool    `json:"entry" bson:"entry"`
	Bearings []int     `json:"bearings" bson:"bearings"`
	Location []float64 `json:"location" bson:"location"`
}

type Maneuver struct {
	Location      []float64 `json:"location" bson:"location"`
	BearingAfter  int       `json:"bearing_after" bson:"bearing_after"`
	BearingBefore int       `json:"bearing_before" bson:"bearing_before"`
	Type          string    `json:"type" bson:"type"`
	Modifier      string    `json:"modifier" bson:"modifier"`
	Instruction   string    `json:"instruction" bson:"instruction"`
}

type VoiceInstructions struct {
	Announcement          string  `json:"announcement" bson:"announcement"`
	DistanceAlongGeometry float64 `json:"distanceAlongGeometry" bson:"distanceAlongGeometry"`
	SsmlAnnouncement      string  `json:"ssmlAnnouncement" bson:"ssmlAnnouncement"`
}

type Component struct {
	Text string `json:"text" bson:"text"`
}

type Primary struct {
	Text       string       `json:"text" bson:"text"`
	Type       string       `json:"type" bson:"type"`
	Modifier   string       `json:"modifier" bson:"modifier"`
	Components []*Component `json:"components" bson:"components"`
}

type BannerInstructions struct {
	DistanceAlongGeometry float64  `json:"distanceAlongGeometry" bson:"distanceAlongGeometry"`
	Primary               *Primary `json:"primary" bson:"primary"`
	Secondary             *Primary `json:"secondary" bson:"secondary"`
}

type Step struct {
	Intersections      []*Intersection       `json:"intersections" bson:"intersections"`
	DrivingSide        string                `json:"driving_side" bson:"driving_side"`
	Geometry           string                `json:"geometry" bson:"geometry"`
	Mode               string                `json:"mode" bson:"mode"`
	Maneuver           *Maneuver             `json:"maneuver" bson:"maneuver"`
	Weight             float64               `json:"weight" bson:"weight"`
	Duration           float64               `json:"duration" bson:"duration"`
	Name               string                `json:"name" bson:"name"`
	Distance           float64               `json:"distance" bson:"distance"`
	Ref                string                `json:"ref" bson:"ref"`
	VoiceInstructions  []*VoiceInstructions  `json:"voiceInstructions" bson:"voiceInstructions"`
	BannerInstructions []*BannerInstructions `json:"bannerInstructions" bson:"bannerInstructions"`
}

type Legs struct {
	Steps    []Step  `json:"steps" bson:"steps"`
	Weight   float64 `json:"weight" bson:"weight"`
	Distance float64 `json:"distance" bson:"distance"`
	Summary  string  `json:"summary" bson:"summary"`
	Duration float64 `json:"duration" bson:"duration"`
}

type Geometry struct {
	Coordinates []Point `json:"coordinates" bson:"coordinates"`
}

type Waitpoint struct {
	Location []float64 `json:"location" bson:"location"`
	Name     string    `json:"name" bson:"name"`
}

type Route struct {
	Legs       []*Legs      `json:"legs" bson:"legs"`
	WeightName string       `json:"weight_name" bson:"weight_name"`
	Geometry   string       `json:"geometry" bson:"geometry"`
	Weight     float64      `json:"weight" bson:"weight"`
	Distance   float64      `json:"distance" bson:"distance"`
	Duration   float64      `json:"duration" bson:"duration"`
	Waitpoints []*Waitpoint `json:"waypoints" bson:"waypoints"`
}

type DirectionResponse struct {
	Geometry  string       `json:"geometry" bson:"geometry"`
	Duration  float64      `json:"duration" bson:"duration"`
	Distance  float64      `json:"distance" bson:"distance"`
	Waypoints []*Waitpoint `json:"waypoints" bson:"waypoints"`
	Routes    []*Route     `json:"routes" bson:"routes"`
}

type DirectionRequest struct {
	Points []*Point `json:"points" bson:"points"`
}

func (r *DirectionRequest) Valid() bool {
	return len(r.Points) <= 1
}

func (r *DirectionRequest) String() string {
	var response []string

	for _, point := range r.Points {
		response = append(response, fmt.Sprintf("%f,%f", point.Lat, point.Lng))
	}
	return strings.Join(response, ";")
}
