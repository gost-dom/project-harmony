package authrepo_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"harmony/internal/core/corerepo"
	"harmony/internal/couchdb"
	"harmony/internal/domain"
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
	reloaded, err := repo.Get(t.Context(), acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, acc.Account, reloaded)

	foundByEmail, err := repo.FindByEmail(t.Context(), acc.Email.String())
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

func TestUpdate(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	email := domaintest.NewAddress()
	pwacc := auth.UseCaseOfEntity(domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email)))
	err := repo.Insert(ctx, pwacc)
	assert.NoError(t, err, "Error inserting account")

	acc, err := repo.Get(ctx, pwacc.Entity.Account.ID)
	assert.NoError(t, err, "Error loading account")

	acc.DisplayName = "New name"
	reloaded, err := repo.Update(ctx, acc)
	assert.NoError(t, err, "Error updating account")
	assert.Equal(t, "New name", reloaded.DisplayName)
}

type TimeoutTest struct {
	t testing.TB
	f func(context.Context)
}

func withTimeout(t testing.TB, f func(ctx context.Context)) TimeoutTest {
	return TimeoutTest{t, f}
}

func (t TimeoutTest) Run() {
	t.RunWithErrorf("Timeout")
}

func (t TimeoutTest) RunWithErrorf(format string, args ...any) {
	t.t.Helper()
	ctx, cancel := context.WithTimeout(t.t.Context(), time.Second)

	go func() {
		defer cancel()
		t.f(ctx)
	}()

	<-ctx.Done()
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.t.Errorf(format, args...)
	}
}

func TestInsertDomainEvents(t *testing.T) {
	var actual []domain.Event
	withTimeout(t, func(ctx context.Context) {
		repo := initRepository()
		coreRepo := corerepo.DefaultMessageSource
		assert.NoError(t, coreRepo.StartListener(ctx))

		// Insert an entity with two domain events
		acc := auth.UseCaseOfEntity(domaintest.InitPasswordAuthAccount())
		event1 := acc.Entity.StartEmailValidationChallenge()
		event2 := authdomain.CreateAccountRegisteredEvent(acc.Entity.Account)
		acc.AddEvent(event1)
		acc.AddEvent(event2)
		assert.NoError(t, repo.Insert(ctx, acc))

		ch, err := corerepo.DefaultDomainEventRepo.StreamOfEvents(ctx)
		assert.NoError(t, err)

		// Wait for the domain events to appear. Ignore other events,
		expected := []domain.Event{event1, event2}
		for e := range ch {
			e.Rev = ""
			if reflect.DeepEqual(e, event1) || reflect.DeepEqual(e, event2) {
				actual = append(actual, e)
			}
			if len(actual) == 2 {
				assert.ElementsMatch(t, expected, actual)
				break
			}
		}
	}).RunWithErrorf("Failed finding all expected events. Found %+v", actual)
}
