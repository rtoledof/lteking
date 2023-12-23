package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

// {
// 	"waypoints": [
// 	  {
// 		"location": [
// 		  -84.518399,
// 		  39.134126
// 		],
// 		"name": ""
// 	  },
// 	  {
// 		"location": [
// 		  -84.511987,
// 		  39.102638
// 		],
// 		"name": "East 6th Street"
// 	  }
// 	],
// 	"routes": [
// 	  {
// 		"legs": [
// 		  {
// 			"steps": [],
// 			"weight": 1332.6,
// 			"distance": 4205,
// 			"summary": "",
// 			"duration": 1126
// 		  }
// 		],
// 		"weight_name": "cyclability",
// 		"geometry": {
// 		  "coordinates": [
// 			[
// 			  -84.518399,
// 			  39.134126
// 			],
// 			[
// 			  -84.51841,
// 			  39.133781
// 			],
// 			[
// 			  -84.520024,
// 			  39.133456
// 			],
// 			[
// 			  -84.520321,
// 			  39.132597
// 			],
// 			[
// 			  -84.52085,
// 			  39.128019
// 			],
// 			[
// 			  -84.52036,
// 			  39.127901
// 			],
// 			[
// 			  -84.52094,
// 			  39.122783
// 			],
// 			[
// 			  -84.52022,
// 			  39.122713
// 			],
// 			[
// 			  -84.520768,
// 			  39.120841
// 			],
// 			[
// 			  -84.519639,
// 			  39.120268
// 			],
// 			[
// 			  -84.51233,
// 			  39.114141
// 			],
// 			[
// 			  -84.512652,
// 			  39.11311
// 			],
// 			[
// 			  -84.512399,
// 			  39.112216
// 			],
// 			[
// 			  -84.513232,
// 			  39.112084
// 			],
// 			[
// 			  -84.512127,
// 			  39.107599
// 			],
// 			[
// 			  -84.512904,
// 			  39.107489
// 			],
// 			[
// 			  -84.511692,
// 			  39.102682
// 			],
// 			[
// 			  -84.511987,
// 			  39.102638
// 			]
// 		  ],
// 		  "type": "LineString"
// 		},
// 		"weight": 1332.6,
// 		"distance": 4205,
// 		"duration": 1126
// 	  }
// 	],
// 	"code": "Ok"
//   }

type OrderItem struct {
	PickUp  Point `json:"pick_up" bson:"pick_up"`
	DropOff Point `json:"drop_off" bson:"drop_off"`
}

type CategoryPrice struct {
	Category VehicleCategory `json:"category"`
	Price    uint64          `json:"price"`
}

type Order struct {
	ID               string                `json:"id" bson:"_id"`
	Items            DirectionRequest      `json:"items" bson:"items"`
	History          []*Point              `json:"history,omitempty" bson:"history,omitempty"`
	Driver           string                `json:"driver,omitempty" bson:"driver,omitempty"`
	Rider            string                `json:"rider" bson:"rider"`
	Status           OrderStatus           `json:"status" bson:"status"`
	StatusHistory    []*OrderStatusHistory `json:"status_history,omitempty" bson:"status_history,omitempty"`
	Rate             int                   `json:"rate" bson:"rate"`
	Price            uint64                `json:"price" bson:"price"`
	Coupon           string                `json:"coupon,omitempty" bson:"coupon,omitempty"`
	StartAt          int64                 `json:"start_at" bson:"start_at"`
	EndAt            int64                 `json:"end_at" bson:"end_at"`
	Review           string                `json:"review,omitempty" bson:"review"`
	CreatedAt        int64                 `json:"created_at" bson:"created_at"`
	UpdatedAt        int64                 `json:"updated_at" bson:"updated_at"`
	Route            *DirectionResponse    `json:"route,omitempty" bson:"route,omitempty"`
	Distance         float64               `json:"distance,omitempty" bson:"distance,omitempty"`
	Duration         float64               `json:"duration,omitempty" bson:"duration,omitempty"`
	SelectedCategory CategoryPrice         `json:"selected_category,omitempty" bson:"selected_category,omitempty"`
	CategoryPrice    []*CategoryPrice      `json:"categories_prices,omitempty" bson:"categories_prices,omitempty"`
	RouteString      string                `json:"route_string,omitempty" bson:"route_string,omitempty"`
}

type Item struct {
	PickUp  *PointInput   `json:"pick_up"`
	DropOff *PointInput   `json:"drop_off"`
	Seconds int           `json:"seconds"`
	M       float64       `json:"m"`
	Route   []*PointInput `json:"route"`
}

func AssambleOrderItem(items []*Item) []*OrderItem {
	var resp []*OrderItem
	for _, v := range items {
		var i = OrderItem{
			PickUp: Point{
				Lat: v.PickUp.Lat,
				Lng: v.PickUp.Lng,
			},
			DropOff: Point{
				Lat: v.DropOff.Lat,
				Lng: v.DropOff.Lng,
			},
		}
		resp = append(resp, &i)
	}
	return resp
}

type PointInput struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lon"`
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
	Limit  int       `json:"limit,omitempty"`
	Token  *string   `json:"token,omitempty"`
	Ids    []*string `json:"ids,omitempty"`
	Rider  *string   `json:"rider,omitempty"`
	Driver *string   `json:"driver,omitempty"`
	Status *string   `json:"status,omitempty"`
}

type OrderService interface {
	Create(context.Context, *DirectionRequest) (*Order, error)
	Update(context.Context, *DirectionRequest) (*Order, error)
	FindByID(context.Context, string) (*Order, error)
	FindAll(context.Context, *OrderFilter) (*OrderList, error)

	AcceptOrder(context.Context, string) (*Order, error)
	CancelOrder(context.Context, string) (*Order, error)
	CompleteOrder(context.Context, string) (*Order, error)
	StartOrder(context.Context, string) (*Order, error)
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
	OrderStatusNew       OrderStatus = "NEW"
	OrderStatusPickUp    OrderStatus = "PICKED_UP"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusOnTheWay  OrderStatus = "ON_THE_WAY"
	OrderStatusDropOff   OrderStatus = "DROPED_OFF"
	OrderStatusCancel    OrderStatus = "CANCELED"
)

var AllOrderStatus = []OrderStatus{
	OrderStatusNew,
	OrderStatusPickUp,
	OrderStatusConfirmed,
	OrderStatusOnTheWay,
	OrderStatusDropOff,
	OrderStatusCancel,
}

func (e OrderStatus) IsValid() bool {
	switch e {
	case OrderStatusNew, OrderStatusPickUp, OrderStatusOnTheWay, OrderStatusDropOff, OrderStatusCancel:
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
	Route    []*PointInput `json:"route"`
	Coupon   *string       `json:"coupon,omitempty"`
	Riders   *int          `json:"riders,omitempty"`
	Baggages *int          `json:"baggages,omitempty"`
}

type CreateOrderResponse struct {
	Order *Order          `json:"order"`
	Cost  []*CostPerBrand `json:"cost"`
	Price int             `json:"price"`
}
