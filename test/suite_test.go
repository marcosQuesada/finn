// +build integration

package test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcosQuesada/finn"
	"github.com/marcosQuesada/finn/http"
)

var userID = uuid.New().String()

func TestAccountSuite(t *testing.T) {
	t.Run("CreateAccountWithValidParametersDoesNotThrowError", testCreateAccountWithValidParametersDoesNotThrowError)
	t.Run("FetchAccountOnAlreadyCreatedUserDoesNotThrowError", testFetchAccountOnAlreadyCreatedUserDoesNotThrowError)
	t.Run("FetchAccountOnNonCreatedUserThrowNonExistentError", testFetchAccountOnNonCreatedUserThrowNonExistentError)
	t.Run("DeleteAccountOnNonExistentAccountDoesThrowNotFoundError", testDeleteAccountOnNonExistentAccountThrowsNotFoundError)
	t.Run("DeleteAccountOnInvalidVersionThrowConflictError", testDeleteAccountOnInvalidVersionThrowConflictError)
	t.Run("DeleteAccountOnExistentAccountDoesNotThrowError", testDeleteAccountOnExistentAccountDoesNotThrowError)
	t.Run("ListAccountsWithPaginationDoesNotThrowErrorAndReturnsASubSetOfResults", testListAccountsWithPaginationDoesNotThrowErrorAndReturnsASubSetOfResults)
}

func testCreateAccountWithValidParametersDoesNotThrowError(t *testing.T) {
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	acc := &finn.Account{
		AccoundData: &finn.AccoundData{
			Type:           "accounts",
			ID:             userID,
			OrganisationID: uuid.New().String(),
			Version:        0,
			Attributes: &finn.Attributes{
				Country:      "ES",
				BaseCurrency: "EUR",
			},
		}}

	ac, err := a.Create(ctx, acc)
	if err != nil {
		t.Fatalf("unexpected error creating account, error %v", err)
	}

	if ac == nil || ac.AccoundData == nil {
		t.Fatal("nil account")
	}
	if ac.AccoundData.ID != userID {
		t.Errorf("Account userID does not match, expected %s got %s", userID, ac.AccoundData.ID)
	}
}

func testFetchAccountOnAlreadyCreatedUserDoesNotThrowError(t *testing.T) {
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	acc, err := a.Fetch(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error creating account, error %v", err)
	}

	if acc == nil || acc.AccoundData == nil {
		log.Fatal("Fetched account is nil")
	}

	if acc.AccoundData.ID != userID {
		t.Errorf("Account userID does not match, expected %s got %s", userID, acc.AccoundData.ID)
	}
}

func testFetchAccountOnNonCreatedUserThrowNonExistentError(t *testing.T) {
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userNotFound := uuid.New().String()
	_, err := a.Fetch(ctx, userNotFound)
	if err == nil {
		t.Fatal("expected user not found error")
	}

	if !errors.Is(err, http.ErrContentNotFound) {
		t.Errorf("unexpected error type, expected content not found, error %v", err)
	}
}

func testListAccountsWithPaginationDoesNotThrowErrorAndReturnsASubSetOfResults(t *testing.T) {
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	total := 15
	userIDs, err := setup(ctx, a, total)
	if err != nil {
		t.Fatalf("unexpected error populating db, error %v", err)
	}

	p := finn.NewPagination(0, 10)
	acc, err := a.List(ctx, p)
	if err != nil {
		t.Fatalf("unexpected error creating account, error %v", err)
	}

	if acc == nil || acc.Accounts == nil {
		t.Fatalf("account list is nil")
	}

	if got, want := len(acc.Accounts), p.Size; got != want {
		t.Fatalf("unexpected account size, expected %d got %d", want, got)
	}

	p = finn.NewPagination(1, 10)
	acc, err = a.List(ctx, p)
	if err != nil {
		t.Fatalf("unexpected error creating account, error %v", err)
	}

	if got, want := len(acc.Accounts), total-10; got != want {
		t.Fatalf("unexpected account size, expected %d got %d", want, got)
	}

	err = tearDown(ctx, a, userIDs)
	if err != nil {
		t.Errorf("unexpected error doing tear down, error %v", err)
	}
}

func testDeleteAccountOnExistentAccountDoesNotThrowError(t *testing.T) {
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	version := 0
	err := a.Delete(ctx, userID, version)
	if err != nil {
		t.Fatalf("unexpected error creating account, error %v", err)
	}
}

func testDeleteAccountOnNonExistentAccountThrowsNotFoundError(t *testing.T) {
	t.Skip("api server returning 204 status codes on expected not existent accounts")
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	version := 0
	err := a.Delete(ctx, uuid.New().String(), version)
	if err == nil {
		t.Fatal("Expected content not found")
	}

	if !errors.Is(err, http.ErrContentNotFound) {
		t.Errorf("Unexpected error response, expected conflict got %v", err)
	}
}

func testDeleteAccountOnInvalidVersionThrowConflictError(t *testing.T) {
	t.Skip("api server returning 404 status codes on expected conflict error")
	cl := http.NewClient()
	a := finn.NewAPIClient(cl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	version := 2
	err := a.Delete(ctx, userID, version)
	if err == nil {
		t.Fatal("expected conflict error")
	}

	if !errors.Is(err, finn.ErrVersionConflict) {
		t.Errorf("unexpected error type, exepcted conflict got %v", err)
	}
}

func setup(ctx context.Context, cl *finn.APIClient, size int) ([]string, error) {
	userIDs := make([]string, 0)
	for i := 0; i < size; i++ {
		acc := &finn.Account{
			AccoundData: &finn.AccoundData{
				Type:           "accounts",
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Version:        0,
				Attributes: &finn.Attributes{
					Country:      "ES",
					BaseCurrency: "EUR",
				},
			}}

		_, err := cl.Create(ctx, acc)
		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, acc.AccoundData.ID)
	}

	return userIDs, nil
}

func tearDown(ctx context.Context, cl *finn.APIClient, userIDs []string) error {
	version := 0
	for _, u := range userIDs {
		err := cl.Delete(ctx, u, version)
		if err != nil {
			return err
		}
	}

	return nil
}
