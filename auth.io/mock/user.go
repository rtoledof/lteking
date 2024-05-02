package mock

import (
	"context"

	"auth.io/models"
)

var _ models.UserService = &UserService{}

type UserService struct {
	AddDeviceFn             func(context.Context, string) error
	AddFavoritePlaceFn      func(context.Context, string, models.Point) (*models.Location, error)
	AddFavoriteVehicleFn    func(context.Context, string, *string) error
	CreateUserFn            func(context.Context, *models.User) error
	FavoritePlacesFn        func(context.Context) ([]*models.Location, error)
	FavoriteVehiclesFn      func(context.Context) ([]string, error)
	FindAllFn               func(context.Context, *models.UserFilter) (*models.UserList, error)
	FindByEmailFn           func(context.Context, string) (*models.User, error)
	FindByIDFn              func(context.Context, string) (*models.User, error)
	GetUserDevicesFn        func(context.Context, models.UserFilter) ([]string, error)
	LastNAddressFn          func(context.Context, int) ([]*models.Location, error)
	LoginFn                 func(context.Context, string, string, ...string) (*models.User, error)
	MeFn                    func(context.Context) (*models.User, error)
	SetAvailabilityFn       func(context.Context, bool) error
	UpdateFn                func(context.Context, *models.User) error
	UpdatePlaceFn           func(context.Context, *models.UpdatePlace) (*models.Location, error)
	UpdateProfileFn         func(context.Context, *models.UpdateProfile) error
	TokenFn                 func(context.Context, *models.User) (string, error)
	LogoutFn                func(context.Context) error
	AddVehicleFn            func(context.Context, *models.Vehicle) error
	UpdateVehicleFn         func(context.Context, *models.Vehicle) error
	DeleteVehicleFn         func(context.Context, string) error
	AddDeviceTokenFn        func(context.Context, string, string) error
	DeleteFavoriteVehicleFn func(context.Context, string) error
	DeviceTokensFn          func(context.Context) ([]string, error)
	FavoritePlaceFn         func(context.Context, string) (*models.Location, error)
	RemoveDeviceTokenFn     func(context.Context, string) error
	SetActiveVehicleFn      func(context.Context, string) error
	SetPreferedCurrencyFn   func(context.Context, string) error
	VehicleFn               func(context.Context, string) (*models.Vehicle, error)
	VehiclesFn              func(context.Context) ([]*models.Vehicle, error)
	DeleteFavoritePlaceFn   func(context.Context, string) error
	UpdateFavoritePlaceFn   func(context.Context, string, models.UpdatePlace) error
}

// AddDeviceToken implements models.UserService.
func (s *UserService) AddDeviceToken(ctx context.Context, token string, name string) error {
	return s.AddDeviceTokenFn(ctx, token, name)
}

// DeleteFavoriteVehicle implements models.UserService.
func (s *UserService) DeleteFavoriteVehicle(ctx context.Context, v string) error {
	return s.DeleteFavoriteVehicleFn(ctx, v)
}

// DeviceTokens implements models.UserService.
func (s *UserService) DeviceTokens(ctx context.Context) ([]string, error) {
	return s.DeviceTokensFn(ctx)
}

// FavoritePlace implements models.UserService.
func (s *UserService) FavoritePlace(ctx context.Context, place string) (*models.Location, error) {
	return s.FavoritePlaceFn(ctx, place)
}

// RemoveDeviceToken implements models.UserService.
func (s *UserService) RemoveDeviceToken(ctx context.Context, token string) error {
	return s.RemoveDeviceTokenFn(ctx, token)
}

// SetActiveVehicle implements models.UserService.
func (s *UserService) SetActiveVehicle(ctx context.Context, plate string) error {
	return s.SetActiveVehicleFn(ctx, plate)
}

// SetPreferedCurrency implements models.UserService.
func (s *UserService) SetPreferedCurrency(ctx context.Context, cur string) error {
	return s.SetPreferedCurrencyFn(ctx, cur)
}

// Vehicle implements models.UserService.
func (s *UserService) Vehicle(ctx context.Context, palte string) (*models.Vehicle, error) {
	return s.VehicleFn(ctx, palte)
}

// Vehicles implements models.UserService.
func (s *UserService) Vehicles(ctx context.Context) ([]*models.Vehicle, error) {
	return s.VehiclesFn(ctx)
}

