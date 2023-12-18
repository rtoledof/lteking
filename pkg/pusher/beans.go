package pusher

import (
	"cubawheeler.io/pkg/errors"
	"fmt"
	pn "github.com/pusher/push-notifications-go"
	"log/slog"
)

type Notification struct {
	Title    string
	Body     string
	Metadata map[string]string
}

func (n Notification) Request() map[string]any {
	notification := map[string]any{
		"title": n.Title,
		"body":  n.Body,
	}
	for k, v := range n.Metadata {
		notification[k] = v
	}
	return map[string]any{
		"apns": map[string]any{
			"aps": map[string]any{
				"alert": notification,
			},
		},
		"fcm": map[string]any{
			"notification": notification,
		},
		"web": map[string]any{
			"notification": notification,
		},
	}
}

type PushNotification struct {
	client pn.PushNotifications
}

func NewPushNotification(instanceId, secretKey string) *PushNotification {
	pusher, err := pn.New(instanceId, secretKey)
	if err != nil {
		panic(err)
	}
	return &PushNotification{client: pusher}
}

func (pn *PushNotification) PublishToInterest(interest []string, notification Notification) error {
	pubId, err := pn.client.PublishToInterests(interest, notification.Request())
	if err != nil {
		return fmt.Errorf("unable to send push notification: %v: %w", err, errors.ErrInternal)
	}
	slog.Info("push notification with publish id: %s was succesfull sent to interest: %s", pubId, interest)
	return nil
}

func (pn *PushNotification) PublishToUser(users []string, notification Notification) error {
	pubId, err := pn.client.PublishToUsers(users, notification.Request())
	if err != nil {
		return fmt.Errorf("unable to send push notification to users: %v: %w", err, errors.ErrInternal)
	}
	slog.Info(fmt.Sprintf("push notification with publish id: %s was successfull sent to users: %v", pubId, users))
	return nil
}

func (pn *PushNotification) GenerateToken(userID string) (map[string]any, error) {
	beansToken, err := pn.client.GenerateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to generate beansToken: %v: %w", err, errors.ErrInternal)
	}
	return beansToken, nil
}

func (pn *PushNotification) DeleteUser(user string) error {
	if err := pn.client.DeleteUser(user); err != nil {
		return fmt.Errorf("unable to delete the user: %v: %w", err, errors.ErrInternal)
	}
	return nil
}
