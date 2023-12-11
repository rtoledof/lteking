package cubawheeler

import "context"

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

type UserLocationRequest struct {
	User     string   `json:"user" bson:"user"`
	Location Location `json:"location" bson:"location"`
}

type RealTimeService interface {
	FindNearByDrivers(context.Context, GeoLocation) ([]*User, error)
	NotifyDrivers(context.Context, []*User) error
	UpdateDriversLocation(context.Context, []UserLocationRequest) error
}
