package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"cubawheeler.io/cmd/driver/graph/model"
	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/currency"
	"cubawheeler.io/pkg/redis"
)

type OrderHandler struct {
	Service cubawheeler.OrderService
	Redis   *redis.Redis
}

func NewOrderHandler(service cubawheeler.OrderService, redis *redis.Redis) *OrderHandler {
	return &OrderHandler{
		Service: service,
		Redis:   redis,
	}
}

func (o *OrderHandler) Create(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Creating order")
	if !canDo(r, cubawheeler.RoleRider) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	req, err := processRequest(r)
	if err != nil {
		if err, ok := err.(*cubawheeler.Error); ok {
			w.WriteHeader(err.StatusCode)
		}
		return err
	}
	order, err := o.Service.Create(r.Context(), req)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order created: %s", order.ID))
	return json.NewEncoder(w).Encode(order)
}

func (o *OrderHandler) Update(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Updating order")
	if !canDo(r, cubawheeler.RoleRider) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	idParams := chi.URLParam(r, "id")
	order, err := o.Service.FindByID(r.Context(), idParams)
	if err != nil {
		return err
	}
	req, err := processRequest(r)
	if err != nil {
		return err
	}
	req.ID = order.ID
	order, err = o.Service.Update(r.Context(), req)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order updated: %s", order.ID))
	return json.NewEncoder(w).Encode(assambleOrder(order))
}

func (o *OrderHandler) List(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Listing orders")
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	filter, err := getOrderFilter(r)
	if err != nil {
		return err
	}
	if user.Role == cubawheeler.RoleDriver {
		filter.Driver = &user.ID
	}
	if user.Role == cubawheeler.RoleRider {
		filter.Rider = &user.ID
	}
	orders, err := o.Service.FindAll(r.Context(), filter)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Orders listed: %d", len(orders.Data)))
	return json.NewEncoder(w).Encode(orders)
}

func (o *OrderHandler) FindByID(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Finding order")
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	idParams := chi.URLParam(r, "id")
	order, err := o.Service.FindByID(r.Context(), idParams)
	if err != nil {
		return err
	}
	if (user.Role == cubawheeler.RoleDriver && (order.Driver != "" && order.Driver != user.ID)) ||
		(user.Role == cubawheeler.RoleRider && order.Rider != user.ID) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	logger.Info(fmt.Sprintf("Order found: %s", order.ID))
	return json.NewEncoder(w).Encode(order)
}

func (o *OrderHandler) Accept(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Accepting order")
	if !canDo(r, cubawheeler.RoleDriver) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	idParams := chi.URLParam(r, "id")
	order, err := o.Service.AcceptOrder(r.Context(), idParams)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order accepted: %s", order.ID))
	return json.NewEncoder(w).Encode(order)
}

func (o *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Canceling order")
	if !canDo(r, []cubawheeler.Role{cubawheeler.RoleRider, cubawheeler.RoleDriver}...) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	idParams := chi.URLParam(r, "id")
	order, err := o.Service.CancelOrder(r.Context(), idParams)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order canceled: %s", order.ID))
	return json.NewEncoder(w).Encode(order)
}

func (o *OrderHandler) Complete(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Completing order")
	if !canDo(r, cubawheeler.RoleDriver) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	idParams := chi.URLParam(r, "id")
	order, err := o.Service.FinishOrder(r.Context(), idParams)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order completed: %s", order.ID))
	return json.NewEncoder(w).Encode(order)
}

func (o *OrderHandler) Start(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Starting order")
	if !canDo(r, cubawheeler.RoleDriver) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	idParams := chi.URLParam(r, "id")
	order, err := o.Service.StartOrder(r.Context(), idParams)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order started: %s", order.ID))
	return json.NewEncoder(w).Encode(order)
}

func (o *OrderHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	logger := cannon.LoggerFromContext(r.Context())
	logger.Info("Confirming order")
	if !canDo(r, cubawheeler.RoleRider) {
		return cubawheeler.NewError(nil, http.StatusForbidden, "you are not allowed to do this")
	}
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %v: %w", err, cubawheeler.ErrInvalidInput)
	}
	order := chi.URLParam(r, "id")
	var req = cubawheeler.ConfirmOrder{
		OrderID:  order,
		Currency: r.FormValue("currency"),
	}
	cat := r.FormValue("category")
	if len(cat) == 0 {
		return cubawheeler.NewInvalidParameter("category", "category is required")
	}
	req.Category = cubawheeler.VehicleCategory(cat)
	if !req.Category.IsValid() {
		return cubawheeler.NewInvalidParameter("category", "category is not valid")
	}
	method := r.FormValue("method")
	if len(method) == 0 {
		return cubawheeler.NewInvalidParameter("method", "method is required")
	}
	req.Method = cubawheeler.ChargeMethod(method)
	if !req.Method.IsValid() {
		return cubawheeler.NewInvalidParameter("method", "method is not valid")
	}
	if err := o.Service.ConfirmOrder(r.Context(), req); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Order confirmed: %s", order))
	return nil
}

