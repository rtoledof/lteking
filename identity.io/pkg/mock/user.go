package mock

import (
	"context"

	"identity.io/pkg/identity"
)

var _ identity.UserService = &UserService{}

type UserService struct {
	AddDeviceFn             func(context.Context, string) error
	AddFavoritePlaceFn      func(context.Context, string, identity.Point) (*identity.Location, error)
	AddFavoriteVehicleFn    func(context.Context, string, *string) error
	CreateUserFn            func(context.Context, *identity.User) error
	FavoritePlacesFn        func(context.Context) ([]*identity.Location, error)
	FavoriteVehiclesFn      func(context.Context) ([]string, error)
	FindAllFn               func(context.Context, *identity.UserFilter) (*identity.UserList, error)
	FindByEmailFn           func(context.Context, string) (*identity.User, error)
	FindByIDFn              func(context.Context, string) (*identity.User, error)
	GetUserDevicesFn        func(context.Context, identity.UserFilter) ([]string, error)
	LastNAddressFn          func(context.Context, int) ([]*identity.Location, error)
	LoginFn                 func(context.Context, string, string, ...string) (*identity.User, error)
	MeFn                    func(context.Context) (*identity.Profile, error)
	SetAvailabilityFn       func(context.Context, bool) error
	UpdateFn                func(context.Context, *identity.User) error
	UpdatePlaceFn           func(context.Context, *identity.UpdatePlace) (*identity.Location, error)
	UpdateProfileFn         func(context.Context, *identity.UpdateProfile) error
	TokenFn                 func(context.Context, *identity.User) (string, error)
	LogoutFn                func(context.Context) error
	AddVehicleFn            func(context.Context, *identity.Vehicle) error
	UpdateVehicleFn         func(context.Context, *identity.Vehicle) error
	DeleteVehicleFn         func(context.Context, string) error
	AddDeviceTokenFn        func(context.Context, string, string) error
	DeleteFavoriteVehicleFn func(context.Context, string) error
	DeviceTokensFn          func(context.Context) ([]string, error)
	FavoritePlaceFn         func(context.Context, string) (*identity.Location, error)
	RemoveDeviceTokenFn     func(context.Context, string) error
	SetActiveVehicleFn      func(context.Context, string) error
	SetPreferedCurrencyFn   func(context.Context, string) error
	VehicleFn               func(context.Context, string) (*identity.Vehicle, error)
	VehiclesFn              func(context.Context) ([]*identity.Vehicle, error)
	DeleteFavoritePlaceFn   func(context.Context, string) error
	UpdateFavoritePlaceFn   func(context.Context, string, identity.UpdatePlace) error
}

// AddDeviceToken implements identity.UserService.
func (s *UserService) AddDeviceToken(ctx context.Context, token string, name string) error {
	return s.AddDeviceTokenFn(ctx, token, name)
}

// DeleteFavoriteVehicle implements identity.UserService.
func (s *UserService) DeleteFavoriteVehicle(ctx context.Context, v string) error {
	return s.DeleteFavoriteVehicleFn(ctx, v)
}

// DeviceTokens implements identity.UserService.
func (s *UserService) DeviceTokens(ctx context.Context) ([]string, error) {
	return s.DeviceTokensFn(ctx)
}

// FavoritePlace implements identity.UserService.
func (s *UserService) FavoritePlace(ctx context.Context, place string) (*identity.Location, error) {
	return s.FavoritePlaceFn(ctx, place)
}

// RemoveDeviceToken implements identity.UserService.
func (s *UserService) RemoveDeviceToken(ctx context.Context, token string) error {
	return s.RemoveDeviceTokenFn(ctx, token)
}

// SetActiveVehicle implements identity.UserService.
func (s *UserService) SetActiveVehicle(ctx context.Context, plate string) error {
	return s.SetActiveVehicleFn(ctx, plate)
}

// SetPreferedCurrency implements identity.UserService.
func (s *UserService) SetPreferedCurrency(ctx context.Context, cur string) error {
	return s.SetPreferedCurrencyFn(ctx, cur)
}

// Vehicle implements identity.UserService.
func (s *UserService) Vehicle(ctx context.Context, palte string) (*identity.Vehicle, error) {
	return s.VehicleFn(ctx, palte)
}

// Vehicles implements identity.UserService.
func (s *UserService) Vehicles(ctx context.Context) ([]*identity.Vehicle, error) {
	return s.VehiclesFn(ctx)
}

// DeleteFavoritePlace implements identity.UserService.
func (s *UserService) DeleteFavoritePlace(ctx context.Context, place string) error {
	return s.DeleteFavoritePlaceFn(ctx, place)
}

