package mongo

import (
	"context"
	"fmt"

	"github.com/go-oauth2/oauth2/v4"
	"go.mongodb.org/mongo-driver/bson"

	"cubawheeler.io/pkg/cannon"
	"cubawheeler.io/pkg/cubawheeler"
)

var _ cubawheeler.ApplicationService = &ApplicationService{}
var _ oauth2.ClientStore = &ApplicationService{}

const ApplicationCollection Collections = "applications"

type ApplicationService struct {
	db *DB
}

// GetByID implements oauth2.ClientStore.
func (s *ApplicationService) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return findApplicationByID(ctx, s.db, id)
}

func NewApplicationService(db *DB) *ApplicationService {
	return &ApplicationService{db: db}
}

func (s *ApplicationService) FindByClient(ctx context.Context, client string) (*cubawheeler.Application, error) {
	apps, _, err := findApplications(ctx, s.db, cubawheeler.ApplicationFilter{
		Client: &client,
		Limit:  1,
	})
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return apps[0], nil
}

func (s *ApplicationService) FindByID(ctx context.Context, input string) (*cubawheeler.Application, error) {
	apps, _, err := findApplications(ctx, s.db, cubawheeler.ApplicationFilter{
		Ids:   []*string{&input},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return apps[0], nil
}

// CreateApplication implements cubawheeler.ApplicationService.
func (s *ApplicationService) CreateApplication(ctx context.Context, input cubawheeler.ApplicationRequest) (*cubawheeler.Application, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	app := &cubawheeler.Application{
		Name:   input.Name,
		Type:   input.Type,
		Client: input.Client,
		Secret: input.Secret,
	}
	err := createApplication(ctx, s.db, app)
	return app, err
}

// FindApplications implements cubawheeler.ApplicationService.
func (s *ApplicationService) FindApplications(ctx context.Context, input *cubawheeler.ApplicationFilter) (*cubawheeler.ApplicationList, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	apps, token, err := findApplications(ctx, s.db, *input)
	if err != nil {
		return nil, err
	}
	return &cubawheeler.ApplicationList{
		Token: &token,
		Data:  apps,
	}, nil
}

// UpdateApplicationCredentials implements cubawheeler.ApplicationService.
func (s *ApplicationService) UpdateApplicationCredentials(ctx context.Context, application string) (*cubawheeler.Application, error) {
	usr := cubawheeler.UserFromContext(ctx)
	if usr == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if usr.Role != cubawheeler.RoleAdmin {
		return nil, cubawheeler.ErrAccessDenied
	}
	applications, _, err := findApplications(ctx, s.db, cubawheeler.ApplicationFilter{
		Ids:   []*string{&application},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	app := applications[0]
	private, public, err := cubawheeler.NewKeyPair()
	if err != nil {
		return nil, fmt.Errorf("unable to generate a new key pair: %v: %w", err, cubawheeler.ErrInternal)
	}
	app.Client = private.X.String()
	app.Secret = string(public)
	f := bson.D{
		{Key: "private", Value: app.Client},
		{Key: "public", Value: app.Secret},
	}
	return app, updateApplications(ctx, application, s.db, f)
}

func createApplication(ctx context.Context, db *DB, app *cubawheeler.Application) error {
	collection := db.client.Database(database).Collection(ApplicationCollection.String())
	app.ID = cubawheeler.NewID().String()
	if len(app.Client) == 0 {
		private, public, err := cubawheeler.NewKeyPair()
		if err != nil {
			return err
		}
		app.Client = private.X.String()
		app.Secret = string(public)
	}

	_, err := collection.InsertOne(ctx, app)
	if err != nil {
		return cubawheeler.ErrInternal
	}
	return nil
}

func findApplications(ctx context.Context, db *DB, filter cubawheeler.ApplicationFilter) (_ []*cubawheeler.Application, _ string, err error) {
	logger := cannon.LoggerFromContext(ctx)
	if logger != nil {
		logger.Info("find applications")
	}
	defer func() {
		if err != nil && logger != nil {
			logger.Info(fmt.Sprintf("error finding applications: %v", err))
		}
	}()
	var applications []*cubawheeler.Application
	var token string
	f := bson.D{}
	if filter.Name != nil {
		f = append(f, bson.E{Key: "name", Value: *filter.Name})
	}
	if filter.Client != nil {
		f = append(f, bson.E{Key: "client", Value: *filter.Client})
	}
	if filter.Token != nil {
		f = append(f, bson.E{Key: "_id", Value: bson.A{"$gt", filter.Token}})
	}
	collection := db.client.Database(database).Collection(ApplicationCollection.String())
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", fmt.Errorf("unable to find the applications: %v: %w", err, cubawheeler.ErrInternal)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var app cubawheeler.Application
		err := cur.Decode(&app)
		if err != nil {
			return nil, "", fmt.Errorf("unable to decode the application: %v: %w", err, cubawheeler.ErrInternal)
		}
		applications = append(applications, &app)
		if len(applications) > filter.Limit {
			token = applications[filter.Limit].ID
			applications = applications[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", fmt.Errorf("an error processing the applications: %v: %w", err, cubawheeler.ErrInternal)
	}
	return applications, token, nil
}

func updateApplications(ctx context.Context, app string, db *DB, f bson.D) error {
	collection := db.client.Database(database).Collection(ApplicationCollection.String())
	_, err := collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: app}}, f)
	if err != nil {
		return fmt.Errorf("unable to update the application: %s, %v: %w", app, err, cubawheeler.ErrInternal)
	}
	return nil
}

func findApplicationByID(ctx context.Context, db *DB, id string) (*cubawheeler.Application, error) {
	app, _, err := findApplications(ctx, db, cubawheeler.ApplicationFilter{Ids: []*string{&id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(app) == 0 {
		return nil, cubawheeler.ErrNotFound
	}
	return app[0], nil
}
