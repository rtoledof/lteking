package models

import (
	"context"
)

type TenantStatus int64

const (
	TenantStatusActive TenantStatus = iota + 1
	TenantStatusInactive
	TenantStatusSuspended
	TenantStatusDeleted
)

type Tenant struct {
	ID        ID           `bson:"_id"`
	Name      string       `bson:"name"`
	Domain    string       `bson:"domain"`
	Status    TenantStatus `bson:"status"`
	CreatedAt int64        `bson:"created,omitempty"`
	UpdatedAt int64        `bson:"updated,omitempty"`
	DeletedAt int64        `bson:"deleted,omitempty"`
}

type TenantFilter struct {
	ID      []ID           `bson:"_id,omitempty"`
	Name    []string       `bson:"name,omitempty"`
	Domain  []string       `bson:"domain,omitempty"`
	Status  []TenantStatus `bson:"status,omitempty"`
	Clients []ID           `bson:"clients,omitempty"`
	Limit   int            `bson:"limit,omitempty"`
	Token   string         `bson:"token,omitempty"`
}

type TenantService interface {
	Create(context.Context, Tenant) (*Tenant, error)
	FindAll(context.Context, TenantFilter) ([]Tenant, string, error)
	FindByID(context.Context, ID) (Tenant, error)
	Update(context.Context, Tenant) (Tenant, error)
	Delete(context.Context, ID) error
}
