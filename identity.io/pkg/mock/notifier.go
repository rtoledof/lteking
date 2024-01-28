package mock

import (
	"context"

	"github.com/ably/ably-go/ably"

	"identity.io/pkg/realtime"
)

var _ realtime.Notifier = &Notifier{}

type Notifier struct {
	NotifyRiderOrderAcceptedFn func(context.Context, []string, realtime.OrderNotification) error
	NotifyToDevicesFn          func(context.Context, []string, realtime.OrderNotification, *ably.Realtime, *ably.REST) error
}

// NotifyRiderOrderAccepted implements realtime.Notifier.
func (s Notifier) NotifyRiderOrderAccepted(ctx context.Context, devices []string, notification realtime.OrderNotification) error {
	return s.NotifyRiderOrderAcceptedFn(ctx, devices, notification)
}

// NotifyToDevices implements realtime.Notifier.
func (s Notifier) NotifyToDevices(ctx context.Context, devices []string, notification realtime.OrderNotification, realTIme *ably.Realtime, rest *ably.REST) error {
	return s.NotifyToDevicesFn(ctx, devices, notification, realTIme, rest)
}
