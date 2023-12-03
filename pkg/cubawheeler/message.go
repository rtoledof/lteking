package cubawheeler

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

type Message struct {
	ID      string        `json:"id" bson:"_id"`
	Trip    string        `json:"trip" bson:"trip"`
	From    string        `json:"from" bson:"from"`
	To      string        `json:"to" bson:"to"`
	Message string        `json:"message" bson:"message"`
	Status  MessageStatus `json:"status" bson:"status"`
}

type MessageRequest struct {
	Trip    string
	From    string
	To      string
	Message string
	Status  MessageStatus
}

type MessageFilter struct {
	Limit   int
	Token   string
	Ids     []string
	Trip    string
	From    string
	To      string
	Message string
	Status  MessageStatus
}

type MessageService interface {
	Create(context.Context, *MessageRequest) (*Message, error)
	Update(context.Context, *MessageRequest) (*Message, error)
	FindByID(context.Context, string) (*Message, error)
	FindAll(context.Context, *MessageFilter) ([]*Message, string, error)
}

type MessageStatus string

const (
	MessageStatusNew     MessageStatus = "NEW"
	MessageStatusRead    MessageStatus = "READ"
	MessageStatusDeleted MessageStatus = "DELETED"
)

var AllMessageStatus = []MessageStatus{
	MessageStatusNew,
	MessageStatusRead,
	MessageStatusDeleted,
}

func (e MessageStatus) IsValid() bool {
	switch e {
	case MessageStatusNew, MessageStatusRead, MessageStatusDeleted:
		return true
	}
	return false
}

func (e MessageStatus) String() string {
	return string(e)
}

func (e *MessageStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageStatus", str)
	}
	return nil
}

func (e MessageStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
