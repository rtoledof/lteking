package ably

import (
	"context"
	"encoding/json"

	"cubawheeler.io/pkg/realtime"
	"github.com/ably/ably-go/ably"
)

type PushType string

const (
	NewTripRequest        PushType = "NEW_TRIP_REQUEST"
	TripRequestCancelled  PushType = "TRIP_REQUEST_CANCELLED"
	NewChatMessage        PushType = "NEW_CHAT_MESSAGE"
	UserDocumentsReviewed PushType = "USER_DOCUMENTS_REVIEWED"
)

type Notifier service

func (n Notifier) NotifyToDevices(ctx context.Context, devices []string, order realtime.OrderNotification, client *ably.Realtime, ablyRest *ably.REST) error {
	for _, deviceID := range devices {
		body := map[string]interface{}{
			"recipient": map[string]any{
				"deviceId": deviceID,
			},
			"notification": map[string]interface{}{
				"title": "New Trip",
				"body":  "You have a new Trip",
			},
			"data": map[string]interface{}{
				"order": order,
			},
		}
		data, _ := json.Marshal(body)
		if err := n.client.Push("/push/publish", data); err != nil {
			return err
		}
	}

	return nil
}

func (n Notifier) NotifyRiderOrderAccepted(ctx context.Context, deviceIds []string, order realtime.OrderNotification) error {
	for _, deviceId := range deviceIds {
		body := map[string]interface{}{
			"recipient": map[string]any{
				"deviceId": deviceId,
			},
			"notification": map[string]interface{}{
				"title": "Order Accepted",
				"body":  "Your order has been accepted by a driver",
			},
			"data": map[string]interface{}{
				"order": order,
			},
		}
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		if err := n.client.Push("/push/publish", data); err != nil {
			return err
		}
	}
	return nil
}
