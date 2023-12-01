package cubawheeler

import (
	"context"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID        string  `json:"id" gorm:"primaryKey;varchar(36);not null"`
	Name      string  `json:"name"`
	URL       *string `json:"url,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Whatsapp  *string `json:"whatsapp,omitempty"`
	Telegram  *string `json:"telegram,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Ads       []*Ads  `json:"ads,omitempty" gorm:"foreignKey:Owner"`
}

func (c *Client) BeforeSave(*gorm.DB) error {
	if c.ID == "" {
		c.ID = NewID().String()
	}
	return nil
}

type ClientStore interface {
	Store(context.Context, *Client) error
	Update(context.Context, *Client) error
	FindById(context.Context, string) (*Client, error)
	FindAll(context.Context) ([]Client, error)
}
