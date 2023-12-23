package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.CouponService = &CouponService{}

type CouponService struct {
	db         *DB
	collection *mongo.Collection
}

func NewCouponService(db *DB) *CouponService {
	return &CouponService{
		db:         db,
		collection: db.client.Database(database).Collection("coupons"),
	}
}

func (s *CouponService) Create(ctx context.Context, request *cubawheeler.CouponRequest) (*cubawheeler.Coupon, error) {
	coupon := &cubawheeler.Coupon{
		ID:         cubawheeler.NewID().String(),
		Code:       request.Code,
		Percent:    request.Percent,
		Amount:     request.Amount,
		Status:     request.Status,
		ValidFrom:  request.ValidFrom,
		ValidUntil: request.ValidUntil,
	}
	_, err := s.collection.InsertOne(ctx, coupon)
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

// FindByCode implements cubawheeler.CouponService.
func (s *CouponService) FindByCode(ctx context.Context, code string) (*cubawheeler.Coupon, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	return findByCode(ctx, s.collection, code)
}

// Redeem implements cubawheeler.CouponService.
func (s *CouponService) Redeem(ctx context.Context, code string) (*cubawheeler.Coupon, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}

	coupon, err := findByCode(ctx, s.collection, code)
	if err != nil {
		return nil, err
	}
	coupon.Status = cubawheeler.CouponStatusRedeemed
	tx, err := s.db.client.StartSession()
	if err != nil {
		return nil, err
	}
	if err := tx.StartTransaction(); err != nil {
		return nil, err
	}
	err = updateCoupon(ctx, s.db, coupon)
	if err != nil {
		tx.AbortTransaction(ctx)
		return nil, err
	}
	w, err := findWalletByOwner(ctx, s.db, user.ID)
	if err != nil {
		tx.AbortTransaction(ctx)
		return nil, err
	}
	w.Deposit(*coupon.Amount)
	err = updateWallet(ctx, s.db, w)
	if err != nil {
		tx.AbortTransaction(ctx)
		return nil, err
	}
	return coupon, tx.CommitTransaction(ctx)
}

func (s *CouponService) FindByID(ctx context.Context, id string) (*cubawheeler.Coupon, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	users, _, err := findAllCoupon(ctx, s.collection, &cubawheeler.CouponRequest{Ids: []string{id}})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (s *CouponService) FindAll(ctx context.Context, request *cubawheeler.CouponRequest) ([]*cubawheeler.Coupon, string, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, "", cubawheeler.ErrAccessDenied
	}
	return findAllCoupon(ctx, s.collection, request)
}

func findAllCoupon(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.CouponRequest) ([]*cubawheeler.Coupon, string, error) {
	var coupons []*cubawheeler.Coupon
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{"_id", primitive.A{"$in", filter.Ids}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{"_id", primitive.E{"$gt", filter.Token}})
	}
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var coupon cubawheeler.Coupon
		err := cur.Decode(&coupon)
		if err != nil {
			return nil, "", err
		}
		coupons = append(coupons, &coupon)
		if len(coupons) == filter.Limit+1 {
			token = coupons[filter.Limit].ID
			coupons = coupons[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return coupons, token, nil
}

func updateCoupon(ctx context.Context, db *DB, coupon *cubawheeler.Coupon) error {
	collection := db.Collection("coupons")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": coupon.ID}, bson.M{"$set": coupon})
	if err != nil {
		return err
	}
	return nil
}

func findByCode(ctx context.Context, collection *mongo.Collection, code string) (*cubawheeler.Coupon, error) {
	var coupon cubawheeler.Coupon
	err := collection.FindOne(ctx, bson.M{"code": code}).Decode(&coupon)
	if err != nil {
		return nil, fmt.Errorf("error finding coupon: %v: %w", err, cubawheeler.ErrNotFound)
	}
	return &coupon, nil
}
