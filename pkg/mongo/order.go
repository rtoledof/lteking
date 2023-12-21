package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
	"cubawheeler.io/pkg/mapbox"
)

var _ cubawheeler.OrderService = &OrderService{}

var OrderCollection Collections = "orders"

type OrderService struct {
	db         *DB
	collection *mongo.Collection
	orderChan  chan *cubawheeler.Order
	mutex      sync.Map
}

func NewOrderService(db *DB) *OrderService {
	return &OrderService{
		db:         db,
		orderChan:  make(chan *cubawheeler.Order, 10000),
		collection: db.client.Database(database).Collection(OrderCollection.String()),
	}
}

func (s *OrderService) Create(ctx context.Context, input []*cubawheeler.OrderItem) (_ *cubawheeler.Order, err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.Create")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}

	// Must connect to to mapbox and get the route.
	mapbox := mapbox.NewClient(os.Getenv("MAPBOX_TOKEN"))
	// Get brands y rate to be aplied
	// Calculate the price of the trip and store it in the order
	// Send the order, brand and price to the client
	order := cubawheeler.Order{
		ID:        cubawheeler.NewID().String(),
		Items:     input,
		Rider:     usr.ID,
		Status:    cubawheeler.OrderStatusNew,
		CreatedAt: time.Now().UTC().Unix(),
	}

	var req cubawheeler.DirectionRequest
	for _, l := range input {
		req.AddPoint(&l.PickUp)
		req.AddPoint(&l.DropOff)
	}

	route, err := mapbox.Directions.GetRoute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to get the route: %w", err)
	}
	order.Distance = route.Distance
	order.Duration = route.Duration
	order.Route = route

	order.Price, err = s.CalculatePrice(&order)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate the price: %w", err)
	}

	_, err = s.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("unable to store the trip: %w", err)
	}

	// realtime.OrderChan <- &order

	return &order, nil
}

func (s *OrderService) CalculatePrice(o *cubawheeler.Order) (uint64, error) {
	// Get brands y rate to be aplied
	// Calculate the price of the trip and store it in the order
	rates, _, err := findRates(context.Background(), s.db, &cubawheeler.RateFilter{})
	if err != nil {
		return 0, err
	}
	var rate *cubawheeler.Rate
	for _, r := range rates {
		if checkRate(r) {
			rate = r
			break
		}
	}
	return uint64(rate.BasePrice) + uint64(float64(rate.PricePerKm)*o.Distance) + uint64(float64(rate.PricePerMin)*o.Duration), nil
}

func checkRate(r *cubawheeler.Rate) bool {
	currentTime := time.Now().UTC()
	if (r.StartDate != 0 && currentTime.Unix() < r.StartDate || currentTime.Unix() > r.EndDate) ||
		(r.EndDate != 0 && currentTime.Unix() > r.EndDate || r.StartDate < currentTime.Unix()) {
		return false
	}
	if r.StartTime != "" {
		rateStartTime, err := time.Parse("15:04", r.StartTime)
		if err != nil {
			return false
		}
		if currentTime.Before(rateStartTime) {
			return false
		}
	}
	if r.EndTime != "" {
		rateEndTime, err := time.Parse("15:04", r.EndTime)
		if err != nil {
			return false
		}
		if currentTime.After(rateEndTime) {
			return false
		}
	}
	return true
}

// Request with brand and price
// Accept the order for the client
// Send the order to the near by drivers with are riding a vehicle of the same brand

func (s *OrderService) Update(ctx context.Context, trip *cubawheeler.UpdateOrder) (*cubawheeler.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OrderService) FindByID(ctx context.Context, id string) (*cubawheeler.Order, error) {
	return findOrderById(ctx, s.db, id)
}

func (s *OrderService) FindAll(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	trips, token, err := findOrders(ctx, s.collection, filter)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.OrderList{Data: trips, Token: token}, nil
}

