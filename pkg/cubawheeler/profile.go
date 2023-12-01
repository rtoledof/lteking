package cubawheeler

import (
	"context"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	ID       string `json:"id" gorm:"type:varchar(36);primaryKey" faker:"-"`
	Name     string `json:"name" faker:"name"`
	LastName string `json:"last_name" faker:"last_name"`
	UserID   string `faker:"-"`
	Gender   Gender `json:"gender"`
	Phone    string `json:"phone" faker:"phone_number"`
	Photo    string `json:"photo" faker:"url"`
	Licence  string `json:"licence,omitempty"`
	Dni      string `json:"dni,omitempty"`
}

func (u *Profile) BeforeSave(*gorm.DB) error {
	if u.ID == "" {
		u.ID = NewID().String()
	}
	return nil
}

type ProfileService interface {
	FindByUser(context.Context, string) (*Profile, error)
}
