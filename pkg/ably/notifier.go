package ably

import (
	"context"
	"fmt"
)

type Notifier service

func (n Notifier) NotifyToDevices(ctx context.Context, devices []string) error {
	for _, d := range devices {
		fmt.Println(d)
	}
	return nil
}
