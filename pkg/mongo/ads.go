package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

var _ cubawheeler.AdsService = &AdsService{}

const AdsCollection Collections = "ads"

type AdsService struct {
	db *DB
}

func NewAdsService(db *DB) *AdsService {
	return &AdsService{
		db: db,
	}
}

func (s *AdsService) Create(ctx context.Context, request *cubawheeler.AdsRequest) (_ *cubawheeler.Ads, err error) {
	defer derrors.Wrap(&err, "mongo.AdsService.Create")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleClient && usr.Role != cubawheeler.RoleAdmin {
		return nil, errors.New("access denied")
	}
	var ads cubawheeler.Ads
	assambleAds(&ads, request)
	if ads.ID == "" {
		ads.ID = cubawheeler.NewID().String()
	}
	if ads.Owner == "" && usr.Role == cubawheeler.RoleClient {
		ads.Owner = usr.ID
	}
	if !ads.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s: %w", ads.Status, cubawheeler.ErrInvalidInput)
	}

	return &ads, insertAds(ctx, s.db, &ads)
}

func (s *AdsService) Update(ctx context.Context, request *cubawheeler.AdsRequest) (_ *cubawheeler.Ads, err error) {
	defer derrors.Wrap(&err, "mongo.AdsService.Update")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if usr.Role != cubawheeler.RoleClient && usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	ads, err := s.FindById(ctx, request.ID)
	if err != nil {
		return nil, err
	}
	if ads.Owner != usr.ID && usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	assambleAds(ads, request)
	return ads, updateAds(ctx, s.db, ads)
}

func (s *AdsService) FindById(ctx context.Context, id string) (_ *cubawheeler.Ads, err error) {
	defer derrors.Wrap(&err, "mongo.AdsService.FindById")
	return findAdsById(ctx, s.db, id)
}

func (s *AdsService) FindAll(ctx context.Context, request *cubawheeler.AdsRequest) (_ []*cubawheeler.Ads, _ string, err error) {
	defer derrors.Wrap(&err, "mongo.AdsService.FindAll")
	return findAllAds(ctx, s.db, request)
}

func findAllAds(ctx context.Context, db *DB, filter *cubawheeler.AdsRequest) ([]*cubawheeler.Ads, string, error) {
	var adses []*cubawheeler.Ads
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: bson.D{{Key: "$in", Value: filter.Ids}}})
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
	if len(filter.Status) > 0 {
		if !filter.Status.IsValid() {
			return nil, "", fmt.Errorf("invalid status: %s: %w", filter.Status, cubawheeler.ErrInvalidInput)
		}
		f = append(f, bson.E{Key: "status", Value: filter.Status})
	}
	cur, err := db.Collection(AdsCollection).Find(ctx, f)
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
		if len(adses) == filter.Limit+1 && filter.Limit > 0 {
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

func findAdsById(ctx context.Context, db *DB, id string) (*cubawheeler.Ads, error) {
	adses, _, err := findAllAds(ctx, db, &cubawheeler.AdsRequest{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(adses) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return adses[0], nil
}

func insertAds(ctx context.Context, db *DB, ads *cubawheeler.Ads) error {
	_, err := db.Collection(AdsCollection).InsertOne(ctx, ads)
	if err != nil {
		return fmt.Errorf("insert ads: %v: %w", err, cubawheeler.ErrInternal)
	}
	return err
}

func updateAds(ctx context.Context, db *DB, ads *cubawheeler.Ads) error {
	_, err := db.Collection(AdsCollection).UpdateOne(ctx, bson.M{"_id": ads.ID}, bson.M{"$set": ads})
	if err != nil {
		return fmt.Errorf("update ads: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}

func assambleAds(ads *cubawheeler.Ads, request *cubawheeler.AdsRequest) {
	ads.ID = request.ID
	ads.Name = request.Name
	ads.Description = request.Description
	ads.Photo = request.Photo
	ads.Owner = request.Owner
	ads.Status = request.Status
	ads.Priority = request.Priority
	ads.ValidFrom = request.ValidFrom
	ads.ValidUntil = request.ValidUntil
}
