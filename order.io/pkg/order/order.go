package order

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Point struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lon" bson:"lon"`
}

func (p *Point) String() string {
	return fmt.Sprintf("%f,%f", p.Lng, p.Lat)
}

func (p *Point) Valid() bool {
	return p.Lat >= -90 && p.Lat <= 90 && p.Lng >= -180 && p.Lng <= 180
}

type Item struct {
	Points   []*Point `json:"points" bson:"points"`
	Coupon   string   `json:"coupon" bson:"coupon"`
	Riders   int      `json:"riders" bson:"riders"`
	Baggages bool     `json:"baggages" bson:"baggages"`
	Currency string   `json:"currency,omitempty" bson:"currency,omitempty"`
}

type CategoryPrice struct {
	Category VehicleCategory `json:"category"`
	Price    int             `json:"price,omitempty"`
	Currency string          `json:"currency,omitempty"`
}

type Order struct {
	ID               string                `json:"id" bson:"_id"`
	Item             Item                  `json:"item" bson:"item"`
	History          []*Point              `json:"history,omitempty" bson:"history,omitempty"`
	Driver           string                `json:"driver,omitempty" bson:"driver,omitempty"`
	Rider            string                `json:"rider" bson:"rider"`
	Status           OrderStatus           `json:"status" bson:"status"`
	StatusHistory    []*OrderStatusHistory `json:"status_history,omitempty" bson:"status_history,omitempty"`
	Rate             float64               `json:"rate" bson:"rate"`
	Price            int                   `json:"price,omitempty" bson:"price,omitempty"`
	Currency         string                `json:"currency,omitempty" bson:"currency,omitempty"`
	Coupon           string                `json:"coupon,omitempty" bson:"coupon,omitempty"`
	StartAt          int64                 `json:"start_at" bson:"start_at"`
	EndAt            int64                 `json:"end_at" bson:"end_at"`
	Review           string                `json:"review,omitempty" bson:"review"`
	CreatedAt        int64                 `json:"created_at" bson:"created_at"`
	UpdatedAt        int64                 `json:"updated_at" bson:"updated_at"`
	Route            *Route                `json:"route,omitempty" bson:"route,omitempty"`
	Distance         float64               `json:"distance,omitempty" bson:"distance,omitempty"`
	Duration         float64               `json:"duration,omitempty" bson:"duration,omitempty"`
	SelectedCategory *CategoryPrice        `json:"selected_category,omitempty" bson:"selected_category,omitempty"`
	CategoryPrice    []*CategoryPrice      `json:"categories_prices,omitempty" bson:"categories_prices,omitempty"`
	RouteString      string                `json:"route_string,omitempty" bson:"route_string,omitempty"`
	ChargeMethod     ChargeMethod          `json:"charge_method,omitempty" bson:"charge_method,omitempty"`
	ChargeID         string                `json:"charge_id,omitempty" bson:"charge_id,omitempty"`
	BannedDrivers    map[string]bool       `json:"banned_drivers,omitempty" bson:"banned_drivers,omitempty"`
}

func AssambleOrderItem(items *Item) Item {
	return Item{
		Points:   items.Points,
		Riders:   items.Riders,
		Baggages: items.Baggages,
		Currency: items.Currency,
	}
}

type UpdateOrder struct {
	Driver *string      `json:"driver,omitempty"`
	Items  []*Item      `json:"items,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
}

type OrderList struct {
	Token string   `json:"token"`
	Data  []*Order `json:"data"`
}

type OrderFilter struct {
	Limit  int         `json:"limit,omitempty"`
	Token  string      `json:"token,omitempty"`
	IDs    []string    `json:"ids,omitempty"`
	Rider  string      `json:"rider,omitempty"`
	Driver string      `json:"driver,omitempty"`
	Status OrderStatus `json:"status,omitempty"`
}

type OrderService interface {
	Create(context.Context, Item) (*Order, error)
	Update(context.Context, string, Item) (*Order, error)
	FindAll(context.Context, OrderFilter) (*OrderList, error)
	FindByID(context.Context, string) (*Order, error)

	ConfirmOrder(context.Context, ConfirmOrder) error

	AcceptOrder(context.Context, string) error
	StartOrder(context.Context, string) error
	CancelOrder(context.Context, string) error
	FinishOrder(context.Context, string) error
	RateOrder(context.Context, string, float64, string) error

	Categories(context.Context, string) ([]*CategoryPrice, error)
}

type AddPlace struct {
	Name  string `json:"name"`
	Point *Point `json:"point"`
}

type OrderStatusHistory struct {
	Status    OrderStatus `json:"status" bson:"status"`
	ChangedAt string      `json:"changed_at" bson:"changed_at"`
}

type OrderStatus string

const (
	OrderStatusNew           OrderStatus = "NEW"
	OrderStatusWaitingDriver OrderStatus = "WAITING_DRIVER"
	OrderStatusPickUp        OrderStatus = "PICKED_UP"
	OrderStatusConfirmed     OrderStatus = "CONFIRMED"
	OrderStatusOnTheWay      OrderStatus = "ON_THE_WAY"
	OrderStatusDropOff       OrderStatus = "DROPED_OFF"
	OrderStatusCancel        OrderStatus = "CANCELED"
)

var AllOrderStatus = []OrderStatus{
	OrderStatusNew,
	OrderStatusWaitingDriver,
	OrderStatusPickUp,
	OrderStatusConfirmed,
	OrderStatusOnTheWay,
	OrderStatusDropOff,
	OrderStatusCancel,
}

func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusNew, OrderStatusPickUp, OrderStatusOnTheWay, OrderStatusDropOff, OrderStatusCancel, OrderStatusWaitingDriver:
		return true
	}
	return false
}

func (e OrderStatus) String() string {
	return string(e)
}

func (e *OrderStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrderStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrderStatus", str)
	}
	return nil
}

func (e OrderStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type CostPerBrand struct {
	ID    string `json:"id"`
	Brand Brand  `json:"brand"`
	Price int    `json:"price"`
}

type CreateOrderRequest struct {
	Route    []*Point `json:"route"`
	Coupon   *string  `json:"coupon,omitempty"`
	Riders   *int     `json:"riders,omitempty"`
	Baggages *int     `json:"baggages,omitempty"`
}

type CreateOrderResponse struct {
	Order *Order          `json:"order"`
	Cost  []*CostPerBrand `json:"cost"`
	Price int             `json:"price"`
}

type ConfirmOrder struct {
	OrderID  string          `json:"order_id"`
	Category VehicleCategory `json:"category"`
	Method   ChargeMethod    `json:"method"`
	Currency string          `json:"currency"`
}
