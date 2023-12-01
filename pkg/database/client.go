package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"cubawheeler.io/pkg/cubawheeler"
)

type ClientStore struct {
	db *gorm.DB
}

func NewClientStore(db *gorm.DB) *ClientStore {
	return &ClientStore{db: db}
}

func (s *ClientStore) Store(_ context.Context, client *cubawheeler.Client) error {
	return s.db.Create(client).Error
}

func (s *ClientStore) Update(_ context.Context, client *cubawheeler.Client) error {
	return s.db.Updates(client).Error
}

func (s *ClientStore) FindById(_ context.Context, strID string) (*cubawheeler.Client, error) {
	var c = cubawheeler.Client{ID: strID}
	if err := s.db.Find(&c).Error; err != nil {
		return nil, fmt.Errorf("not found client: %w", err)
	}
	return &c, nil
}

func (s *ClientStore) FindAll(_ context.Context) ([]cubawheeler.Client, error) {
	var clients []cubawheeler.Client
	if err := s.db.Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("unable to find the clients: %w", err)
	}
	return clients, nil
}
