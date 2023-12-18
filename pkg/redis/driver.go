package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/errors"
	"cubawheeler.io/pkg/realtime"
)

var _ realtime.FinderUpdater = &RealTimeService{}

var key = "drivers"

type RealTimeService struct {
	redis *Redis
}

func NewRealTimeService(redis *Redis) *RealTimeService {
	return &RealTimeService{redis: redis}
}

func (s *RealTimeService) FindNearByDrivers(ctx context.Context, location cubawheeler.GeoLocation) ([]*cubawheeler.Location, error) {
	//TODO implement me
	res, _ := s.redis.client.GeoRadius(ctx, key, location.Long, location.Lat, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithGeoHash: true,
		WithCoord:   true,
		WithDist:    true,
		Count:       20,
		Sort:        "ASC",
	}).Result()
	var locations []*cubawheeler.Location
	for _, l := range res {
		var riderLocation = cubawheeler.Location{
			User: l.Name,
			Geolocation: cubawheeler.GeoLocation{
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

func (s *RealTimeService) UpdateLocation(ctx context.Context, user string, location cubawheeler.GeoLocation) error {
	var geoLocations []*redis.GeoLocation
	geoLocations = append(geoLocations, &redis.GeoLocation{
		Name:      user,
		Longitude: location.Long,
		Latitude:  location.Lat,
	})
	if err := s.redis.client.GeoAdd(ctx, key, geoLocations...); err != nil {
		return fmt.Errorf("unable to update driver locations: %v: %w", err, errors.ErrInternal)
	}
	return nil
}