// UpdateFavoritePlace implements identity.UserService.
func (s *UserService) UpdateFavoritePlace(ctx context.Context, palce string, update identity.UpdatePlace) error {
	return s.UpdateFavoritePlaceFn(ctx, palce, update)
}

// DeleteVehicle implements identity.UserService.
func (s *UserService) DeleteVehicle(ctx context.Context, id string) error {
	return s.DeleteVehicleFn(ctx, id)
}

// UpdateVehicle implements identity.UserService.
func (s *UserService) UpdateVehicle(ctx context.Context, v *identity.Vehicle) error {
	return s.UpdateVehicleFn(ctx, v)
}

// AddVehicle implements identity.UserService.
func (s *UserService) AddVehicle(ctx context.Context, v *identity.Vehicle) error {
	return s.AddVehicleFn(ctx, v)
}

// Logout implements identity.UserService.
func (s *UserService) Logout(ctx context.Context) error {
	return s.LogoutFn(ctx)
}

// Token implements identity.UserService.
func (s *UserService) Token(ctx context.Context, user *identity.User) (string, error) {
	return s.TokenFn(ctx, user)
}

// AddDevice implements identity.UserService.
func (s *UserService) AddDevice(ctx context.Context, device string) error {
	return s.AddDeviceFn(ctx, device)
}

// AddFavoritePlace implements identity.UserService.
func (s *UserService) AddFavoritePlace(ctx context.Context, name string, point identity.Point) (*identity.Location, error) {
	return s.AddFavoritePlaceFn(ctx, name, point)
}

// AddFavoriteVehicle implements identity.UserService.
func (s *UserService) AddFavoriteVehicle(ctx context.Context, vehicle string, name *string) error {
	return s.AddFavoriteVehicleFn(ctx, vehicle, name)
}

// CreateUser implements identity.UserService.
func (s *UserService) CreateUser(ctx context.Context, user *identity.User) error {
	return s.CreateUserFn(ctx, user)
}

// FavoritePlaces implements identity.UserService.
func (s *UserService) FavoritePlaces(ctx context.Context) ([]*identity.Location, error) {
	return s.FavoritePlacesFn(ctx)
}

// FavoriteVehicles implements identity.UserService.
func (s *UserService) FavoriteVehicles(ctx context.Context) ([]string, error) {
	return s.FavoriteVehiclesFn(ctx)
}

// FindAll implements identity.UserService.
func (s *UserService) FindAll(ctx context.Context, filter *identity.UserFilter) (*identity.UserList, error) {
	return s.FindAllFn(ctx, filter)
}

// FindByEmail implements identity.UserService.
func (s *UserService) FindByEmail(ctx context.Context, email string) (*identity.User, error) {
	return s.FindByEmailFn(ctx, email)
}

// FindByID implements identity.UserService.
func (s *UserService) FindByID(ctz context.Context, id string) (*identity.User, error) {
	return s.FindByIDFn(ctz, id)
}

// GetUserDevices implements identity.UserService.
func (s *UserService) GetUserDevices(ctx context.Context, filter identity.UserFilter) ([]string, error) {
	return s.GetUserDevicesFn(ctx, filter)
}

// LastNAddress implements identity.UserService.
func (s *UserService) LastNAddress(ctz context.Context, limit int) ([]*identity.Location, error) {
	return s.LastNAddressFn(ctz, limit)
}

// Login implements identity.UserService.
func (s *UserService) Login(ctx context.Context, email, otp string, role ...string) (*identity.User, error) {
	return s.LoginFn(ctx, email, otp)
}

// Me implements identity.UserService.
func (s *UserService) Me(ctx context.Context) (*identity.Profile, error) {
	return s.MeFn(ctx)
}

// SetAvailability implements identity.UserService.
func (s *UserService) SetAvailability(ctx context.Context, available bool) error {
	return s.SetAvailabilityFn(ctx, available)
}

// Update implements identity.UserService.
func (s *UserService) Update(ctx context.Context, user *identity.User) error {
	return s.UpdateFn(ctx, user)
}

// UpdatePlace implements identity.UserService.
func (s *UserService) UpdatePlace(ctx context.Context, place *identity.UpdatePlace) (*identity.Location, error) {
	return s.UpdatePlaceFn(ctx, place)
}

// UpdateProfile implements identity.UserService.
func (s *UserService) UpdateProfile(ctx context.Context, Profile *identity.UpdateProfile) error {
	return s.UpdateProfileFn(ctx, Profile)
}
