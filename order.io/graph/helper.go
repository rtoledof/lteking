package graph

import (
	"order.io/graph/model"
	"order.io/pkg/order"
)

func assembleOrderItem(input *model.RideInput) order.Item {
	item := order.Item{
		Points: assemblePoints(input.Item),
	}
	if input.Currency != nil {
		item.Currency = *input.Currency
	}
	if input.Coupon != nil {
		item.Coupon = *input.Coupon
	}
	if input.Riders != nil {
		item.Riders = *input.Riders
	}
	if input.Baggages != nil {
		item.Baggages = *input.Baggages
	}
	return item
}

func assemblePoint(input *model.PointInput) order.Point {
	return order.Point{
		Lat: input.Lat,
		Lng: input.Lng,
	}
}

func assemblePoints(input []*model.PointInput) []*order.Point {
	points := make([]*order.Point, len(input))
	for i, p := range input {
		point := assemblePoint(p)
		points[i] = &point
	}
	return points
}

func assembleModelOrder(o *order.Order) (*model.Order, error) {
	ord := &model.Order{
		ID:       o.ID,
		Items:    assembleModelItem(o.Item),
		Rate:     &o.Rate,
		Price:    &o.Price,
		Currency: &o.Currency,
		Route:    &o.RouteString,
		Distance: &o.Distance,
		Duration: &o.Duration,
		ChargeID: &o.ChargeID,
		History:  assembleModelPoints(o.History),
	}
	status := model.OrderStatus(o.Status)
	if !status.IsValid() {
		return nil, order.NewInvalidParameter("status", "invalid status")
	}
	ord.Status = &status
	pm := model.PaymentMethod(o.ChargeMethod)
	if !pm.IsValid() {
		return nil, order.NewInvalidParameter("payment_method", "invalid payment method")
	}
	ord.PaymentMethod = &pm

	return ord, nil
}

func assembleModelItem(item order.Item) *model.Item {
	return &model.Item{
		Points:   assembleModelPoints(item.Points),
		Baggages: &item.Baggages,
		Currency: &item.Currency,
		Coupon:   &item.Coupon,
		Riders:   &item.Riders,
	}
}

func assembleModelPoint(point *order.Point) *model.Point {
	return &model.Point{
		Lat: point.Lat,
		Lng: point.Lng,
	}
}

func assembleModelPoints(points []*order.Point) []*model.Point {
	modelPoints := make([]*model.Point, len(points))
	for i, p := range points {
		point := assembleModelPoint(p)
		modelPoints[i] = point
	}
	return modelPoints
}

func assembleOrderFilter(filter model.OrderListFilter) order.OrderFilter {
	f := order.OrderFilter{}
	f.Limit = 10
	if filter.Limit != nil {
		f.Limit = *filter.Limit
	}
	if filter.Status != nil {
		f.Status = order.OrderStatus(*filter.Status)
		if !f.Status.IsValid() {
			f.Status = ""
		}
	}
	if filter.Token != nil {
		f.Token = *filter.Token
	}

	return f
}

func assembleCategoryPrice(category order.VehicleCategory, price int, currency string) (*model.CategoryPrice, error) {
	catPrice := &model.CategoryPrice{
		Price:    float64(price),
		Currency: currency,
	}
	catPrice.Category = model.Category(category)
	if !catPrice.Category.IsValid() {
		return nil, order.NewInvalidParameter("category", "invalid category")
	}

	return catPrice, nil
}

func assembleCategoryPrices(categories []*order.CategoryPrice) []*model.CategoryPrice {
	catPrices := make([]*model.CategoryPrice, len(categories))
	for i, c := range categories {
		catPrice, _ := assembleCategoryPrice(c.Category, c.Price, c.Currency)
		catPrices[i] = catPrice
	}
	return catPrices
}
