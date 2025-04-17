package authrepo_test

import (
	"testing"
	"time"

	"harmony/internal/couchdb"
	"harmony/internal/domain"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	. "harmony/internal/features/auth/authrepo"
	_ "harmony/internal/testing/couchtest" // clear database before tests
	"harmony/internal/testing/domaintest"

	"github.com/onsi/gomega"
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

	var err error
	ch, closer, err := couchdb.DefaultConnection.StartListener(ctx)
	assert.NoError(t, err)
	defer closer.Close()
	// defer func() {
	// 	close()
	// }()
	acc := auth.UseCaseOfEntity(domaintest.InitPasswordAuthAccount())
	event := authdomain.CreateValidationRequestEvent(acc.Entity.Account)
	event2 := authdomain.CreateAccountRegisteredEvent(acc.Entity.Account)
	acc.AddEvent(event)
	acc.AddEvent(event2)
	assert.NoError(t, repo.Insert(ctx, acc))
	// var res couchdb.ViewResult[domain.Event]
	// v := make(url.Values)
	// v.Set("key", `"`+string(event.ID)+`"`)
	// _, err = repo.Connection.GetPath(
	// 	"_design/events/_view/unpublished_events",
	// 	v,
	// 	&res,
	// )
	assert.NoError(t, err)
	// assert.Equal(t, []domain.Event{event}, res.Values())

	var actual []domain.Event
	select {
	case actual1 := <-ch:
		actual = append(actual, actual1)
		if len(actual) == 2 {
			gomega.NewWithT(t).Expect(actual).To(gomega.ConsistOf(
				event, event2))
		}
	case <-time.After(time.Second):
		t.Errorf("Timeout waiting for event. Existing events: %v", actual)
	}
	// actual2 := <-ch
	// actual := []domain.Event{actual1}
	// gomega.NewWithT(t).Expect(actual).To(gomega.ConsistOf(
	// 	event, event2))
	// time.Sleep(time.Second / 2)
}
