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

	"order.io/pkg/currency"
	"order.io/pkg/derrors"
	"order.io/pkg/mapbox"
	"order.io/pkg/order"
	"order.io/pkg/redis"
)

var _ order.OrderService = &OrderService{}

var OrderCollection Collections = "orders"

type OrderService struct {
	db        *DB
	orderChan chan *order.Order
	mutex     sync.Map
	redis     *redis.Redis
	direction order.DirectionService
}

func NewOrderService(
	db *DB,
	redis *redis.Redis,
) *OrderService {
	client := mapbox.NewClient(os.Getenv("MAPBOX_TOKEN"))

	return &OrderService{
		db:        db,
		orderChan: make(chan *order.Order, 10000),
		redis:     redis,
		direction: client.Directions,
	}
}

func (s *OrderService) Create(ctx context.Context, req order.Item) (_ *order.Order, err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.Create")
	user := order.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("nil user in context: %w", order.ErrAccessDenied)
	}

	o, err := s.prepareOrder(ctx, nil, order.DirectionRequest{
		Points: req.Points,
	})
	if err != nil {
		return nil, err
	}

	err = storeOrder(ctx, s.db, o)
	if err != nil {
		return nil, fmt.Errorf("unable to store the trip: %w", err)
	}

	return o, nil
}

// ConfirmOrder implements order.OrderService.
func (s *OrderService) ConfirmOrder(ctx context.Context, req order.ConfirmOrder) error {
	s.orderLock(req.OrderID)
	defer s.orderUnlock(req.OrderID)
	usr := order.UserFromContext(ctx)
	if usr == nil || usr.Role != order.RoleRider {
		return order.ErrAccessDenied
	}
	ord, err := findOrderById(ctx, s.db, req.OrderID)
	if err != nil {
		return err
	}
	if ord.Rider != usr.ID {
		return order.ErrAccessDenied
	}
	if ord.Status != order.OrderStatusNew {
		return order.ErrNotFound
	}
	ord.Status = order.OrderStatusConfirmed
	for _, c := range ord.CategoryPrice {
		if c.Category == req.Category {
			ord.Price = int(c.Price)
			ord.SelectedCategory = c
			break
		}
	}
	ord.ChargeMethod = req.Method
	if err := updateOrder(ctx, s.db, ord.ID, ord); err != nil {
		return err
	}
	if err := updateOrder(ctx, s.db, ord.ID, ord); err != nil {
		return err
	}

	if err := s.redis.Publish(ctx, "orders", ord); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) CalculatePrice(o *order.Order) error {
	// Get brands y rate to be aplied
	brands, _, err := findVehicleCategoriesRate(context.Background(), s.db, order.VehicleCategoryRateFilter{})
	if err != nil {
		return err
	}
	// Calculate the price of the trip and store it in the order
	rates, _, err := findRates(context.Background(), s.db, &order.RateFilter{})
	if err != nil {
		return err
	}
	var rate *order.Rate = rates[0]
	for _, r := range rates {
		if checkRate(r) {
			rate = r
			break
		}
	}
	if rate == nil {
		return fmt.Errorf("no rate found: %w", order.ErrNotFound)
	}
	price := price(o.Distance, o.Duration, *rate, o.Item.Riders)

	for _, b := range brands {
		o.CategoryPrice = append(o.CategoryPrice, &order.CategoryPrice{
			Category: b.Category,
			Price:    int(float64(price) * b.Factor),
			Currency: o.Currency,
		})
	}
	return nil
}

// Request with brand and price
// Accept the order for the client
// Send the order to the near by drivers with are riding a vehicle of the same brand

