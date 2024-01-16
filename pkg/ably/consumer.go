package ably

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/ably/ably-go/ably"
	amqp "github.com/rabbitmq/amqp091-go"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/realtime"
)

type Body struct {
	ID       string              `json:"id"`
	Time     uint64              `json:"timestamp"`
	Name     string              `json:"name"`
	User     string              `json:"ClientId"`
	Encoding string              `json:"encoding"`
	Data     string              `json:"data"`
	Action   ably.PresenceAction `json:"action"`
}

type Source string

const (
	ChannelMessage  Source = "channel.message"
	ChannelPrecense Source = "channel.presence"
)

type ConnectionStatus int

const (
	Joined ConnectionStatus = iota + 1
	Abandon
)

type Message struct {
	Source         Source `json:"source"`
	Channel        string `json:"channel"`
	MessagesEvents []Body `json:"messages,omitempty"`
	PresenceEvents []Body `json:"presence,omitempty"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Bearing   float64 `json:"bearing"`
}

type Consumer service

func (c Consumer) Consume(queue, consumer string, autoAsk, exclusive, noLocal, noWait bool, args amqp.Table) error {
	conn, err := c.client.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := c.client.Channel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Qos(5, 0, false)
	if err != nil {
		return fmt.Errorf("unable to set prefetch: %v: %w", err, cubawheeler.ErrInternal)
	}
	msgs, err := ch.Consume(queue, consumer, autoAsk, exclusive, noLocal, noWait, args)
	if err != nil {
		return fmt.Errorf("unable to consume from the queue: %s: %v: %w", queue, err, cubawheeler.ErrInternal)
	}
	go func() {
		for msg := range msgs {
			var data Message
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				slog.Info("unable to unmarshal data")
				msg.Ack(false)
			}

			fmt.Printf("Message: %s\n", msg.Body)
			switch data.Source {
			case ChannelMessage:
				for _, m := range data.MessagesEvents {
					var location Location
					if m.Encoding == "json" {
						if err := json.Unmarshal([]byte(m.Data), &location); err != nil {
							slog.Info(fmt.Sprintf("unable to unmarshal the message: %s", m.Data))
							continue
						}
						realtime.DriverLocations <- cubawheeler.Location{
							User: m.User,
							Geolocation: cubawheeler.GeoLocation{
								Type:        cubawheeler.ShapeTypePoint,
								Lat:         location.Latitude,
								Long:        location.Longitude,
								Bearing:     location.Bearing,
								Coordinates: []float64{location.Longitude, location.Latitude},
							},
						}
					}
				}
			case ChannelPrecense:
				for _, v := range data.PresenceEvents {
					userStatus := realtime.UserStatus{
						User:      v.User,
						Available: v.Action == ably.PresenceActionEnter,
					}
					switch v.Action {
					case ably.PresenceActionEnter, ably.PresenceActionLeave:
						realtime.UserAvailabilityStatus <- userStatus
					}
				}
			}

			dotCount := bytes.Count(msg.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			msg.Ack(false)
		}
	}()
	<-c.client.done
	return nil
}
