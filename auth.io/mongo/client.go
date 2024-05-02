package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"auth.io/derrors"
	"auth.io/models"
)

var _ models.ClientService = (*ClientService)(nil)

var ClientCollection Collections = "clients"

type ClientService struct {
	db *DB
}

func NewClientService(db *DB) *ClientService {
	index := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}},
	}
	_, err := db.Collection(ClientCollection).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic("unable to create user index")
	}
	return &ClientService{db: db}
}

// Create implements models.ClientService.
func (s *ClientService) Create(ctx context.Context, client *models.Client) (err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.CreateClient")
	user := models.UserFromContext(ctx)
	if user == nil || user.Role != models.RoleAdmin {
		return models.NewError(nil, 401, "not authorized")
	}

	if client.ID == nil {
		client.ID = models.NewID()
	}
	client.CreatedAt = models.Now().UnixNano()
	client.Status = models.ClientStatusActive
	_, err = findClientByID(ctx, s.db, models.ClientFilter{ID: []models.ID{client.ID}})
	if err == nil {
		return err
	}
	return createClient(ctx, s.db, client)
}

// Update implements models.ClientService.
func (s *ClientService) Update(ctx context.Context, client *models.Client) (err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.UpdateClient")
	user := models.UserFromContext(ctx)
	if user == nil || user.Role != models.RoleAdmin {
		return models.NewError(nil, 401, "not authorized")
	}
	filter := models.ClientFilter{ID: []models.ID{client.ID}}
	c, err := findClientByID(ctx, s.db, filter)
	if err != nil {
		return models.NewNotFound(fmt.Sprintf("client %s not found", client.ID))
	}
	if err := c.Update(*client); err != nil {
		return err
	}
	return updateClient(ctx, s.db, c, false)
}

func (s *ClientService) UpdateKey(ctx context.Context, client models.ID, keys bool) (err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.UpdateClient")
	user := models.UserFromContext(ctx)
	if user == nil || user.Role != models.RoleAdmin {
		return models.NewError(nil, 401, "not authorized")
	}
	c, err := findClientByID(ctx, s.db, models.ClientFilter{ID: []models.ID{client}})
	if c == nil {
		return models.NewNotFound(fmt.Sprintf("client %s not found", client))
	}
	return updateClient(ctx, s.db, c, keys)
}

// DeleteByID implements models.ClientService.
func (s *ClientService) DeleteByID(ctx context.Context, client models.ID) (err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.DeleteByID")
	user := models.UserFromContext(ctx)
	if user == nil || user.Role != models.RoleAdmin {
		return models.NewError(nil, 401, "not authorized")
	}
	filter := models.ClientFilter{ID: []models.ID{client}}
	c, err := findClientByID(ctx, s.db, filter)
	if err != nil {
		return models.NewNotFound(fmt.Sprintf("client %s not found", client))
	}
	c.Status = models.ClientStatusDeleted
	c.DeletedAt = models.Now().UnixNano()
	return updateClient(ctx, s.db, c, false)
}

// FindClients implements models.ClientService.
func (s *ClientService) FindClients(ctx context.Context, filter models.ClientFilter) (_ []*models.Client, _ string, err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.FindClients")
	user := models.UserFromContext(ctx)
	if user == nil {
		return nil, "", models.NewError(nil, 401, "not authorized")
	}
	if !user.Role.IsValid() || user.Role != models.RoleAdmin {
		return nil, "", models.NewError(nil, 401, "not authorized")
	}
	return findClients(ctx, s.db, filter)
}

// FindByID implements models.ClientService.
func (s *ClientService) FindByID(ctx context.Context, id models.ID) (_ *models.Client, err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.FindClientByID")
	user := models.UserFromContext(ctx)
	if user == nil {
		return nil, models.NewError(nil, 401, "not authorized")
	}
	if !user.Role.IsValid() || user.Role != models.RoleAdmin {
		return nil, models.NewError(nil, 401, "not authorized")
	}
	filter := models.ClientFilter{ID: []models.ID{id}}
	return findClientByClientID(ctx, s.db, filter)
}

// FindByKey implements models.ClientService.
func (s *ClientService) FindByKey(ctx context.Context, strKey string) (_ *models.Client, err error) {
	defer derrors.Wrap(&err, "mongo.ClientService.FindByKey")
	return findClientByKey(ctx, s.db, strKey)
}

func createClient(ctx context.Context, db *DB, client *models.Client) error {
	if client.ID == nil {
		client.ID = models.NewID()
	}
	client.CreatedAt = models.Now().UnixNano()
	client.Status = models.ClientStatusActive
	var key, public models.AuthKey
	key, public = models.NewKeyPair()

	client.PrivateKey = key.String()
	client.PublicKey = public.String()

	_, err := db.Collection(ClientCollection).InsertOne(ctx, client)
	return err
}

func findClients(ctx context.Context, db *DB, filter models.ClientFilter) ([]*models.Client, string, error) {
	var clients []*models.Client
	var token string
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
	if filter.Key != "" {
		f = append(f, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "private_key", Value: filter.Key}},
				bson.D{{Key: "public_key", Value: filter.Key}},
			},
		})
	}
	if filter.Token != "" {
		f = append(f, bson.E{Key: "token", Value: filter.Token})
	}

	cur, err := db.Collection(ClientCollection).Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var client models.Client
		if err := cur.Decode(&client); err != nil {
			return nil, "", err
		}
		clients = append(clients, &client)
		if filter.Limit > 0 && len(clients) == filter.Limit+1 {
			token = client.ID.String()
			clients = clients[:len(clients)-1]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", models.NewError(err, 500, "internal error")
	}
	return clients, token, nil
}

func findClientByID(ctx context.Context, db *DB, filter models.ClientFilter) (*models.Client, error) {
	clients, _, err := findClients(ctx, db, filter)
	if err != nil {
		return nil, err
	}
	if len(clients) == 1 {
		return clients[0], nil
	}
	return nil, models.NewNotFound("client not found")
}

func findClientByClientID(ctx context.Context, db *DB, filter models.ClientFilter) (*models.Client, error) {
	clients, _, err := findClients(ctx, db, filter)
	if err != nil {
		return nil, err
	}
	if len(clients) == 1 {
		return clients[0], nil
	}
	return nil, models.NewNotFound("client not found")
}

func updateClient(ctx context.Context, db *DB, update *models.Client, key bool) error {
	if _, err := db.Collection(ClientCollection).UpdateOne(ctx, bson.M{"_id": update.ID}, bson.M{"$set": update}); err != nil {
		return models.NewError(err, 500, "internal error")
	}
	return nil
}

func findClientByKey(ctx context.Context, db *DB, key string) (*models.Client, error) {
	filter := models.ClientFilter{Key: key, Limit: 1}
	clients, _, err := findClients(ctx, db, filter)
	if err != nil {
		return nil, err
	}

	if len(clients) == 1 {
		return clients[0], nil
	}

	return nil, models.NewNotFound("client not found")
}
