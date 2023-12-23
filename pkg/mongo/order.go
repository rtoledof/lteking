package mongo

import (
	"context"
	"encoding/base64"
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

func (s *OrderService) Create(ctx context.Context, req *cubawheeler.DirectionRequest) (_ *cubawheeler.Order, err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.Create")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleRider {
		return nil, fmt.Errorf("invalid user to create the order: %w", cubawheeler.ErrAccessDenied)
	}

	order, err := s.prepareOrder(ctx, nil, req)
	if err != nil {
		return nil, err
	}

	_, err = s.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("unable to store the trip: %w", err)
	}

	return order, nil
}

func (s *OrderService) CalculatePrice(o *cubawheeler.Order) error {
	// Get brands y rate to be aplied
	brands, _, err := findVehicleCategoriesRate(context.Background(), s.db, cubawheeler.VehicleCategoryRateFilter{})
	if err != nil {
		return err
	}
	// Calculate the price of the trip and store it in the order
	rates, _, err := findRates(context.Background(), s.db, &cubawheeler.RateFilter{})
	if err != nil {
		return err
	}
	var rate *cubawheeler.Rate = rates[0]
	for _, r := range rates {
		if checkRate(r) {
			rate = r
			break
		}
	}
	if rate == nil {
		return fmt.Errorf("no rate found: %w", cubawheeler.ErrNotFound)
	}
	price := uint64(rate.BasePrice) + uint64(float64(rate.PricePerKm)*(o.Distance/1000)) + uint64(float64(rate.PricePerMin)*o.Duration/60) + uint64(float64(rate.PricePerPassenger)*float64(o.Items.Riders))
	if o.Items.Baggages {
		price += uint64(rate.PricePerBaggage)
	}
	if rate.PricePerCarryPet != 0 {
		price += uint64(rate.PricePerCarryPet)
	}

	for _, b := range brands {
		o.CategoryPrice = append(o.CategoryPrice, &cubawheeler.CategoryPrice{
			Category: b.Category,
			Price:    uint64(float64(price) * b.Factor),
		})
	}
	return nil
}

// Request with brand and price
// Accept the order for the client
// Send the order to the near by drivers with are riding a vehicle of the same brand

func (s *OrderService) Update(ctx context.Context, req *cubawheeler.DirectionRequest) (_ *cubawheeler.Order, err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.Update")
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != cubawheeler.RoleRider {
		return nil, fmt.Errorf("invalid user to update the order: %w", cubawheeler.ErrAccessDenied)
	}
	order, err := findOrderById(ctx, s.db, req.ID)
	if err != nil {
		return nil, err
	}
	order, err = s.prepareOrder(ctx, order, req)
	if err != nil {
		return nil, err
	}

	if err := updateOrder(ctx, s.db, order.ID, order); err != nil {
		return nil, err
	}

	return order, nil

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
	order.Driver = usr.ID
	order.Status = cubawheeler.OrderStatusConfirmed
	if err := updateOrder(ctx, s.db, order.ID, order); err != nil {
		return nil, err
	}
	// TODO: send the order to the drivers
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
	if err = updateOrder(ctx, s.db, order.ID, order); err != nil {
		return nil, err
	}
	// TODO: if the order was sent to the drivers, cancel it
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
	order.EndAt = time.Now().UTC().Unix()
	if err = updateOrder(ctx, s.db, order.ID, order); err != nil {
		return nil, err
	}
	lastPoint := order.Items.Points[len(order.Items.Points)-1]
	user.LastLocations = append(user.LastLocations, &cubawheeler.Location{
		Name: "",
		Geolocation: cubawheeler.GeoLocation{
			Lat:  lastPoint.Lat,
			Long: lastPoint.Lng,
		},
	})
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
	order.StartAt = time.Now().UTC().Unix()
	if err = updateOrder(ctx, s.db, order.ID, order); err != nil {
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

func updateOrder(ctx context.Context, db *DB, id string, order *cubawheeler.Order) error {
	collection := db.client.Database(database).Collection(OrderCollection.String())
	if _, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: order}}); err != nil {
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

func assambleOrder(o *cubawheeler.Order, req *cubawheeler.DirectionRequest) *cubawheeler.Order {
	if req.Riders == 0 {
		req.Riders = 1
	}
	o.Items = *req
	o.Distance = 0
	o.Duration = 0
	return o
}

func (s *OrderService) prepareOrder(ctx context.Context, order *cubawheeler.Order, req *cubawheeler.DirectionRequest) (*cubawheeler.Order, error) {
	if order == nil {
		order = &cubawheeler.Order{
			ID:        cubawheeler.NewID().String(),
			Status:    cubawheeler.OrderStatusNew,
			CreatedAt: time.Now().UTC().Unix(),
		}
	}
	mapbox := mapbox.NewClient(os.Getenv("MAPBOX_TOKEN"))
	assambleOrder(order, req)

	order.Items = *req

	routes, strBody, err := mapbox.Directions.GetRoute(ctx, *req)
	if err != nil {
		return nil, fmt.Errorf("unable to get the route: %w", err)
	}

	if len(routes.Routes) == 0 {
		return nil, fmt.Errorf("no routes found: %w", cubawheeler.ErrNotFound)
	}

	route := routes.Routes[0]
	order.Distance = route.Distance
	order.Duration = route.Duration
	order.RouteString = base64.StdEncoding.EncodeToString([]byte(strBody))

	err = s.CalculatePrice(order)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate the price: %w", err)
	}
	return order, nil
}

func checkRate(r *cubawheeler.Rate) bool {
	currentTime := time.Now().UTC()
	if r.StartDate != 0 {
		if (currentTime.Unix() < r.StartDate || currentTime.Unix() > r.EndDate) ||
			(r.EndDate != 0 && currentTime.Unix() > r.EndDate || r.StartDate < currentTime.Unix()) {
			return false
		}
	}
	today := currentTime.Format("2006-01-02")
	if r.StartTime != "" {
		rateStartTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", today, r.StartTime))
		if err != nil {
			return false
		}
		if currentTime.Before(rateStartTime) {
			return false
		}
	}
	if r.EndTime != "" {
		rateEndTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", today, r.EndTime))
		if err != nil {
			return false
		}
		if currentTime.After(rateEndTime) {
			return false
		}
	}
	return true
}
