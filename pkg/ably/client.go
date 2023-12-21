package ably

import (
	"fmt"

	"cubawheeler.io/pkg/cubawheeler"
	"github.com/ably/ably-go/ably"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connectionString = "amqps://z8TS2w.nyBXtw:7QTI5Uq-wOaCLoazlWNigoh9LEbGIaNdIx4nRb2ZWKM@us-east-1-a-queue.ably.io:5671/shared"
)

type Client struct {
	connectionString string
	done             chan struct{}
	rest             *ably.REST
	common           service

	Consumer *Consumer
	Notifier *Notifier
}

type service struct {
	client *Client
}

func (c *Client) Dial() (*amqp.Connection, error) {
	conn, err := amqp.Dial(c.connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connecto to the queue: %v: %w", err, cubawheeler.ErrInternal)
	}
	return conn, nil
}

func (c *Client) Channel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("unable to connecto to the channel: %v: %w", err, cubawheeler.ErrInternal)
	}
	return ch, nil
}

func NewClient(connection string, done chan struct{}, abyKey string) *Client {
	if len(connection) == 0 {
		connection = connectionString
	}

	client, err := ably.NewREST(ably.WithKey(abyKey))
	if err != nil {
		panic(err)
	}

	c := &Client{
		connectionString: connection,
		done:             done,
		rest:             client,
	}
	c.common.client = c
	c.Consumer = (*Consumer)(&c.common)
	c.Notifier = (*Notifier)(&c.common)
	return c
}
