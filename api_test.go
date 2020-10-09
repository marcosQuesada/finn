package finn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	client "github.com/marcosQuesada/finn/http"
)

func TestCreateAccountReturnsNoErrorOnStatusAccepted(t *testing.T) {
	userID := uuid.New().String()
	acc := &Account{
		AccoundData: &AccoundData{
			Type:           "accounts",
			ID:             userID,
			OrganisationID: uuid.New().String(),
			Version:        0,
		}}
	rawAccount, err := json.Marshal(acc)
	if err != nil {
		t.Fatalf("unexpected error marshalling account, error %v", err)
	}
	h := &fakeHTTPClient{
		statusCode: http.StatusCreated,
		body:       rawAccount,
	}
	api := NewAPIClient(h)
	a, err := api.Create(context.Background(), acc)
	if err != nil {
		t.Errorf("unexpected error creating user, error %v", err)
	}

	if a == nil || a.AccoundData == nil {
		t.Fatal("created account response is nil")
	}

	if a.AccoundData.ID != userID {
		t.Errorf("Account userID does not match, expected %s got %s", userID, a.AccoundData.ID)
	}
}

func TestCreateAccountReturnsInternalServerErrorOnNonStatusAccepted(t *testing.T) {
	h := &fakeHTTPClient{}
	api := NewAPIClient(h)
	acc := &Account{
		AccoundData: &AccoundData{
			Type:           "accounts",
			ID:             uuid.New().String(),
			OrganisationID: uuid.New().String(),
			Version:        0,
		}}

	_, err := api.Create(context.Background(), acc)
	if err == nil {
		t.Error("expected error")
	}

	if !errors.Is(err, client.ErrInternalServer) {
		t.Errorf("Unexpected error type, got %v", err)
	}
}

func TestFetchAccountReturnsAFullPopulatedAccountOnValidStatusCode(t *testing.T) {
	userID := uuid.New().String()
	acc := &Account{
		AccoundData: &AccoundData{
			Type:           "accounts",
			ID:             userID,
			OrganisationID: uuid.New().String(),
			Version:        0,
		}}
	rawAccount, err := json.Marshal(acc)
	if err != nil {
		t.Fatalf("unexpected error marshalling account, error %v", err)
	}

	h := &fakeHTTPClient{
		body: rawAccount,
	}
	api := NewAPIClient(h)

	a, err := api.Fetch(context.Background(), userID)
	if err != nil {
		t.Fatalf("error fetching user, error %v", err)
	}

	if !reflect.DeepEqual(acc, a) {
		t.Error("accounts are not equal")
	}
}

func TestFetchAccountReturnsErrorOnNotFoundStatusCode(t *testing.T) {
	userID := uuid.New().String()
	h := &fakeHTTPClient{
		err: client.ErrContentNotFound,
	}
	api := NewAPIClient(h)

	_, err := api.Fetch(context.Background(), userID)
	if err == nil {
		t.Fatal("expected not found error")
	}

	if !errors.Is(err, client.ErrContentNotFound) {
		t.Errorf("error is not content not found got %v", err)
	}
}

func TestListAccountsReturnASliceOfAccountsOnStatusOk(t *testing.T) {
	accs := make([]*AccoundData, 10)
	for i := 0; i < 10; i++ {
		acc := &AccoundData{
			Type:           "accounts",
			ID:             fmt.Sprintf("fakeUserID%d", i),
			OrganisationID: uuid.New().String(),
			Version:        0,
		}
		accs[i] = acc
	}

	accList := &AccountList{
		Accounts: accs,
		Links: &LinkList{
			First: "fakeFirstLint",
			Last:  "fakeLastLink",
			Self:  "fakeSelfLink",
		},
	}
	rawAccount, err := json.Marshal(accList)
	if err != nil {
		t.Fatalf("unexpected error marshalling account, error %v", err)
	}

	h := &fakeHTTPClient{
		body: rawAccount,
	}
	api := NewAPIClient(h)

	pags := NewPagination(0, 10)
	a, err := api.List(context.Background(), pags)
	if err != nil {
		t.Fatalf("error fetching user, error %v", err)
	}

	if !reflect.DeepEqual(accList, a) {
		t.Error("account slice are not equal")
	}
}

func TestListAccountsReturnAnErrorOnUnexpectedStatusCode(t *testing.T) {
	h := &fakeHTTPClient{
		err: client.ErrContentNotFound,
	}
	api := NewAPIClient(h)

	pags := NewPagination(0, 10)
	_, err := api.List(context.Background(), pags)
	if err == nil {
		t.Fatal("expected not found error")
	}

	if !errors.Is(err, client.ErrContentNotFound) {
		t.Errorf("error is not content not found got %v", err)
	}
}

func TestDeleteAccountReturnsNoErrorOnSuccessfulRequest(t *testing.T) {
	userID := uuid.New().String()
	h := &fakeHTTPClient{}
	api := NewAPIClient(h)

	version := 0
	err := api.Delete(context.Background(), userID, version)
	if err != nil {
		t.Errorf("unexpected error deleting, got %v", err)
	}
}

func TestDeleteAccountReturnsNotFoundErrorOnStatusNotFound(t *testing.T) {
	userID := uuid.New().String()
	h := &fakeHTTPClient{
		err: client.ErrContentNotFound,
	}
	api := NewAPIClient(h)

	version := 0
	err := api.Delete(context.Background(), userID, version)
	if err == nil {
		t.Fatal("expected not found error")
	}

	if !errors.Is(err, client.ErrContentNotFound) {
		t.Errorf("error is not content not found got %v", err)
	}
}

func TestDeleteAccountReturnsConflictErrorOnWrongVersion(t *testing.T) {
	userID := uuid.New().String()
	h := &fakeHTTPClient{
		err: ErrVersionConflict,
	}
	api := NewAPIClient(h)

	version := 99
	err := api.Delete(context.Background(), userID, version)
	if err == nil {
		t.Fatal("expected not found error")
	}

	if !errors.Is(err, ErrVersionConflict) {
		t.Errorf("error is not content not found got %v", err)
	}
}

type fakeHTTPClient struct {
	statusCode int
	body       []byte
	err        error
}

func (f *fakeHTTPClient) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if v != nil && f.body != nil {
		err := json.Unmarshal(f.body, &v)
		if err != nil {
			return nil, err
		}
	}
	return &http.Response{
		StatusCode: f.statusCode,
		Body:       ioutil.NopCloser(bytes.NewBuffer(f.body)),
		Header:     make(http.Header),
	}, f.err
}

func (f *fakeHTTPClient) CreateRequest(method, url string, body interface{}) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}
