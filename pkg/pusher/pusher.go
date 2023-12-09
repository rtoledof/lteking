package pusher

import (
	"cubawheeler.io/pkg/errors"
	"encoding/json"
	"fmt"
	"github.com/pusher/pusher-http-go/v5"

	"cubawheeler.io/pkg/cubawheeler"
)

var (
	prefixFilter = "presence-"
)

type Pusher struct {
	client pusher.Client
}

func NewPusher(app, key, secret, cluster string, secure bool) *Pusher {
	return &Pusher{
		client: pusher.Client{
			AppID:   app,
			Key:     key,
			Secret:  secret,
			Cluster: cluster,
			Secure:  secure,
		},
	}
}

func (p *Pusher) Authenticate(params []byte, usr *cubawheeler.User, trip *cubawheeler.Trip) (string, error) {

	memberData := pusher.MemberData{
		UserID: usr.ID,
		UserInfo: map[string]string{
			"name": usr.Name,
		},
	}
	var data struct {
		Auth string
	}
	// "presence-"+trip.ID
	presAuthData, err := p.client.AuthorizePresenceChannel(params, memberData)
	if err != nil {
		return "", fmt.Errorf("unable to authorize presence channel: %v: %w", err, errors.ErrInvalidInput)
	}
	if err := json.Unmarshal(presAuthData, &data); err != nil {
		return "", fmt.Errorf("unable to unmarshal presence auth data: %v: %w", err, errors.ErrInternal)
	}

	return data.Auth, nil
}
