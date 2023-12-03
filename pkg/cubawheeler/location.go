package cubawheeler

import "context"

type Location struct {
	ID   string  `json:"id" bson:"_id"`
	Name string  `json:"name" bson:"name"`
	Lat  float64 `json:"lat" bson:"lat"`
	Long float64 `json:"long" bson:"long"`
	User string  `json:"-" bson:"user_id"`
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

type UpdatePlace struct {
	Name     string         `json:"name"`
	Location *LocationInput `json:"location,omitempty"`
}

type LocationInput struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
