package mongo

import (
	"context"
	"cubawheeler.io/pkg/cubawheeler"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ cubawheeler.MessageService = &MessageService{}

type MessageService struct {
	db         *DB
	collection *mongo.Collection
}

func NewMessageService(db *DB) *MessageService {
	return &MessageService{
		db:         db,
		collection: db.client.Database(database).Collection("messages"),
	}
}

func (s *MessageService) Create(ctx context.Context, request *cubawheeler.MessageRequest) (*cubawheeler.Message, error) {
	message := &cubawheeler.Message{
		ID:      cubawheeler.NewID().String(),
		Trip:    request.Trip,
		From:    request.From,
		To:      request.To,
		Message: request.Message,
		Status:  request.Status,
	}
	if _, err := s.collection.InsertOne(ctx, message); err != nil {
		return nil, fmt.Errorf("unable to store the message: %w", err)
	}
	return message, nil
}

func (s *MessageService) Update(ctx context.Context, request *cubawheeler.MessageRequest) (*cubawheeler.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MessageService) FindByID(ctx context.Context, id string) (*cubawheeler.Message, error) {
	messages, _, err := findAllMessages(ctx, s.collection, &cubawheeler.MessageFilter{Ids: []string{id}, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, errors.New("message not found ")
	}
	return messages[0], nil
}

func (s *MessageService) FindAll(ctx context.Context, filter *cubawheeler.MessageFilter) ([]*cubawheeler.Message, string, error) {
	return findAllMessages(ctx, s.collection, filter)
}

func findAllMessages(ctx context.Context, collection *mongo.Collection, filter *cubawheeler.MessageFilter) ([]*cubawheeler.Message, string, error) {
	var messages []*cubawheeler.Message
	var token string
	f := bson.D{}
	if len(filter.Ids) > 0 {
		f = append(f, primitive.E{"_id", primitive.D{{"$in", filter.Ids}}})
	}
	if len(filter.Token) > 0 {
		f = append(f, bson.E{"token", primitive.E{"$gt", filter.Token}})
	}
	if len(filter.Trip) > 0 {
		f = append(f, bson.E{"trip", filter.Trip})
	}
	if len(filter.From) > 0 {
		f = append(f, bson.E{"from", filter.From})
	}
	if len(filter.To) > 0 {
		f = append(f, bson.E{"to", filter.To})
	}
	cur, err := collection.Find(ctx, f)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var message cubawheeler.Message
		err := cur.Decode(&message)
		if err != nil {
			return nil, "", err
		}
		messages = append(messages, &message)
		if len(messages) == filter.Limit+1 {
			token = messages[filter.Limit].ID
			messages = messages[:filter.Limit]
			break
		}
	}
	if err := cur.Err(); err != nil {
		return nil, "", err
	}
	return messages, token, nil
}
