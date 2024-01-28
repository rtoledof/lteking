package mapbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"order.io/pkg/derrors"
	"order.io/pkg/order"
)

const (
	defaultBaseURL   = "https://api.mapbox.com"
	defaultMediaType = "application/json"
)

type Client struct {
	AccessToken string
	client      *http.Client
	BaseURL     *url.URL

	common service

	Directions *DirectionService
}

type service struct {
	client *Client
}

func NewClient(accessToken string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		AccessToken: accessToken,
		BaseURL:     baseURL,
		client:      http.DefaultClient,
	}

	c.common.client = c
	c.Directions = (*DirectionService)(&c.common)

	return c
}

func (c *Client) NewRequest(method, urlStr string, body any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return c.newRequest(method, u, body)
}

func (c *Client) newRequest(method string, u *url.URL, body any) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}
	uid := order.NewID().String()
	req.Header.Add("request-id", uid)
	// if body != nil {
	// 	req.Header.Add("Content-Type", defaultMediaType)
	// }
	return req, nil
}

func (c *Client) Do(req *http.Request, v any) (*http.Response, string, error) {
	fmt.Println(req.URL.String())
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	if err := CheckResponse(resp); err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	var bodyResponse string
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			data, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(data, v)
			if err == io.EOF {
				err = nil
			}
			bodyResponse = string(data)
		}
	}
	return resp, bodyResponse, err
}

func CheckResponse(r *http.Response) (err error) {
	defer derrors.WrapStack(&err, "mapbox.CheckResponse")
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}
	type errorResponse struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	data, err := io.ReadAll(r.Body)
	fmt.Printf("%s\n", data)
	if err == nil && data != nil {
		r.Body.Close()
		var errRsp = struct {
			Error errorResponse `json:"error"`
		}{}
		json.Unmarshal(data, &errRsp)
		return &order.Error{Message: errRsp.Error.Message}
	}
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	return nil
}
