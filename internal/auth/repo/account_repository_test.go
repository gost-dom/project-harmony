package repo_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"harmony/internal/auth"
	"harmony/internal/auth/domain"
	"harmony/internal/auth/domain/password"
	. "harmony/internal/auth/repo"
	"harmony/internal/core"
	"harmony/internal/core/corerepo"
	_ "harmony/internal/testing/couchtest" // clear database before tests
	"harmony/internal/testing/domaintest"

	"github.com/stretchr/testify/assert"
)

func initRepository() AccountRepository {
	corerepo.AssertInitialized()
	conn := corerepo.DefaultConnection
	return AccountRepository{conn}
}

func insertAccount(c context.Context, repo AccountRepository, acc auth.AccountUseCaseResult) error {
	_, err := repo.Insert(c, acc)
	return err
}

func TestAccountRoundtrip(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	acc := domaintest.InitPasswordAuthAccount(domaintest.WithPassword("foobar"))
	uc := auth.AccountUseCaseResult{Entity: acc}
	inserted, err := repo.Insert(ctx, uc)
	assert.NoError(t, err, "Error inserting account")
	reloaded, err := repo.Get(t.Context(), acc.ID)
	assert.NoError(t, err)
	assert.Equal(t,
		inserted.Account, reloaded,
		"The account retrieved by a GET should be identical to its state right after insert")

	t.Run("FindPWAuthByEmail", func(t *testing.T) {
		foundByEmail, err := repo.FindPWAuthByEmail(t.Context(), acc.Email.String())
		assert.NoError(t, err, "Error finding by email")
		assert.Equal(t, inserted, foundByEmail, "Entity found by email")
		assert.True(t, foundByEmail.Validate(password.Parse("foobar")), "Password validates")
	})

	t.Run("FindByEmail", func(t *testing.T) {
		foundByEmail, err := repo.FindByEmail(t.Context(), acc.Email.String())
		assert.NoError(t, err, "Error finding by email")
		assert.Equal(t, inserted.Account, foundByEmail, "Entity found by email")
	})
}

func TestDuplicateEmail(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	email := domaintest.NewAddress()
	acc1 := core.UseCaseOfEntity(
		domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email)))
	acc2 := core.UseCaseOfEntity(
		domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email)))
	assert.NoError(t, insertAccount(ctx, repo, acc1))
	assert.ErrorIs(t, insertAccount(ctx, repo, acc2), ErrConflict)
}

func TestAccountRepositoryUpdate(t *testing.T) {
	ctx := t.Context()
	repo := initRepository()

	email := domaintest.NewAddress()
	pwacc := core.UseCaseOfEntity(domaintest.InitPasswordAuthAccount(domaintest.WithEmail(email)))
	tmp, err := repo.Insert(ctx, pwacc)
	inserted := tmp.Account
	if !assert.NoError(t, err, "Error inserting account") {
		return
	}

	if !t.Run("Modify the returned value and update", func(t *testing.T) {
		inserted.DisplayName = "New name"
		reloaded, err := repo.Update(ctx, inserted)
		assert.NoError(t, err, "Error updating account")
		assert.Equal(t, "New name", reloaded.DisplayName)
		assert.NotEqual(
			t,
			reloaded.Rev,
			inserted.Rev,
			"After update, the revision should have changed",
		)
		assert.NotEmpty(t, reloaded.Rev, "Update returns a revision")

		reloaded.DisplayName = "2nd update"
		_, err = repo.Update(ctx, reloaded)
		assert.NoError(t, err, "Error updating the value returned from Update")
	}) {
		return
	}

	if !t.Run("Modify the original returned value, which is now stale", func(t *testing.T) {
		inserted.DisplayName = "2nd name"
		_, err = repo.Update(ctx, inserted)
		assert.ErrorIs(t, err, ErrConflict, "Update should fail with a conflict error")
	}) {
		return
	}

	t.Run("Update document returned from Get", func(t *testing.T) {
		reloaded, err := repo.Get(ctx, inserted.ID)
		assert.NoError(t, err, "Error loading account")
		reloaded.DisplayName = "Update after Get"
		_, err = repo.Update(ctx, reloaded)
		assert.NoError(t, err, "Error updating after Get")
	})
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
	var actual []core.DomainEvent
	withTimeout(t, func(ctx context.Context) {
		repo := initRepository()
		coreRepo := corerepo.DefaultMessageSource
		assert.NoError(t, coreRepo.StartListener(ctx))

		// Insert an entity with two domain events
		acc := core.UseCaseOfEntity(domaintest.InitPasswordAuthAccount())
		event1 := acc.Entity.StartEmailValidationChallenge()
		event2 := domain.CreateAccountRegisteredEvent(acc.Entity.Account)
		acc.AddEvent(event1)
		acc.AddEvent(event2)
		assert.NoError(t, insertAccount(ctx, repo, acc))

		ch, err := corerepo.DefaultDomainEventRepo.StreamOfEvents(ctx)
		assert.NoError(t, err)

		// Wait for the domain events to appear. Ignore other events,
		expected := []core.DomainEvent{event1, event2}
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