func (s *OrderService) AcceptOrder(ctx context.Context, id string) (*cubawheeler.Order, error) {
	s.orderLock(id)
	defer s.orderUnlock(id)
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", cubawheeler.ErrAccessDenied)
	}
	if usr.Role != cubawheeler.RoleDriver {
		return nil, fmt.Errorf("invalid user to acept the order: %w", cubawheeler.ErrAccessDenied)
	}
	order, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if order.Driver != "" {
		return nil, cubawheeler.ErrOrderAccepted
	}
	f := bson.D{}
	order.Driver = usr.ID
	f = append(f, bson.E{Key: "driver", Value: usr.ID})
	order.Status = cubawheeler.OrderStatusOnTheWay
	f = append(f, bson.E{Key: "status", Value: order.Status})
	if err := updateOrder(ctx, s.db, order.ID, f); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) (*cubawheeler.Order, error) {
	s.orderLock(id)
	defer s.orderUnlock(id)
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if user.Role != cubawheeler.RoleDriver && user.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	order, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	order.Status = cubawheeler.OrderStatusDropOff
	f := bson.D{{Key: "status", Value: order.Status}}
	if err = updateOrder(ctx, s.db, order.ID, f); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CompleteOrder(ctx context.Context, id string) (*cubawheeler.Order, error) {
	s.orderLock(id)
	defer s.orderUnlock(id)
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if user.Role != cubawheeler.RoleDriver && user.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	order, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	order.Status = cubawheeler.OrderStatusDropOff
	f := bson.D{{Key: "status", Value: order.Status}}
	order.EndAt = time.Now().UTC().Unix()
	f = append(f, bson.E{Key: "end_at", Value: order.EndAt})
	if err = updateOrder(ctx, s.db, order.ID, f); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) StartOrder(ctx context.Context, id string) (*cubawheeler.Order, error) {
	s.orderLock(id)
	defer s.orderUnlock(id)
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if user.Role != cubawheeler.RoleDriver && user.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	order, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if order.Driver != user.ID {
		return nil, cubawheeler.ErrAccessDenied
	}
	order.Status = cubawheeler.OrderStatusPickUp
	f := bson.D{{Key: "status", Value: order.Status}}
	order.StartAt = time.Now().UTC().Unix()
	f = append(f, bson.E{Key: "start_at", Value: order.EndAt})
	if err = updateOrder(ctx, s.db, order.ID, f); err != nil {
		return nil, err
	}
	return order, nil
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
		f = append(f, bson.E{Key: "_id", Value: primitive.E{Key: "$gt", Value: filter.Token}})
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
		if len(trips) == filter.Limit+1 {
			token = trips[filter.Limit].ID
			trips = trips[:filter.Limit]
			break
		}
	}

	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	cur.Close(ctx)
	return trips, token, nil
}

func findOrderById(ctx context.Context, db *DB, id string) (*cubawheeler.Order, error) {
	collection := db.client.Database(database).Collection(OrderCollection.String())
	trips, _, err := findOrders(ctx, collection, &cubawheeler.OrderFilter{
		Ids:   []*string{&id},
		Limit: 1,
	})
	if err != nil && len(trips) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return trips[0], nil
}

func updateOrder(ctx context.Context, db *DB, id string, f bson.D) error {
	collection := db.client.Database(database).Collection(OrderCollection.String())
	if _, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: f}}); err != nil {
		return fmt.Errorf("unabe to update the order: %v: %w", err, cubawheeler.ErrInternal)
	}
	return nil
}

func (s *OrderService) orderLock(key string) {
	mu := sync.Mutex{}
	l_, _ := s.mutex.LoadOrStore(key, &mu)
	l := l_.(*sync.Mutex)
	l.Lock()
	if l != &mu {
		l.Unlock()
		s.orderLock(key)
	}
}

func (s *OrderService) orderUnlock(key string) {
	l_, ok := s.mutex.Load(key)
	if !ok {
		return
	}
	l := l_.(*sync.Mutex)
	s.mutex.Delete(key)
	l.Unlock()
}
