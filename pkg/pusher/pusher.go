package pusher

import (
	"fmt"
	"log"
	"strconv"

	"github.com/pusher/pusher-http-go/v5"

	"cubawheeler.io/pkg/cubawheeler"
)

var (
	presencePrefixFilter = "presence"
	privatePrefixFilter  = "private"
	NewOrdersChannel     = make(chan *cubawheeler.Order, 1000)
	CanceledOrderChannel = make(chan *cubawheeler.Order, 1000)
)

type Pusher struct {
	client pusher.Client
	key    string
}

func NewPusher(app, key, secret, cluster string, secure bool) *Pusher {
	pusher := &Pusher{
		client: pusher.Client{
			AppID:   app,
			Key:     key,
			Secret:  secret,
			Cluster: cluster,
			Secure:  secure,
		},
		key: key,
	}

	go pusher.createOrderChannel()
	go pusher.cancelOrderChannel()

	return pusher
}

func (p *Pusher) Authenticate(params []byte, usr *cubawheeler.User, presence bool) ([]byte, error) {
	var err error
	var data []byte
	if presence {
		memberData := pusher.MemberData{
			UserID: usr.ID,
			UserInfo: map[string]string{
				"name": fmt.Sprintf("%s %s", usr.Profile.Name, usr.Profile.LastName),
			},
		}
		data, err = p.client.AuthorizePresenceChannel(params, memberData)
		if err != nil {
			return nil, fmt.Errorf("unable to authorize presence channel: %v: %w", err, cubawheeler.ErrInvalidInput)
		}
	} else {
		data, err = p.client.AuthorizePrivateChannel(params)

		if err != nil {
			return nil, fmt.Errorf("unable to authorize private channel: %v: %w", err, cubawheeler.ErrInvalidInput)
		}
	}

	return data, nil
}

func (p *Pusher) createOrderChannel() {
	// 1. Create the presence channel for the order
	for order := range NewOrdersChannel {
		params := map[string]string{
			"rider":       order.Rider,
			"pickup_lat":  strconv.Itoa(int(order.Items[0].PickUp.Lat)),
			"pickup_lon":  strconv.Itoa(int(order.Items[0].PickUp.Lon)),
			"dropoff_lat": strconv.Itoa(int(order.Items[0].DropOff.Lat)),
			"dropoff_lon": strconv.Itoa(int(order.Items[0].DropOff.Lon)),
			"price":       strconv.FormatInt(int64(order.Price), 10),
			"distance":    strconv.FormatInt(int64(order.Items[0].Meters), 10),
			"time":        strconv.FormatInt(int64(order.Items[0].Seconds), 10),
		}
		channelName := fmt.Sprintf("%s-%s", presencePrefixFilter, order.ID)
		if err := p.client.Trigger(channelName, cubawheeler.ChannelEventNewOrder.String(), params); err != nil {
			log.Printf("unable to create channel for order: %s", order.ID)
		}
	}
	// 2. Send new order event to nearby drivers - (Push notification)
}

func (p *Pusher) cancelOrderChannel() {
	// it(read from channel)
	for order := range CanceledOrderChannel {
		params := map[string]string{
			"status": order.Status.String(),
		}
		channelName := fmt.Sprintf("%s-%s", presencePrefixFilter, order.ID)
		if err := p.client.Trigger(channelName, cubawheeler.ChannelEventUpdateStatus.String(), params); err != nil {
			log.Printf("unable to create channel for order: %s", order.ID)
		}
	}
	p.client.Webhook(nil, nil)
}

func (p *Pusher) getDriverChannels() {
	//	prefixFilter := "presence-"
	//	attributes := "user_count"
	//	params := pusher.ChannelsParams{FilterByPrefix: &prefixFilter, Info: &attributes}
	//	channels, err := client.Channels(params)
}
