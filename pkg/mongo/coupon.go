package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
	"cubawheeler.io/pkg/derrors"
)

var _ cubawheeler.CouponService = &CouponService{}

const CouponCollection Collections = "coupons"

type CouponService struct {
	db *DB
}

func NewCouponService(db *DB) *CouponService {
	index := mongo.IndexModel{
		Keys: bson.D{{Key: "code", Value: 1}},
	}
	_, err := db.Collection(CouponCollection).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic("unable to create coupon index")
	}
	return &CouponService{
		db: db,
	}
}

func (s *CouponService) Create(ctx context.Context, request *cubawheeler.CouponRequest) (_ *cubawheeler.Coupon, err error) {
	defer derrors.Wrap(&err, "mongo.CouponService.Create")
	coupon := &cubawheeler.Coupon{
		ID:         request.ID,
		Code:       request.Code,
		Percent:    request.Percent,
		Status:     request.Status,
		ValidFrom:  *request.ValidFrom,
		ValidUntil: *request.ValidUntil,
		CreatedAt:  time.Now().UTC().Unix(),
	}
	if coupon.ID == "" {
		coupon.ID = cubawheeler.NewID().String()
	}
	if !coupon.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %w", cubawheeler.ErrInvalidInput)
	}
	coupon.Amount = currency.Amount{
		Amount: request.Amount,
	}
	coupon.Amount.Currency, err = currency.Parse(request.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", cubawheeler.ErrInvalidInput)
	}

	return coupon, insertCoupon(ctx, s.db, coupon)
}

// FindByCode implements cubawheeler.CouponService.
func (s *CouponService) FindByCode(ctx context.Context, code string) (*cubawheeler.Coupon, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	return findByCode(ctx, s.db, code)
}

// Redeem implements cubawheeler.CouponService.
func (s *CouponService) Redeem(ctx context.Context, code string) (*cubawheeler.Coupon, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}

	coupon, err := findByCode(ctx, s.db, code)
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
	w.Deposit(coupon.Amount.Amount)
	err = updateWallet(ctx, s.db, w)
	if err != nil {
		tx.AbortTransaction(ctx)
		return nil, err
	}
	return coupon, tx.CommitTransaction(ctx)
}

func (s *CouponService) FindByID(ctx context.Context, id string) (*cubawheeler.Coupon, error) {
	return findById(ctx, s.db, id)
}

func (s *CouponService) FindAll(ctx context.Context, request *cubawheeler.CouponRequest) ([]*cubawheeler.Coupon, string, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, "", cubawheeler.ErrAccessDenied
	}
	return findCoupons(ctx, s.db, request)
}

func findCoupons(ctx context.Context, db *DB, filter *cubawheeler.CouponRequest) ([]*cubawheeler.Coupon, string, error) {
	var coupons []*cubawheeler.Coupon
	collection := db.Collection(CouponCollection)
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.M{"$in": filter.Ids}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.E{Key: "$gt", Value: filter.Token}})
	}
	if len(filter.Code) > 0 {
		f = append(f, bson.E{Key: "code", Value: filter.Code})
	}
	if filter.Status.IsValid() {
		f = append(f, bson.E{Key: "status", Value: filter.Status})
	}
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var coupon cubawheeler.Coupon
		err := cur.Decode(&coupon)
		if err != nil {
			return nil, "", err
		}
		coupons = append(coupons, &coupon)
		if len(coupons) == filter.Limit+1 && filter.Limit > 0 {
			token = coupons[filter.Limit].ID
			coupons = coupons[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
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

func findByCode(ctx context.Context, db *DB, code string) (*cubawheeler.Coupon, error) {
	coupons, _, err := findCoupons(ctx, db, &cubawheeler.CouponRequest{Code: code})
	if err != nil {
		return nil, fmt.Errorf("error finding coupon: %v: %w", err, cubawheeler.ErrNotFound)
	}
	if len(coupons) == 0 {
		return nil, fmt.Errorf("error finding coupon: %w", cubawheeler.ErrNotFound)
	}
	return coupons[0], nil
}

func findById(ctx context.Context, db *DB, id string) (*cubawheeler.Coupon, error) {
	coupons, _, err := findCoupons(ctx, db, &cubawheeler.CouponRequest{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, fmt.Errorf("error finding coupon: %v: %w", err, cubawheeler.ErrNotFound)
	}
	if len(coupons) == 0 {
		return nil, fmt.Errorf("error finding coupon: %w", cubawheeler.ErrNotFound)
	}
	return coupons[0], nil
}

func insertCoupon(ctx context.Context, db *DB, coupon *cubawheeler.Coupon) error {
	collection := db.Collection(CouponCollection)
	_, err := collection.InsertOne(ctx, coupon)
	if err != nil {
		return err
	}
	return nil
}
