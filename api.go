package finn

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	client "github.com/marcosQuesada/finn/http"
)

const apVersion = "v1"
const path = "organisation/accounts"

// ErrVersionConflict happens on version conflict error
var ErrVersionConflict = errors.New("version conflict")

// httpClient defines http transport
type httpClient interface {
	Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error)
	CreateRequest(method, url string, body interface{}) (*http.Request, error)
}

// APIClient defines an account http api client
type APIClient struct {
	api httpClient
}

// NewAPIClient instantiates api client
func NewAPIClient(api httpClient) *APIClient {
	return &APIClient{
		api: api,
	}
}

// Create invokes account creation
func (c *APIClient) Create(ctx context.Context, account *Account) (*Account, error) {
	uri := fmt.Sprintf("%s/%s", apVersion, path)
	req, respErr := c.api.CreateRequest(http.MethodPost, uri, account)
	if respErr != nil {
		return nil, fmt.Errorf("unexpected error creating request, error %v", respErr)
	}

	acc := &Account{}
	resp, respErr := c.api.Do(ctx, req, acc)
	if respErr != nil {
		return nil, fmt.Errorf("unexpected error executing http request, error %v", respErr)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, client.ErrInternalServer
	}

	return acc, nil
}

// Fetch gets user account by uuid
func (c *APIClient) Fetch(ctx context.Context, uuid string) (*Account, error) {
	uri := fmt.Sprintf("%s/%s/%s", apVersion, path, uuid)
	req, err := c.api.CreateRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	a := &Account{}
	_, err = c.api.Do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// List accounts with pagination
func (c *APIClient) List(ctx context.Context, pags *Pagination) (*AccountList, error) {
	uri := fmt.Sprintf("%s/%s?%s", apVersion, path, pags.QueryString())

	req, err := c.api.CreateRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	a := &AccountList{}
	_, err = c.api.Do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Delete removes account by user uuid and version
func (c *APIClient) Delete(ctx context.Context, uuid string, version int) error {
	uri := fmt.Sprintf("%s/%s/%s?version=%d", apVersion, path, uuid, version)

	req, err := c.api.CreateRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	resp, respErr := c.api.Do(ctx, req, nil)
	if resp != nil && resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if resp != nil && resp.StatusCode == http.StatusConflict {
		return ErrVersionConflict
	}

	return respErr
}

// Pagination defines list subset with page offset and size
type Pagination struct {
	Page int
	Size int
}

// NewPagination instantiates pagination
func NewPagination(page, size int) *Pagination {
	return &Pagination{
		Page: page,
		Size: size,
	}
}

// QueryString translate pagination to query string
func (p *Pagination) QueryString() string {
	return fmt.Sprintf("page[number]=%d&page[size]=%d", p.Page, p.Size)
}
