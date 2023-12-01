package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cubawheeler.io/pkg/cubawheeler"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// func UserCreate(ctx context.Context, input model.NewUser) (*cubawheeler.User, error) {
// 	user := cubawheeler.User{
// 		Name:     input.Name,
// 		Email:    strings.ToLower(input.Email),
// 		Password: []byte(input.Password),
// 	}

// 	if err := Db.Model(user).Create(&user).Error; err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

func (s *UserService) FindByID(ctx context.Context, id string) (*cubawheeler.User, error) {
	users, _, err := findAll(ctx, s.db, cubawheeler.UserFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}

	return users[0], nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*cubawheeler.User, error) {
	users, _, err := findAll(ctx, s.db, cubawheeler.UserFilter{Email: email, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}

	return users[0], nil
}

func (s *UserService) FindAll(ctx context.Context, filter cubawheeler.UserFilter) ([]*cubawheeler.User, string, error) {
	user := cubawheeler.UserForContext(ctx)
	if user == nil {
		return nil, "", errors.New("permission denied")
	}
	return findAll(ctx, s.db, filter)
}

func (s *UserService) UpdateOTP(ctx context.Context, email string, otp uint64) error {
	user, err := s.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	user.OTP = strconv.Itoa(int(otp))
	return s.db.UpdateColumn("otp", &user).Error
}

func (s *UserService) CreateUser(_ context.Context, user *cubawheeler.User) error {
	return s.db.Create(&user).Error
}

func findAll(ctx context.Context, db *gorm.DB, filter cubawheeler.UserFilter) ([]*cubawheeler.User, string, error) {
	var users []*cubawheeler.User
	var token string
	if len(filter.Ids) > 0 {
		db = db.Where("id IN ?", filter.Ids)
	}
	if len(filter.Email) > 0 {
		db = db.Where("email = ?", filter.Email)
	}
	if len(filter.Name) > 0 {
		db = db.Where("profile.name LIKE ?", fmt.Sprintf("%%s%", filter.Name))
	}
	if len(filter.Token) > 0 {
		db = db.Where("id >= ?", filter.Token)
	}
	if filter.Limit > 0 {
		db = db.Limit(int(filter.Limit) + 1)
	}
	if err := db.Find(&users).Error; err != nil {
		return nil, "", errors.New("unable to fetch the users")
	}
	if len(users) > int(filter.Limit) && filter.Limit > 0 {
		token = users[filter.Limit].ID
		users = users[:filter.Limit]
	}
	return users, token, db.Find(&users).Error
}
