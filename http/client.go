package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const defaultBaseURL = "http://accountapi:8080/"
const jsonContentType = "application/vnd.api+json"

// ErrContentNotFound happens on 404 response status code
var ErrContentNotFound = errors.New("content not found")

// ErrNotAuthorized happens on forbidden access
var ErrNotAuthorized = errors.New("not authorized")

// ErrBadRequest happens on wrong request arguments
var ErrBadRequest = errors.New("bad request")

// ErrInternalServer happens on not controlled error
var ErrInternalServer = errors.New("internal server error")

// Client takes care on the whole http execution
type Client struct {
	client  *http.Client
	baseURL *url.URL
}

// NewClient creates an http client that points to default base url
func NewClient() *Client {
	baseUrl, _ := url.Parse(defaultBaseURL)

	return &Client{
		client:  &http.Client{},
		baseURL: baseUrl,
	}
}

// NewClientWithUrl creates an http client with specific url
func NewClientWithUrl(u *url.URL) *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: u,
	}
}

// Do executes an http.Request, when v is provided response body gets json unmarshalled
// response status code is validated against basic rules
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = c.validateStatusCode(resp.StatusCode)
	if err != nil {
		return resp, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if v != nil {
		unMarshallErr := json.NewDecoder(resp.Body).Decode(v)
		if unMarshallErr != nil {
			return resp, unMarshallErr
		}
	}

	return resp, nil
}

// CreateRequest creates an http API request, applies json encoding to body
func (c *Client) CreateRequest(method, url string, body interface{}) (*http.Request, error) {
	uri, err := c.baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf []byte
	if body != nil {
		var marshallErr error
		buf, marshallErr = json.Marshal(body)
		if marshallErr != nil {
			return nil, marshallErr
		}
	}
	req, reqErr := http.NewRequest(method, uri.String(), bytes.NewReader(buf))
	if reqErr != nil {
		return nil, reqErr
	}

	if body != nil {
		req.Header.Set("Content-Type", jsonContentType)
	}

	return req, nil
}

// validateStatusCode validates common status codes
func (c *Client) validateStatusCode(statusCode int) error {
	if http.StatusOK <= statusCode && statusCode < http.StatusMultipleChoices {
		return nil
	}

	if statusCode == http.StatusNotFound {
		return ErrContentNotFound
	}

	if statusCode == http.StatusBadRequest {
		return ErrBadRequest
	}

	if statusCode == http.StatusForbidden {
		return ErrNotAuthorized
	}

	return ErrInternalServer
}
