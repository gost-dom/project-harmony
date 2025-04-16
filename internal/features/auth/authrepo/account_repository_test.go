package authrepo_test

import (
	"fmt"
	"testing"

	"harmony/internal/couchdb"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	. "harmony/internal/features/auth/authrepo"
	_ "harmony/internal/testing/couchtest" // clear database before tests
	"harmony/internal/testing/domaintest"

	"github.com/stretchr/testify/assert"
)

func initRepository() AccountRepository {
	couchdb.AssertInitialized()
	conn := couchdb.DefaultConnection
	return AccountRepository{conn}
}

func TestAccountRoundtrip(t *testing.T) {
	repo := initRepository()

	acc := domaintest.InitPasswordAuthAccount(domaintest.WithPassword("foobar"))
	uc := auth.AccountUseCaseResult{Entity: acc}
	assert.NoError(t, repo.Insert(t.Context(), uc))
	reloaded, err := repo.Get(acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.Account, reloaded)

	foundByEmail, err := repo.FindByEmail(acc.Email.String())
	assert.NoError(t, err, "Error finding by email")
	assert.Equal(t, acc, foundByEmail, "Entity found by email")
	assert.True(t, foundByEmail.Validate(password.Parse("foobar")), "Password validates")
}

func TestDuplicateEmail(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	email := domaintest.NewAddress()
	acc1 := auth.UseCaseOfEntity(
		domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email)))
	acc2 := auth.UseCaseOfEntity(
		domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email)))
	assert.NoError(t, repo.Insert(ctx, acc1))
	assert.ErrorIs(t, repo.Insert(ctx, acc2), ErrConflict)
}

func TestInsertDomainEvents(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	acc := auth.UseCaseOfEntity(domaintest.InitPasswordAuthAccount())
	event := authdomain.CreateValidationRequestEvent(acc.Entity.Account)
	acc.AddEvent(event)
	assert.NoError(t, repo.Insert(ctx, acc))
	var res couchdb.ViewResult[authdomain.EmailValidationRequest]
	_, err := repo.Connection.Get("_design/events/_view/unpublished_events", &res)
	fmt.Printf("%+v", res.Rows)
	assert.NoError(t, err)
}
