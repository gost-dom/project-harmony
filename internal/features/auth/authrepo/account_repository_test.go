package authrepo_test

import (
	"harmony/internal/couchdb"
	"harmony/internal/features/auth/authdomain/password"
	. "harmony/internal/features/auth/authrepo"
	"harmony/internal/testing/domaintest"
	"testing"

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
	assert.NoError(t, repo.Insert(t.Context(), acc))
	reloaded, err := repo.Get(acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.Account, reloaded)

	foundByEmail, err := repo.FindByEmail(acc.Email.Address)
	assert.NoError(t, err, "Error finding by email")
	assert.Equal(t, acc, foundByEmail, "Entity found by email")
	assert.True(t, foundByEmail.Validate(password.Parse("foobar")), "Password validates")
}

func TestDuplicateEmail(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	email := domaintest.NewAddress()
	acc1 := domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email))
	acc2 := domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email))
	assert.NoError(t, repo.Insert(ctx, acc1))
	assert.ErrorIs(t, repo.Insert(ctx, acc2), ErrConflict)
}
