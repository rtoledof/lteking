package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type OrderItem struct {
	ID      string  `json:"id" bson:"_id"`
	PickUp  Point   `json:"pick_up" bson:"pick_up"`
	DropOff Point   `json:"drop_off" bson:"drop_off"`
	Route   []Point `json:"route,omitempty" bson:"route,omitempty"`
	Seconds uint64  `json:"seconds,omitempty" bson:"seconds,omitempty"`
	Meters  uint64  `json:"meters,omitempty" bson:"meters,omitempty"`
}

type Order struct {
	ID            string                `json:"id" bson:"_id"`
	Items         []OrderItem           `json:"items" bson:"items"`
	History       []Point               `json:"history,omitempty" bson:"history,omitempty"`
	Driver        string                `json:"driver,omitempty" bson:"driver,omitempty"`
	Rider         string                `json:"rider" bson:"rider"`
	Status        OrderStatus           `json:"status" bson:"status"`
	StatusHistory []*OrderStatusHistory `json:"status_history,omitempty" bson:"status_history,omitempty"`
	Rate          int                   `json:"rate" bson:"rate"`
	Price         uint64                `json:"price" bson:"price"`
	Coupon        string                `json:"coupon,omitempty" bson:"coupon,omitempty"`
	StartAt       int                   `json:"start_at" bson:"start_at"`
	EndAt         int                   `json:"end_at" bson:"end_at"`
	Review        string                `json:"review,omitempty" bson:"review"`
	CreatedAt     int64                 `json:"created_at" bson:"created_at"`
	UpdatedAt     int64                 `json:"updated_at" bson:"updated_at"`
}

type Item struct {
	PickUp  *PointInput   `json:"pick_up"`
	DropOff *PointInput   `json:"drop_off"`
	Seconds int           `json:"seconds"`
	M       float64       `json:"m"`
	Route   []*PointInput `json:"route"`
}

type PointInput struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
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
	Limit  *int      `json:"limit,omitempty"`
	Token  *string   `json:"token,omitempty"`
	Ids    []*string `json:"ids,omitempty"`
	Rider  *string   `json:"rider,omitempty"`
	Driver *string   `json:"driver,omitempty"`
	Status *string   `json:"status,omitempty"`
}

type OrderService interface {
	Create(context.Context, []OrderItem) (*Order, error)
	Update(context.Context, *UpdateOrder) (*Order, error)
	FindByID(context.Context, string) (*Order, error)
	FindAll(context.Context, *OrderFilter) (*OrderList, error)
}

type AddPlace struct {
	Name     string         `json:"name"`
	Location *LocationInput `json:"location"`
}

type OrderStatusHistory struct {
	Status    OrderStatus `json:"status" bson:"status"`
	ChangedAt string      `json:"changed_at" bson:"changed_at"`
}

type OrderStatus string

const (
	OrderStatusNew      OrderStatus = "NEW"
	OrderStatusPickUp   OrderStatus = "PICK_UP"
	OrderStatusOnTheWay OrderStatus = "ON_THE_WAY"
	OrderStatusDropOff  OrderStatus = "DROP_OFF"
)

var AllOrderStatus = []OrderStatus{
	OrderStatusNew,
	OrderStatusPickUp,
	OrderStatusOnTheWay,
	OrderStatusDropOff,
}

func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusNew, OrderStatusPickUp, OrderStatusOnTheWay, OrderStatusDropOff:
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
