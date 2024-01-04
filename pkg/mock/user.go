package mock

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.UserService = &UserService{}

type UserService struct {
	AddDeviceFn          func(context.Context, string) error
	AddFavoritePlaceFn   func(context.Context, cubawheeler.AddPlace) (*cubawheeler.Location, error)
	AddFavoriteVehicleFn func(context.Context, *string) (*cubawheeler.Vehicle, error)
	CreateUserFn         func(context.Context, *cubawheeler.User) error
	FavoritePlacesFn     func(context.Context) ([]*cubawheeler.Location, error)
	FavoriteVehiclesFn   func(context.Context) ([]*cubawheeler.Vehicle, error)
	FindAllFn            func(context.Context, *cubawheeler.UserFilter) (*cubawheeler.UserList, error)
	FindByEmailFn        func(context.Context, string) (*cubawheeler.User, error)
	FindByIDFn           func(context.Context, string) (*cubawheeler.User, error)
	GetUserDevicesFn     func(context.Context, []string) ([]string, error)
	LastNAddressFn       func(context.Context, int) ([]*cubawheeler.Location, error)
	LoginFn              func(context.Context, cubawheeler.LoginRequest) (*cubawheeler.User, error)
	MeFn                 func(context.Context) (*cubawheeler.Profile, error)
	OrdersFn             func(context.Context, *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error)
	SetAvailabilityFn    func(context.Context, string, bool) error
	UpdateFn             func(context.Context, *cubawheeler.User) error
	UpdatePlaceFn        func(context.Context, *cubawheeler.UpdatePlace) (*cubawheeler.Location, error)
	UpdateProfileFn      func(context.Context, *cubawheeler.UpdateProfile) error
}

// AddDevice implements cubawheeler.UserService.
func (s *UserService) AddDevice(ctx context.Context, device string) error {
	return s.AddDeviceFn(ctx, device)
}

// AddFavoritePlace implements cubawheeler.UserService.
func (s *UserService) AddFavoritePlace(ctx context.Context, place cubawheeler.AddPlace) (*cubawheeler.Location, error) {
	return s.AddFavoritePlaceFn(ctx, place)
}

// AddFavoriteVehicle implements cubawheeler.UserService.
func (s *UserService) AddFavoriteVehicle(ctx context.Context, vehicle *string) (*cubawheeler.Vehicle, error) {
	return s.AddFavoriteVehicleFn(ctx, vehicle)
}

// CreateUser implements cubawheeler.UserService.
func (s *UserService) CreateUser(ctx context.Context, user *cubawheeler.User) error {
	return s.CreateUserFn(ctx, user)
}

// FavoritePlaces implements cubawheeler.UserService.
func (s *UserService) FavoritePlaces(ctx context.Context) ([]*cubawheeler.Location, error) {
	return s.FavoritePlacesFn(ctx)
}

// FavoriteVehicles implements cubawheeler.UserService.
func (s *UserService) FavoriteVehicles(ctx context.Context) ([]*cubawheeler.Vehicle, error) {
	return s.FavoriteVehiclesFn(ctx)
}

// FindAll implements cubawheeler.UserService.
func (s *UserService) FindAll(ctx context.Context, filter *cubawheeler.UserFilter) (*cubawheeler.UserList, error) {
	return s.FindAllFn(ctx, filter)
}

// FindByEmail implements cubawheeler.UserService.
func (s *UserService) FindByEmail(ctx context.Context, email string) (*cubawheeler.User, error) {
	return s.FindByEmailFn(ctx, email)
}

// FindByID implements cubawheeler.UserService.
func (s *UserService) FindByID(ctz context.Context, id string) (*cubawheeler.User, error) {
	return s.FindByIDFn(ctz, id)
}

// GetUserDevices implements cubawheeler.UserService.
func (s *UserService) GetUserDevices(ctx context.Context, devices []string) ([]string, error) {
	return s.GetUserDevicesFn(ctx, devices)
}

// LastNAddress implements cubawheeler.UserService.
func (s *UserService) LastNAddress(ctz context.Context, limit int) ([]*cubawheeler.Location, error) {
	return s.LastNAddressFn(ctz, limit)
}

// Login implements cubawheeler.UserService.
func (s *UserService) Login(ctx context.Context, req cubawheeler.LoginRequest) (*cubawheeler.User, error) {
	return s.LoginFn(ctx, req)
}

// Me implements cubawheeler.UserService.
func (s *UserService) Me(ctx context.Context) (*cubawheeler.Profile, error) {
	return s.MeFn(ctx)
}

// Orders implements cubawheeler.UserService.
func (s *UserService) Orders(ctx context.Context, filter *cubawheeler.OrderFilter) (*cubawheeler.OrderList, error) {
	return s.OrdersFn(ctx, filter)
}

// SetAvailability implements cubawheeler.UserService.
func (s *UserService) SetAvailability(ctx context.Context, user string, available bool) error {
	return s.AddDeviceFn(ctx, user)
}

// Update implements cubawheeler.UserService.
func (s *UserService) Update(ctx context.Context, user *cubawheeler.User) error {
	return s.UpdateFn(ctx, user)
}

// UpdatePlace implements cubawheeler.UserService.
func (s *UserService) UpdatePlace(ctx context.Context, place *cubawheeler.UpdatePlace) (*cubawheeler.Location, error) {
	return s.UpdatePlaceFn(ctx, place)
}

// UpdateProfile implements cubawheeler.UserService.
func (s *UserService) UpdateProfile(ctx context.Context, Profile *cubawheeler.UpdateProfile) error {
	return s.UpdateProfileFn(ctx, Profile)
}
