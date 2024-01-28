package identity

import (
	"context"
	"encoding/base64"
)

type ClientStatus int

const (
	ClientStatusActive ClientStatus = iota
	ClientStatusInactive
	ClientStatusDeleted
)

type ClientType int

const (
	ClientTypePublic ClientType = iota + 1
	ClientTypeConfidential
	ClientTypeHybrid
	ClientTypeDriver
	ClientTypeRider
)

type Client struct {
	ID          ID           `bson:"_id,omitempty"`
	Name        string       `bson:"name"`
	Domain      string       `bson:"domain"`
	Secret      string       `bson:"secret"`
	Status      ClientStatus `bson:"status"`
	Scopes      []Scope      `bson:"scopes"`
	Type        ClientType   `bson:"type"`
	PrivateKey  string       `bson:"private_key"`
	PublicKey   string       `bson:"public_key"`
	RedirectURI string       `bson:"redirect_uri,omitempty"`
	CreatedAt   int64        `bson:"created_at,omitempty"`
	UpdatedAt   int64        `bson:"updated_at,omitempty"`
	DeletedAt   int64        `bson:"deleted_at,omitempty"`
}

func (c *Client) Update(cli Client) error {
	if !c.ID.Equal(cli.ID) {
		return NewInvalidParameter("id", "id must match client id")
	}
	c.Name = cli.Name
	c.Domain = cli.Domain
	c.RedirectURI = cli.RedirectURI
	c.UpdatedAt = Now().UnixNano()
	return nil
}

type ClientFilter struct {
	ID     []ID     `bson:"_id,omitempty"`
	Name   []string `bson:"name,omitempty"`
	Domain []string `bson:"domain,omitempty"`
	Key    string   `bson:"private_key,omitempty"`
	Limit  int      `bson:"limit,omitempty"`
	Token  string   `bson:"token,omitempty"`
}

type ClientService interface {
	Create(context.Context, *Client) error
	Update(context.Context, *Client) error
	UpdateKey(context.Context, ID, bool) error
	DeleteByID(context.Context, ID) error
	FindByID(context.Context, ID) (*Client, error)
	FindByKey(context.Context, string) (*Client, error)
	FindClients(context.Context, ClientFilter) ([]*Client, string, error)
}

type AuthKey []byte

func (k AuthKey) MarshalText() ([]byte, error) {
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len(k)))
	base64.RawURLEncoding.Encode(dst, k)
	return dst, nil
}

func (k *AuthKey) UnmarshalText(b []byte) error {
	*k = make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	_, err := base64.RawURLEncoding.Decode(*k, b)
	return err
}

func (k *AuthKey) String() string {
	return base64.RawURLEncoding.EncodeToString(*k)
}
