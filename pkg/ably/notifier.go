package ably

import (
	"context"
	"os"

	"github.com/ably/ably-go/ably"
)

type Notifier service

func (n Notifier) NotifyToDevices(ctx context.Context, devices []string, order string) error {
	client, err := ably.NewRealtime(ably.WithKey(os.Getenv("ABLY_API_KEY")))
	if err != nil {
		return err
	}

	for _, deviceID := range devices {
		channel := client.Channels.Get("pushenabled:" + deviceID)

		err := channel.Publish(ctx, "New Trip", order)
		if err != nil {
			return err
		}
	}

	return nil
}
