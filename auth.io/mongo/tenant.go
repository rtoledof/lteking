package mongo

import (
	"auth.io/derrors"
	"auth.io/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var _ models.TenantService = &TenantService{}

const tenantCollection = "tenants"

type TenantService struct {
	*DB
}

func NewTenantService(db *DB) *TenantService {
	return &TenantService{db}
}

func (s *TenantService) Create(ctx context.Context, tenant models.Tenant) (_ *models.Tenant, err error) {
	defer derrors.Wrap(&err, "mongo.TenantService.Create")
	if err := createTenant(ctx, s.DB, &tenant); err != nil {
		return nil, models.NewInternalError(err)
	}
	return &tenant, nil
}

func (s *TenantService) FindAll(ctx context.Context, filter models.TenantFilter) (_ []models.Tenant, _ string, err error) {
	defer derrors.Wrap(&err, "mongo.TenantService.FindAll")
	return findAll(ctx, s.DB, filter)
}

func (s *TenantService) FindByID(ctx context.Context, id models.ID) (_ models.Tenant, err error) {
	defer derrors.Wrap(&err, "mongo.TenantService.FindByID")
	return findTenantByID(ctx, s.DB, id)
}

func (s *TenantService) Update(ctx context.Context, tenant models.Tenant) (_ models.Tenant, err error) {
	defer derrors.Wrap(&err, "mongo.TenantService.Update")
	if err := updateTenant(ctx, s.DB, &tenant); err != nil {
		return models.Tenant{}, models.NewInternalError(err)
	}
	return tenant, nil
}

func (s *TenantService) Delete(ctx context.Context, id models.ID) (err error) {
	defer derrors.Wrap(&err, "mongo.TenantService.Delete")
	if _, err := findTenantByID(ctx, s.DB, id); err != nil {
		return err
	}
	return deleteTenant(ctx, s.DB, id)
}

func findAll(ctx context.Context, db *DB, filter models.TenantFilter) ([]models.Tenant, string, error) {
	var tenants []models.Tenant
	var token string
	f := getFilter(filter)

	cursor, err := db.Collection(tenantCollection).Find(ctx, f)
	if err != nil {
		return nil, token, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var tenant models.Tenant
		if err := cursor.Decode(&tenant); err != nil {
			return nil, token, err
		}
		tenants = append(tenants, tenant)
		if filter.Limit > 0 && len(tenants) == filter.Limit+1 {
			token = tenant.ID.String()
			tenants = tenants[:filter.Limit]
			return tenants, token, nil
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, token, err
	}
	return tenants, token, nil
}

func findTenantByID(ctx context.Context, db *DB, id models.ID) (models.Tenant, error) {
	tenants, _, err := findAll(ctx, db, models.TenantFilter{
		ID:    []models.ID{id},
		Limit: 1,
	})
	if err != nil {
		return models.Tenant{}, err
	}
	if len(tenants) == 0 {
		return models.Tenant{}, models.NewNotFound("tenant")
	}
	return tenants[0], nil
}

func createTenant(ctx context.Context, db *DB, tenant *models.Tenant) error {
	if tenant.ID == nil {
		tenant.ID = models.NewID()
		tenant.CreatedAt = time.Now().UTC().Unix()
		tenant.Status = models.TenantStatusActive
	}
	_, err := db.Collection(tenantCollection).InsertOne(ctx, tenant)
	return err
}

func updateTenant(ctx context.Context, db *DB, tenant *models.Tenant) error {
	tenant.UpdatedAt = time.Now().UTC().Unix()
	_, err := db.Collection(tenantCollection).UpdateOne(ctx, bson.M{"_id": tenant.ID}, bson.M{"$set": tenant})
	return err
}

func deleteTenant(ctx context.Context, db *DB, id models.ID) error {
	_, err := db.Collection(tenantCollection).DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func getFilter(filter models.TenantFilter) bson.D {
	f := bson.D{}
	if len(filter.ID) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$in", Value: filter.ID}}})
	}
	if len(filter.Name) > 0 {
		f = append(f, bson.E{Key: "name", Value: bson.D{{Key: "$in", Value: filter.Name}}})
	}
	if len(filter.Domain) > 0 {
		f = append(f, bson.E{Key: "domain", Value: bson.D{{Key: "$in", Value: filter.Domain}}})
	}
	if len(filter.Status) > 0 {
		f = append(f, bson.E{Key: "status", Value: bson.D{{Key: "$in", Value: filter.Status}}})
	}
	if len(filter.Clients) > 0 {
		f = append(f, bson.E{Key: "clients", Value: bson.D{{Key: "$in", Value: filter.Clients}}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{Key: "_id", Value: bson.D{{Key: "$gte", Value: filter.Token}}})
	}

	return f
}