func (s *OrderService) Update(ctx context.Context, id string, req order.Item) (_ *order.Order, err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.Update")
	usr := order.UserFromContext(ctx)
	if usr == nil {
		return nil, errors.New("invalid token provided")
	}
	if usr.Role != order.RoleRider {
		return nil, fmt.Errorf("invalid user to update the order: %w", order.ErrAccessDenied)
	}
	o, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	o, err = s.prepareOrder(ctx, o, order.DirectionRequest{
		Points: req.Points,
	})
	if err != nil {
		return nil, err
	}
	if err := updateOrder(ctx, s.db, o.ID, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *OrderService) FindByID(ctx context.Context, id string) (_ *order.Order, err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.FindByID")
	usr := order.UserFromContext(ctx)
	if usr == nil {
		return nil, order.NewError(order.ErrAccessDenied, 401, "user not provided")
	}
	ord, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if ord.Driver != "" && usr.Role == order.RoleDriver && ord.Driver != usr.ID {
		return nil, order.NewError(order.ErrAccessDenied, 401, "user not provided")
	}
	return ord, nil
}

func (s *OrderService) FindAll(ctx context.Context, filter order.OrderFilter) (*order.OrderList, error) {
	trips, token, err := findOrders(ctx, s.db, filter)
	if err != nil {
		return nil, err
	}
	return &order.OrderList{Data: trips, Token: token}, nil
}

func (s *OrderService) AcceptOrder(ctx context.Context, id string) error {
	s.orderLock(id)
	defer s.orderUnlock(id)
	usr := order.UserFromContext(ctx)
	if usr == nil {
		return fmt.Errorf("nil user in context: %w", order.ErrAccessDenied)
	}
	if usr.Role != order.RoleDriver {
		return fmt.Errorf("invalid user to acept the order: %w", order.ErrAccessDenied)
	}
	ord, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return err
	}
	if ord.Driver != "" {
		return order.ErrOrderAccepted
	}

	ord.Driver = usr.ID
	ord.Status = order.OrderStatusOnTheWay
	if err := updateOrder(ctx, s.db, ord.ID, ord); err != nil {
		return err
	}

	if err := s.redis.Publish(ctx, "order:confirmed", ord); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	s.orderLock(id)
	defer s.orderUnlock(id)
	user := order.UserFromContext(ctx)
	if user == nil {
		return order.ErrAccessDenied
	}
	if user.Role != order.RoleDriver && user.Role != order.RoleAdmin {
		return order.ErrAccessDenied
	}
	ord, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return err
	}
	ord.Status = order.OrderStatusCancel
	if user.Role == order.RoleDriver {
		ord.Status = order.OrderStatusWaitingDriver
		ord.Driver = ""
		ord.BannedDrivers[user.ID] = true
	}

	if err = updateOrder(ctx, s.db, ord.ID, ord); err != nil {
		return err
	}
	if ord.Status == order.OrderStatusWaitingDriver {
		if err := s.redis.Publish(ctx, "orders", ord); err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderService) FinishOrder(ctx context.Context, id string) (err error) {
	defer derrors.Wrap(&err, "mongo.OrderService.FinishOrder")
	s.orderLock(id)
	defer s.orderUnlock(id)
	user := order.UserFromContext(ctx)
	if user == nil {
		return order.ErrAccessDenied
	}
	if user.Role != order.RoleDriver && user.Role != order.RoleAdmin {
		return order.ErrAccessDenied
	}
	ord, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return err
	}
	ord.Status = order.OrderStatusDropOff
	ord.EndAt = time.Now().UTC().Unix()
	if err = updateOrder(ctx, s.db, ord.ID, ord); err != nil {
		return err
	}
	// Create order charge
	switch ord.ChargeMethod {
	case order.ChargeMethodCash,
		order.ChargeMethodMLCTransaction,
		order.ChargeMethodCUPTransaction,
		order.ChargeMethodBalance:
		//TODO: create charge
		// charger, err := s.charge.Charge(ctx, order.ChargeMethod, currency.Amount{
		// 	Amount:   int64(order.Price),
		// 	Currency: currency.MustParse(order.Currency),
		// })
		// if err != nil {
		// 	return nil, err
		// }
		// ord.ChargeID = charger.ID
	default:
		return fmt.Errorf("unsupported charge method: %s", ord.ChargeMethod)
	}

	// TODO: send notification to rider that driver started the ride
	// TODO: update rider last location in the trip

	return nil
}

func (s *OrderService) StartOrder(ctx context.Context, id string) error {
	s.orderLock(id)
	defer s.orderUnlock(id)
	user := order.UserFromContext(ctx)
	if user == nil {
		return order.ErrAccessDenied
	}
	if user.Role != order.RoleDriver && user.Role != order.RoleAdmin {
		return order.ErrAccessDenied
	}
	ord, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return err
	}
	if ord.Driver != user.ID {
		return order.ErrAccessDenied
	}
	ord.Status = order.OrderStatusPickUp
	ord.StartAt = time.Now().UTC().Unix()
	if err = updateOrder(ctx, s.db, ord.ID, ord); err != nil {
		return err
	}
	return nil
}

// RateOrder implements order.OrderService.
func (s *OrderService) RateOrder(ctx context.Context, id string, rate float64, comment string) error {
	user := order.UserFromContext(ctx)
	if user == nil {
		return order.ErrAccessDenied
	}
	o, err := findOrderById(context.Background(), s.db, id)
	if err != nil {
		return err
	}
	if o.Rider != user.ID {
		return order.NewUnauthorized()
	}
	o.Rate = rate
	o.Review = comment
	if err := updateOrder(ctx, s.db, o.ID, o); err != nil {
		return err
	}
	return nil
}

