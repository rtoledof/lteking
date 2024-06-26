package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"order.io/pkg/order"
)

var key = "drivers"

type RealTimeService struct {
	redis *Redis
}

func NewRealTimeService(redis *Redis) *RealTimeService {
	return &RealTimeService{redis: redis}
}

func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location order.GeoLocation) ([]*order.Location, error) {
	res, _ := s.redis.client.GeoRadius(ctx, key, location.Long, location.Lat, &redis.GeoRadiusQuery{
		Radius:      500, // TODO: change this to 5km
		Unit:        "km",
		WithGeoHash: true,
		WithCoord:   true,
		WithDist:    true,
		Count:       20,
		Sort:        "ASC",
	}).Result()
	var locations []*order.Location
	for _, l := range res {
		var riderLocation = order.Location{
			User: l.Name,
			Geolocation: order.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{l.Longitude, l.Latitude},
				Lat:         l.Latitude,
				Long:        l.Longitude,
			},
		}
		locations = append(locations, &riderLocation)
	}
	return locations, nil
}

func (s *RealTimeService) UpdateLocation(ctx context.Context, user string, location order.GeoLocation) error {
	var geoLocations []*redis.GeoLocation
	geoLocations = append(geoLocations, &redis.GeoLocation{
		Name:      user,
		Longitude: location.Long,
		Latitude:  location.Lat,
	})
	if err := s.redis.client.GeoAdd(ctx, key, geoLocations...).Err(); err != nil {
		return fmt.Errorf("unable to update driver locations: %v: %w", err, order.ErrInternal)
	}
	return nil
}
