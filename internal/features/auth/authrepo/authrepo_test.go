package authrepo_test

import (
	"harmony/internal/couchdb"
	"harmony/internal/features/auth/authdomain/password"
	. "harmony/internal/features/auth/authrepo"
	"harmony/internal/testing/domaintest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountRoundtrip(t *testing.T) {
	conn, err := couchdb.NewCouchConnection("http://admin:password@localhost:5984/harmony")
	assert.NoError(t, err)
	repo := AccountRepository{conn}
	acc := domaintest.InitPasswordAuthAccount(domaintest.WithPassword("foobar"))
	assert.NoError(t, repo.Insert(acc))
	reloaded, err := repo.Get(acc.ID)
	assert.Equal(t, acc.Account, reloaded)

	foundByEmail, err := repo.FindByEmail(acc.Email.Address)
	assert.NoError(t, err, "Error finding by email")
	assert.Equal(t, acc, foundByEmail, "Entity found by email")
	assert.True(t, foundByEmail.Validate(password.Parse("foobar")), "Password validates")
}

func TestDuplicateEmail(t *testing.T) {
	email := domaintest.NewAddress()
	acc1 := domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email))
	acc2 := domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email))
	conn, err := couchdb.NewCouchConnection("http://admin:password@localhost:5984/harmony")
	assert.NoError(t, err)
	repo := AccountRepository{conn}
	assert.NoError(t, repo.Insert(acc1))
	assert.ErrorIs(t, repo.Insert(acc2), ErrConflict)
}
