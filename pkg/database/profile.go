package database

import (
	"context"

	"gorm.io/gorm"

	"cubawheeler.io/pkg/cubawheeler"
)

type ProfileService struct {
	db *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

func (s *ProfileService) FindByUser(_ context.Context, user string) (*cubawheeler.Profile, error) {
	var profile cubawheeler.Profile
	profile.UserID = user
	return &profile, s.db.First(&profile).Error
}
