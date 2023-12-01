package cubawheeler

import (
	"fmt"
	"io"
	"strconv"
)

type Message struct {
	ID      string        `json:"id"`
	Trip    string        `json:"trip"`
	From    string        `json:"from"`
	To      string        `json:"to"`
	Message string        `json:"message"`
	Status  MessageStatus `json:"status"`
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
