package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.AdsService = &AdsService{}

const CollectionAds Collections = "ads"

type AdsService struct {
	db         *DB
	collection *mongo.Collection
}

func NewAdsService(db *DB) *AdsService {
	return &AdsService{
		db:         db,
		collection: db.client.Database(database).Collection(CollectionAds.String()),
	}
}

func (s *AdsService) Create(ctx context.Context, request *cubawheeler.AdsRequest) (*cubawheeler.Ads, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleClient && usr.Role != cubawheeler.RoleAdmin {
		return nil, errors.New("access denied")
	}
	var ads = cubawheeler.Ads{
		ID:          cubawheeler.NewID().String(),
		Name:        request.Name,
		Description: request.Description,
		Photo:       request.Photo,
		Owner:       request.Owner,
		Status:      request.Status,
		Priority:    request.Priority,
		ValidFrom:   request.ValidFrom,
		ValidUntil:  request.ValidUntil,
	}
	_, err := s.collection.InsertOne(ctx, &ads)
	if err != nil {
		return nil, errors.New("unable to store the ads")
	}
	return &ads, nil
}

func (s *AdsService) Update(ctx context.Context, request *cubawheeler.AdsRequest) (*cubawheeler.Ads, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AdsService) FindById(ctx context.Context, id string) (*cubawheeler.Ads, error) {
	ads, _, err := findAllAds(ctx, s.collection, &cubawheeler.AdsRequest{Ids: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(ads) == 0 {
		return nil, errors.New("ads not found")
	}
	return ads[0], nil
}

func (s *AdsService) FindAll(ctx context.Context, request *cubawheeler.AdsRequest) ([]*cubawheeler.Ads, string, error) {
	return findAllAds(ctx, s.collection, request)
}

func findAllAds(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.AdsRequest) ([]*cubawheeler.Ads, string, error) {
	var adses []*cubawheeler.Ads
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: primitive.A{"$in", filter.Ids}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{Key: "_id", Value: primitive.E{Key: "$gt", Value: filter.Token}})
	}
	if len(filter.Name) > 0 {
		f = append(f, bson.E{Key: "name", Value: filter.Name})
	}
	if len(filter.Owner) > 0 {
		f = append(f, bson.E{Key: "owner", Value: filter.Owner})
	}
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var ads cubawheeler.Ads
		err := cur.Decode(&ads)
		if err != nil {
			return nil, "", err
		}
		adses = append(adses, &ads)
		if len(adses) == filter.Limit+1 {
			token = adses[filter.Limit].ID
			adses = adses[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return adses, token, nil
}
