package ably

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
	"github.com/ably/ably-go/ably"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connectionString = "amqps://z8TS2w.nyBXtw:7QTI5Uq-wOaCLoazlWNigoh9LEbGIaNdIx4nRb2ZWKM@us-east-1-a-queue.ably.io:5671/shared"
	baseURL          = "https://rest.ably.io"
)

type Client struct {
	client           *http.Client
	connectionString string
	BaseURL          *url.URL
	done             chan struct{}
	rest             *ably.REST
	common           service
	realTime         *ably.Realtime
	ablyKey          string

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

func NewClient(
	connection string,
	done chan struct{},
	ablyKey string,
	ablyRealTime *ably.Realtime,
	cli *http.Client,
) *Client {
	if cli == nil {
		cli = http.DefaultClient
	}
	if len(connection) == 0 {
		connection = connectionString
	}

	client, err := ably.NewREST(ably.WithKey(ablyKey))
	if err != nil {
		panic(err)
	}

	c := &Client{
		client:           cli,
		connectionString: connection,
		done:             done,
		rest:             client,
		realTime:         ablyRealTime,
		ablyKey:          ablyKey,
	}
	c.common.client = c
	c.Consumer = (*Consumer)(&c.common)
	c.Notifier = (*Notifier)(&c.common)
	return c
}

func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := CheckResponse(resp); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err := json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil
			}
		}
	}
	return resp, err
}

func CheckResponse(r *http.Response) (err error) {
	defer derrors.WrapStack(&err, "ably.CheckResponse")
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}
	return cubawheeler.NewError(cubawheeler.ErrInsufficientFunds, r.StatusCode, "unable to process the request")
}

func (c *Client) Push(path string, data []byte) error {
	req, err := http.NewRequest(http.MethodPost, baseURL+path, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("unable to create request: %v: %w", err, cubawheeler.ErrInternal)
	}
	_, err = c.Do(req, nil)
	if err != nil {
		return fmt.Errorf("unable to make request: %v: %w", err, cubawheeler.ErrInternal)
	}

	return nil
}

type AuthTransport struct {
	Token string

	Transport http.RoundTripper
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Token, "")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(t.Token))))
	return t.transport().RoundTrip(req)
}

func (t *AuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *AuthTransport) transport() http.RoundTripper {
	if t.Transport == nil {
		return http.DefaultTransport
	}
	return t.Transport
}
