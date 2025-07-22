package router_test

import (
	"context"
	"harmony/internal/auth/domain"
	"harmony/internal/auth/router"
	"harmony/internal/core"
	"harmony/internal/testing/domaintest"
	"harmony/internal/testing/servertest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type repo struct {
	account domain.Account
}

func (r repo) Get(_ context.Context, id domain.AccountID) (domain.Account, error) {
	if id == r.account.ID {
		return r.account, nil
	}
	return domain.Account{}, core.ErrNotFound
}

func TestSessionManagerReturnsAccount(t *testing.T) {
	acc := domaintest.InitAuthenticatedAccount(domaintest.WithEmail("jd@example.com"))
	if !assert.NotZero(t, acc.ID, "Verification is meaningless if the account ID is zero") {
		return
	}

	store := servertest.NewMemStore()
	mgr := router.SessionManager{SessionStore: store, Repo: repo{*acc.Account}}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	assert.NoError(t, mgr.SetAccount(w, r, acc))

	for _, c := range w.Result().Cookies() {
		r.AddCookie(c)
	}

	got := mgr.LoggedInUser(r)
	assert.NotNil(t, got)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, "jd@example.com", got.Email.String())
}