// Categories implements order.OrderService.
func (s *OrderService) Categories(ctx context.Context, id string) ([]*order.CategoryPrice, error) {
	o, err := findOrderById(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return o.CategoryPrice, nil
}

func findOrders(ctx context.Context, db *DB, filter order.OrderFilter) ([]*order.Order, string, error) {
	collection := db.Collection(OrderCollection)
	user := order.UserFromContext(ctx)
	if user != nil {
		switch user.Role {
		case order.RoleRider:
			filter.Rider = user.ID
		case order.RoleDriver:
			filter.Driver = user.ID
		}
	}

	var trips []*order.Order
	var token string
	f := bson.D{}

	if filter.IDs != nil {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$in", Value: filter.IDs}}})
	}
	if filter.Rider != "" {
		f = append(f, bson.E{Key: "rider", Value: filter.Rider})
	}
	// TODO: Check if the order is linked with the driver or doesn't have driver assigned
	if filter.Driver != "" {
		f = append(f, bson.E{Key: "$or", Value: bson.D{
			{Key: "driver", Value: filter.Driver},
			{Key: "driver", Value: bson.D{{Key: "$exists", Value: false}}}}})
	}
	if filter.Status != "" {
		f = append(f, bson.E{Key: "status", Value: filter.Status})
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Token != "" {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: filter.Token}}})
	}

	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	for cur.Next(ctx) {
		var trip order.Order
		err := cur.Decode(&trip)
		if err != nil {
			return nil, "", err
		}
		trips = append(trips, &trip)
		if len(trips) == filter.Limit+1 && filter.Limit != 0 {
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

func findOrderById(ctx context.Context, db *DB, id string) (*order.Order, error) {
	usr := order.UserFromContext(ctx)
	if usr == nil {
		return nil, fmt.Errorf("nil user in context: %w", order.ErrAccessDenied)
	}
	filter := order.OrderFilter{
		IDs:   []string{id},
		Limit: 1,
	}
	if usr.Role == order.RoleRider {
		filter.Rider = usr.ID
	}
	if usr.Role == order.RoleDriver {
		filter.Driver = usr.ID
	}
	trips, _, err := findOrders(ctx, db, filter)
	if err != nil || len(trips) == 0 {
		return nil, order.ErrNotFound
	}
	return trips[0], nil
}

func updateOrder(ctx context.Context, db *DB, id string, o *order.Order) error {
	o.UpdatedAt = time.Now().UTC().Unix()
	collection := db.client.Database(database).Collection(OrderCollection.String())
	if _, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: o}}); err != nil {
		return fmt.Errorf("unabe to update the order: %v: %w", err, order.ErrInternal)
	}
	return nil
}

func storeOrder(ctx context.Context, db *DB, o *order.Order) error {
	collection := db.Collection(OrderCollection)
	if _, err := collection.InsertOne(ctx, o); err != nil {
		return fmt.Errorf("unable to store the order: %v: %w", err, order.ErrInternal)
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

func (s *OrderService) prepareOrder(ctx context.Context, o *order.Order, req order.DirectionRequest) (*order.Order, error) {
	if o == nil {
		o = &order.Order{
			ID:        order.NewID().String(),
			Status:    order.OrderStatusNew,
			CreatedAt: time.Now().UTC().Unix(),
		}
	}
	if o.ID == "" {
		o.ID = order.NewID().String()
	}
	var err error
	usr := order.UserFromContext(ctx)
	if usr == nil {
		return nil, order.ErrAccessDenied
	}
	o.Rider = usr.ID
	if req.Currency != "" {
		_, err = currency.Parse(req.Currency)
		if err != nil {
			return nil, fmt.Errorf("invalid currency: %v: %w", err, order.ErrInvalidCurrency)
		}
		o.Currency = req.Currency
	}
	if o.Currency == currency.XXX.String() {
		o.Currency = currency.CUP
	}
	if req.Riders == 0 {
		req.Riders = 1
	}

	o.Item = order.Item{
		Points:   req.Points,
		Riders:   req.Riders,
		Baggages: req.Baggages,
		Coupon:   req.Coupon,
		Currency: req.Currency,
	}
	routes, strBody, err := s.direction.GetRoute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to get the route: %w", err)
	}

	if len(routes.Routes) == 0 {
		return nil, fmt.Errorf("no routes found: %w", order.ErrNotFound)
	}

	o.Distance = routes.Distance
	o.Duration = routes.Duration
	o.Route = routes.Routes[0]
	o.RouteString = base64.StdEncoding.EncodeToString([]byte(strBody))

	err = s.CalculatePrice(o)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate the price: %w", err)
	}
	return o, nil
}

func checkRate(r *order.Rate) bool {
	currentTime := time.Now().UTC()
	if r.StartDate != "" && r.EndDate != "" {
		startDate, err := time.Parse("2006-01-02", r.StartDate)
		if err != nil {
			return false
		}
		endDate, err := time.Parse("2006-01-02", r.EndDate)
		if err != nil {
			return false
		}

		if (currentTime.Unix() < startDate.Unix() || currentTime.Unix() > startDate.Unix()) ||
			(r.EndDate != "" && currentTime.Unix() > endDate.Unix() || r.StartDate < r.EndDate) {
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

func price(distance, duration float64, rate order.Rate, riders int) float64 {
	price := float64(rate.BasePrice) +
		float64(rate.PricePerKm)*(float64(distance)/1000) +
		float64(rate.PricePerMin)*float64(duration)/60 +
		float64(rate.PricePerPassenger)*float64(riders)
	if rate.PricePerBaggage != 0 {
		price += float64(rate.PricePerBaggage)
	}
	if rate.PricePerCarryPet != 0 {
		price += float64(rate.PricePerCarryPet)
	}
	return price
}