func processRequest(r *http.Request) (*cubawheeler.DirectionRequest, error) {
	var req cubawheeler.DirectionRequest
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("failed to parse form: %v: %w", err, cubawheeler.ErrInvalidInput)
	}
	if r.FormValue("coupon") != "" {
		req.Coupon = r.FormValue("coupon")
	}
	if r.FormValue("currency") != "" {
		req.Currency = r.FormValue("currency")
	}
	if r.FormValue("riders") != "" {
		riders, err := strconv.Atoi(r.FormValue("riders"))
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("riders", "riders is not valid")
		}
		req.Riders = riders
	}
	if r.FormValue("baggages") != "" {
		baggages, err := strconv.ParseBool(r.FormValue("baggages"))
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("baggages", "baggages is not valid")
		}
		req.Baggages = baggages
	}
	points, err := getPoints(r)
	if err != nil {
		return nil, err
	}
	req.Points = points
	if len(req.Points) == 0 {
		return nil, cubawheeler.NewInvalidParameter("items", "items is required")
	}
	for index, p := range req.Points {
		if !p.Valid() {
			return nil, cubawheeler.NewInvalidParameter(fmt.Sprintf("items[%d]", index), "invalid point")
		}
	}
	if req.Riders < 0 || req.Riders > 6 {
		return nil, cubawheeler.NewInvalidParameter("riders", "riders must be between 0 and 6")
	}
	if req.Currency != "" {
		if _, err := currency.Parse(req.Currency); err != nil {
			return nil, cubawheeler.NewInvalidParameter("currency", "currency is not valid")
		}
	}
	return &req, nil
}

func getOrderFilter(r *http.Request) (_ *cubawheeler.OrderFilter, err error) {
	var filter cubawheeler.OrderFilter
	if l := r.URL.Query().Get("limit"); len(l) > 0 {
		filter.Limit, err = strconv.Atoi(l)
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("limit", "limit is not valid")
		}
	}
	if s := r.URL.Query().Get("status"); len(s) > 0 {
		filter.Status = &s
	}
	if d := r.URL.Query().Get("driver"); len(d) > 0 {
		filter.Driver = &d
	}
	if ids := r.URL.Query()["ids"]; len(ids) > 0 {
		filter.IDs = ids
	}
	if token := r.URL.Query().Get("token"); len(token) > 0 {
		filter.Token = &token
	}

	return &filter, nil
}

func assamblePoints(points []*cubawheeler.Point) []*model.Point {
	var result []*model.Point
	for _, p := range points {
		result = append(result, &model.Point{
			Lat: p.Lat,
			Lng: p.Lng,
		})
	}
	return result
}

func AssambleOrderItem(order *cubawheeler.Order) *model.Item {
	return &model.Item{
		Points:   assamblePoints(order.Items.Points),
		Baggages: &order.Items.Baggages,
		Riders:   &order.Items.Riders,
		Currency: &order.Items.Currency,
	}
}

func assambleOrder(order *cubawheeler.Order) *model.Order {
	return &model.Order{
		ID:     order.ID,
		Status: order.Status.String(),
		Item:   AssambleOrderItem(order),
		Price: &model.Amount{
			Amount:   int(order.Price),
			Currency: order.Currency,
		},
	}
}

func assambleOrders(orders *cubawheeler.OrderList) []*model.Order {
	var result []*model.Order
	for _, o := range orders.Data {
		result = append(result, assambleOrder(o))
	}
	return result
}

func getPoints(r *http.Request) ([]*cubawheeler.Point, error) {
	points := r.Form["points"]
	if len(points) == 0 {
		return nil, cubawheeler.NewInvalidParameter("points", "points is required")
	}
	var result []*cubawheeler.Point
	for _, p := range points {
		var point cubawheeler.Point
		strLat, strLong := strings.Split(p, ",")[0], strings.Split(p, ",")[1]
		lat, err := strconv.ParseFloat(strLat, 64)
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("points", "invalid point")
		}
		long, err := strconv.ParseFloat(strLong, 64)
		if err != nil {
			return nil, cubawheeler.NewInvalidParameter("points", "invalid point")
		}
		point.Lat = lat
		point.Lng = long
		result = append(result, &point)
	}

	return result, nil
}
