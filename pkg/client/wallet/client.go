package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

var defaultURL = "http://wallet:80"

type Client struct {
	client  *http.Client
	BaseURL *url.URL

	common service

	Balance     *BalanceService
	Transaction *TransactionService
	Transfer    *TransferService
	Service     *WalletService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client, service string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	var err error

	baseURL, err := url.Parse(defaultURL)
	if err != nil {
		return nil, err
	}
	if service != "" {
		baseURL, err = url.Parse(service)
		if err != nil {
			return nil, err
		}
	}
	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}
	c.common.client = c

	c.Balance = (*BalanceService)(&c.common)
	c.Transaction = (*TransactionService)(&c.common)
	c.Transfer = (*TransferService)(&c.common)
	c.Service = (*WalletService)(&c.common)

	return c, nil
}

func (c *Client) NewBalanceRequest(currency string) (*http.Request, error) {
	u, err := c.BaseURL.Parse("/balance")
	if err != nil {
		return nil, err
	}
	if currency != "" {
		u.RawQuery = url.Values{"currency": []string{currency}}.Encode()
	}
	return c.newRequest(http.MethodGet, u, nil)
}

func (c *Client) NewRequest(method, path string, v url.Values) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	return c.newRequest(method, u, v)
}

func (c *Client) newRequest(method string, url *url.URL, value url.Values) (*http.Request, error) {
	var body io.Reader
	if value != nil {
		body = bytes.NewBufferString(value.Encode())
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
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
	defer derrors.WrapStack(&err, "client.wallet.CheckResponse")
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}
	type errorResponse struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		r.Body.Close()
		var errRsp = struct {
			Error errorResponse `json:"error"`
		}{}
		json.Unmarshal(data, &errRsp)
		return &cubawheeler.Error{StatusCode: r.StatusCode, Message: errRsp.Error.Message}
	}
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	return nil
}

type AuthTransport struct {
	Token string

	Transport http.RoundTripper
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
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
