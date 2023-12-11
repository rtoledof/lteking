package mongo

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//db.restaurants.find({ location:
//	{ $geoWithin:
//		{ $centerSphere: [ [ -73.93414657, 40.82302903 ], 5 / 3963.2 ] } } })

const DriverRealTimeLocationCollection Collections = "drivers_location"

var _ cubawheeler.RealTimeService = &RealTimeService{}

type RealTimeService struct {
	db *DB
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

func NewRealTimeService(db *DB) *RealTimeService {
	index := mongo.IndexModel{
		Keys: bson.D{{"location.geo", "2dsphere"}},
	}
	_, err := db.client.Database(database).Collection(DriverRealTimeLocationCollection.String()).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic("unable to create geo index")
	}

	return &RealTimeService{db: db}
}

func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location cubawheeler.GeoLocation) ([]*cubawheeler.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RealTimeService) NotifyDrivers(ctx context.Context, users []*cubawheeler.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *RealTimeService) UpdateDriversLocation(context.Context, []cubawheeler.UserLocationRequest) error {
	panic("implement me")
}
