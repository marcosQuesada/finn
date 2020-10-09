package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClient_NewRequestBuildsProperRequestPath(t *testing.T) {
	c := NewClient()
	path := "foobar"
	req, err := c.CreateRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatalf("unexpected error creating request, error %v", err)
	}

	if got, want := req.URL.String(), fmt.Sprintf("%s%s", defaultBaseURL, path); want != got {
		t.Errorf("request url does not match, expected %s got %s", want, got)
	}
}

func TestClient_NewRequestAddsMarshalledValueToBody(t *testing.T) {
	c := NewClient()
	path := "foobar"
	v := &fakeAccount{
		ID:      "fakeID",
		Version: 1,
	}
	req, err := c.CreateRequest(http.MethodPost, path, v)
	if err != nil {
		t.Fatalf("unexpected error creating request, error %v", err)
	}

	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Unexpected request body content, error %v", err)
	}

	rawAccount, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unexpected error marshalling account, error %v", err)
	}

	if got, want := string(rawBody), string(rawAccount); want != got {
		t.Errorf("request body does not match, expected %s got %s", want, got)
	}
}

func TestNewClient_NewRequestAddsProperContentTypeOnEncodedBody(t *testing.T) {
	c := NewClient()
	path := "foobar"
	v := &fakeAccount{
		ID:      "fakeID",
		Version: 1,
	}
	req, err := c.CreateRequest(http.MethodPost, path, v)
	if err != nil {
		t.Fatalf("unexpected error creating request, error %v", err)
	}

	if got, want := req.Header.Get("Content-Type"), jsonContentType; want != got {
		t.Errorf("request content typr does not match, expected %s got %s", want, got)
	}
}

type fakeAccount struct {
	ID      string
	Version int
}

func TestClient_DoInvokesAndHydratesResponse(t *testing.T) {
	key := "foo"
	val := "bar"
	raw := []byte(fmt.Sprintf(`{"%s": "%s"}`, key, val))
	c := NewClient()
	c.client.Transport = &fakeTransport{statusCode: 200, body: raw}

	req, err := c.CreateRequest(http.MethodGet, "foo", nil)
	if err != nil {
		t.Fatalf("unexpected error creating request, error %v", err)
	}

	v := make(map[string]string)
	_, err = c.Do(context.Background(), req, &v)
	if err != nil {
		t.Fatalf("unexpected error executing request, error %v", err)
	}

	value, ok := v[key]
	if !ok {
		t.Fatal("foo key not found")
	}

	if value != val {
		t.Errorf("Unexpected response content, expected %s got %s", val, value)
	}
}
func TestClient_DoResponseStatusCodesWithErrorDetection(t *testing.T) {
	tests := []struct {
		statusCode int
		error      error
	}{
		{statusCode: 200, error: nil},
		{statusCode: 404, error: ErrContentNotFound},
		{statusCode: 500, error: ErrInternalServer},
	}

	c := NewClient()
	req, err := c.CreateRequest(http.MethodGet, "foo", nil)
	if err != nil {
		t.Fatalf("unexpected error creating request, error %v", err)
	}
	for _, test := range tests {
		c.client.Transport = &fakeTransport{statusCode: test.statusCode}
		_, err = c.Do(context.Background(), req, nil)
		if err != test.error {
			t.Fatalf("unexpected error executing request, error %v", err)
		}
	}
}

type fakeTransport struct {
	statusCode int
	body       []byte
}

func (f *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: f.statusCode,
		Body:       ioutil.NopCloser(bytes.NewBuffer(f.body)),
		Header:     make(http.Header),
	}, nil
}
