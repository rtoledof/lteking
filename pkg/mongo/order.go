package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.OrderService = &OrderService{}

var OrderCollection Collections = "orders"

type OrderService struct {
	db         *DB
	collection *mongo.Collection
	orderChan  chan *cubawheeler.Order
}

func NewOrderService(db *DB) *OrderService {
	return &OrderService{
		db:         db,
		orderChan:  make(chan *cubawheeler.Order, 10000),
		collection: db.client.Database(database).Collection(OrderCollection.String()),
	}
}

func (s *OrderService) Create(ctx context.Context, input []cubawheeler.OrderItem) (*cubawheeler.Order, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	priceXsec := 1
	priceXm := 100
	order := cubawheeler.Order{
		ID:        cubawheeler.NewID().String(),
		Items:     input,
		Rider:     usr.ID,
		Status:    cubawheeler.OrderStatusNew,
		CreatedAt: time.Now().UTC().Unix(),
	}

	var price uint64
	for _, v := range order.Items {
		price += v.Seconds*uint64(priceXsec) + uint64(priceXm)*v.Meters
	}
	order.Price = price

	_, err := s.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("unable to store the trip: %w", err)
	}

	s.orderChan <- &order

	return &order, nil
}

func (s *OrderService) Update(ctx context.Context, trip *cubawheeler.UpdateOrder) (*cubawheeler.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OrderService) FindByID(ctx context.Context, id string) (*cubawheeler.Order, error) {
	limit := 1
	trips, _, err := findOrders(ctx, s.collection, &cubawheeler.OrderFilter{
		Ids:   []*string{&id},
		Limit: &limit,
	})
	if err != nil && len(trips) == 0 {
		return nil, errors.New("trip not found")
	}
	return trips[0], nil
}

func (s *OrderService) FindAll(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	trips, token, err := findOrders(ctx, s.collection, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.OrderList{Data: trips, Token: token}, nil
}

func findOrders(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.OrderFilter) ([]*cubawheeler.Order, string, error) {
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, "", errors.New("invalid token provided")
	}
	switch user.Role {
	case cubawheeler.RoleRider:
		filter.Rider = &user.ID
	case cubawheeler.RoleDriver:
		filter.Driver = &user.ID
	}
	var trips []*cubawheeler.Order
	var token string
	f := bson.D{}
	if filter.Rider != nil {
		f = append(f, bson.E{Key: "rider", Value: filter.Rider})
	}
	if filter.Driver != nil {
		f = append(f, bson.E{Key: "driver", Value: filter.Driver})
	}
	if filter.Token != nil {
		f = append(f, bson.E{Key: "_id", Value: primitive.E{"$gt", filter.Token}})
	}

	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var trip cubawheeler.Order
		err := cur.Decode(&trip)
		if err != nil {
			return nil, "", err
		}
		trips = append(trips, &trip)
		if len(trips) == *filter.Limit+1 {
			token = trips[*filter.Limit].ID
			trips = trips[:*filter.Limit]
			break
		}
	}

	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return trips, token, nil
}
