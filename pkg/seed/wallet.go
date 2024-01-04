package seed

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/mongo"
)

var _ Seeder = (*Wallet)(nil)

type Wallet struct {
	service cubawheeler.WalletService
	user    cubawheeler.UserService
}

func NewWallet(db *mongo.DB) *Wallet {
	return &Wallet{
		service: mongo.NewWalletService(db),
		user:    mongo.NewUserService(db, nil),
	}
}

// Down implements Seeder.
func (s *Wallet) Down() error {
	panic("unimplemented")
}

// Up implements Seeder.
func (s *Wallet) Up() error {
	usr := cubawheeler.User{
		Role: cubawheeler.RoleAdmin,
	}
	ctx := cubawheeler.NewContextWithUser(context.TODO(), &usr)
	users, err := s.user.FindAll(ctx, &cubawheeler.UserFilter{})
	if err != nil {
		return err
	}
	for _, v := range users.Data {
		if _, err := s.service.FindByOwner(ctx, v.ID); err != nil {
			if _, err := s.service.Create(ctx, v.ID); err != nil {
				return err
			}
		}
	}
	return nil
}