// DeleteFavoritePlace implements models.UserService.
func (s *UserService) DeleteFavoritePlace(ctx context.Context, place string) error {
	return s.DeleteFavoritePlaceFn(ctx, place)
}

// UpdateFavoritePlace implements models.UserService.
func (s *UserService) UpdateFavoritePlace(ctx context.Context, palce string, update models.UpdatePlace) error {
	return s.UpdateFavoritePlaceFn(ctx, palce, update)
}

// DeleteVehicle implements models.UserService.
func (s *UserService) DeleteVehicle(ctx context.Context, id string) error {
	return s.DeleteVehicleFn(ctx, id)
}

// UpdateVehicle implements models.UserService.
func (s *UserService) UpdateVehicle(ctx context.Context, v *models.Vehicle) error {
	return s.UpdateVehicleFn(ctx, v)
}

// AddVehicle implements models.UserService.
func (s *UserService) AddVehicle(ctx context.Context, v *models.Vehicle) error {
	return s.AddVehicleFn(ctx, v)
}

// Logout implements models.UserService.
func (s *UserService) Logout(ctx context.Context) error {
	return s.LogoutFn(ctx)
}

// Token implements models.UserService.
func (s *UserService) Token(ctx context.Context, user *models.User) (string, error) {
	return s.TokenFn(ctx, user)
}

// AddDevice implements models.UserService.
func (s *UserService) AddDevice(ctx context.Context, device string) error {
	return s.AddDeviceFn(ctx, device)
}

// AddFavoritePlace implements models.UserService.
func (s *UserService) AddFavoritePlace(ctx context.Context, name string, point models.Point) (*models.Location, error) {
	return s.AddFavoritePlaceFn(ctx, name, point)
}

// AddFavoriteVehicle implements models.UserService.
func (s *UserService) AddFavoriteVehicle(ctx context.Context, vehicle string, name *string) error {
	return s.AddFavoriteVehicleFn(ctx, vehicle, name)
}

// CreateUser implements models.UserService.
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.CreateUserFn(ctx, user)
}

// FavoritePlaces implements models.UserService.
func (s *UserService) FavoritePlaces(ctx context.Context) ([]*models.Location, error) {
	return s.FavoritePlacesFn(ctx)
}

// FavoriteVehicles implements models.UserService.
func (s *UserService) FavoriteVehicles(ctx context.Context) ([]string, error) {
	return s.FavoriteVehiclesFn(ctx)
}

// FindAll implements models.UserService.
func (s *UserService) FindAll(ctx context.Context, filter *models.UserFilter) (*models.UserList, error) {
	return s.FindAllFn(ctx, filter)
}

// FindByEmail implements models.UserService.
func (s *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.FindByEmailFn(ctx, email)
}

// FindByID implements models.UserService.
func (s *UserService) FindByID(ctz context.Context, id string) (*models.User, error) {
	return s.FindByIDFn(ctz, id)
}

// GetUserDevices implements models.UserService.
func (s *UserService) GetUserDevices(ctx context.Context, filter models.UserFilter) ([]string, error) {
	return s.GetUserDevicesFn(ctx, filter)
}

// LastNAddress implements models.UserService.
func (s *UserService) LastNAddress(ctz context.Context, limit int) ([]*models.Location, error) {
	return s.LastNAddressFn(ctz, limit)
}

// Login implements models.UserService.
func (s *UserService) Login(ctx context.Context, email, otp string, role ...string) (*models.User, error) {
	return s.LoginFn(ctx, email, otp)
}

// Me implements models.UserService.
func (s *UserService) Me(ctx context.Context) (*models.User, error) {
	return s.MeFn(ctx)
}

// SetAvailability implements models.UserService.
func (s *UserService) SetAvailability(ctx context.Context, available bool) error {
	return s.SetAvailabilityFn(ctx, available)
}

// Update implements models.UserService.
func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.UpdateFn(ctx, user)
}

// UpdatePlace implements models.UserService.
func (s *UserService) UpdatePlace(ctx context.Context, place *models.UpdatePlace) (*models.Location, error) {
	return s.UpdatePlaceFn(ctx, place)
}

// UpdateProfile implements models.UserService.
func (s *UserService) UpdateProfile(ctx context.Context, Profile *models.UpdateProfile) error {
	return s.UpdateProfileFn(ctx, Profile)
}
