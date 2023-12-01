package bolt

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.UserService = &UserService{}

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindByID(context.Context, string) (*User, error) {
	panic("implement me")
}

func (s *UserService) FindByEmail(context.Context, string) (*User, error) {
	panic("implement me")
}

func (s *UserService) FindAll(context.Context, UserFilter) ([]*User, string, error) {
	panic("implement me")
}

func (s *UserService) UpdateOTP(context.Context, string, uint64) error {
	panic("implement me")
}

func (s *UserService) CreateUser(ctx context.Context, u *User) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return fmt.Errorf("create bucket %s: %w", userBucket, err)
		}
		eb, err := tx.CreateBucketIfNotExists(userBucketByEmail)
		if err != nil {
			return fmt.Errorf("create bucket %s: %w", userBucketByEmail, err)
		}
		u.ID = cubawheeler.NewID().String()
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(u); err != nil {
			return fmt.Errorf("encoding gob: %w", err)
		}
		if err := b.Put([]byte(u.ID), buf.Bytes()); err != nil {
			return fmt.Errorf("storing in path: %w", err)
		}

		if err := be.Put([]byte(u.Email), []byte(u.ID)); err != nil {
			return fmt.Errorf("storing in path: %w", err)
		}
		return nil
	})
}
