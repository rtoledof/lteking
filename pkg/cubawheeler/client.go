package cubawheeler

import (
	"context"
)

type Client struct {
	ID        string  `json:"id" bson:"_id"`
	Name      string  `json:"name" bson:"name,omitempty"`
	URL       *string `json:"url,omitempty" bson:"url,omitempty"`
	Facebook  *string `json:"facebook,omitempty" bson:"facebook,omitempty"`
	Whatsapp  *string `json:"whatsapp,omitempty" bson:"whatsapp,omitempty"`
	Telegram  *string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Instagram *string `json:"instagram,omitempty" bson:"instagram,omitempty"`
	Ads       []*Ads  `json:"ads,omitempty" bson:"ads,omitempty"`
}

type ClientRequest struct {
	Ids       []string
	Name      string
	URL       *string
	Facebook  *string
	Whatsapp  *string
	Telegram  *string
	Instagram *string
}

type ClientService interface {
	Create(context.Context, *ClientRequest) error
	Update(context.Context, *ClientRequest) error
	FindById(context.Context, string) (*Client, error)
	FindAll(context.Context, *ClientRequest) ([]Client, string, error)
}
