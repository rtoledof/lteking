package cubawheeler

import "context"

type Statistics struct {
	Total  int      `json:"total"`
	Amount uint64   `json:"amount"`
	User   string   `json:"user"`
	Orders []string `json:"orders"`
}

type OrderStatistics struct {
	ID     string                 `json:"id" bson:"_id"`
	User   string                 `json:"user" bson:"user"`
	Orders map[string]*Statistics `json:"orders" bson:"orders"`
}

func (o *OrderStatistics) AddOrder(order Order, date Time) {
	if o.Orders == nil {
		o.Orders = make(map[string]*Statistics)
	}
	formats := []string{"2006-01-02", "2006-01", "2006"}
	for _, f := range formats {
		interval := date.Format(f)
		if _, ok := o.Orders[interval]; !ok {
			o.Orders[interval] = &Statistics{
				User:   o.User,
				Orders: []string{},
				Total:  0,
				Amount: 0,
			}
		}
		o.Orders[interval].Orders = append(o.Orders[interval].Orders, order.ID)
		o.Orders[interval].Total++
		o.Orders[interval].Amount += order.Price
	}
}

type OrderStatisticsFilter struct {
	Limit int    `json:"limit"`
	Token string `json:"token"`
	User  string `json:"user"`
}

type OrderStatisticsService interface {
	AddOrder(context.Context, OrderStatistics) error
	FindStatistictsByUser(context.Context, string) (*OrderStatistics, error)
	FindAllStatistics(context.Context, OrderStatisticsFilter) ([]*OrderStatistics, string, error)
}
