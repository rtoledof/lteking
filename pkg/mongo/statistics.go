package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

var _ cubawheeler.OrderStatisticsService = (*OrderStatistics)(nil)

var StatisticsCollection Collections = "statistics"

type OrderStatistics struct {
	db *DB
}

func NewOrderStatistics(db *DB) *OrderStatistics {
	return &OrderStatistics{
		db: db,
	}
}

// AddOrder implements cubawheeler.OrderStatisticsService.
func (s *OrderStatistics) AddOrder(ctx context.Context, statistics cubawheeler.OrderStatistics) (err error) {
	defer derrors.Wrap(&err, "unable to add order")
	if statistics.ID == "" {
		statistics.ID = cubawheeler.NewID().String()
	}
	return insertStatistics(ctx, s.db, statistics)
}

// FindAllStatistics implements cubawheeler.OrderStatisticsService.
func (s *OrderStatistics) FindAllStatistics(ctx context.Context, filter cubawheeler.OrderStatisticsFilter) (_ []*cubawheeler.OrderStatistics, _ string, err error) {
	defer derrors.Wrap(&err, "unable to find statistics")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, "", cubawheeler.ErrAccessDenied
	}
	if usr.Role != cubawheeler.RoleAdmin {
		filter.User = usr.ID
	}
	return findStatistics(ctx, s.db, filter)
}

// FindStatistictsByUser implements cubawheeler.OrderStatisticsService.
func (s *OrderStatistics) FindStatistictsByUser(ctx context.Context, user string) (_ *cubawheeler.OrderStatistics, err error) {
	defer derrors.Wrap(&err, "unable to find statistics")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if usr.Role != cubawheeler.RoleAdmin {
		user = usr.ID
	}
	stats, _, err := findStatistics(ctx, s.db, cubawheeler.OrderStatisticsFilter{
		User: user,
	})
	if err != nil {
		return nil, err
	}
	if len(stats) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return stats[0], nil
}

func findStatistics(ctx context.Context, db *DB, filter cubawheeler.OrderStatisticsFilter) ([]*cubawheeler.OrderStatistics, string, error) {
	collection := db.Collection(StatisticsCollection)
	f := bson.D{}
	if len(filter.User) > 0 {
		f = append(f, bson.E{Key: "user", Value: filter.User})
	}
	if len(filter.Token) > 0 {
		f = append(f, primitive.E{Key: "_id", Value: primitive.D{{Key: "$gt", Value: filter.Token}}})
	}
	cursor, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)
	var resp []*cubawheeler.OrderStatistics
	var token string
	for cursor.Next(ctx) {
		var orderStatistics cubawheeler.OrderStatistics
		if err := cursor.Decode(&orderStatistics); err != nil {
			return nil, "", err
		}
		resp = append(resp, &orderStatistics)
		if len(resp) == filter.Limit && filter.Limit > 0 {
			token = resp[len(resp)].ID
			resp = resp[:filter.Limit]

		}
	}
	if err := cursor.Err(); err != nil {
		return nil, "", err
	}
	return resp, token, nil
}

func insertStatistics(ctx context.Context, db *DB, statistics cubawheeler.OrderStatistics) error {
	collection := db.Collection(StatisticsCollection)
	if _, err := collection.InsertOne(ctx, statistics); err != nil {
		return fmt.Errorf("unable to insert statistics: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}
