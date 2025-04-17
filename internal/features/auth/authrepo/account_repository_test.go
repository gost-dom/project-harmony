package authrepo_test

import (
	"reflect"
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

type TimeoutTest struct {
	t         testing.TB
	f         func()
	OnTimeout func()
}

func withTimeout(t testing.TB, f func()) TimeoutTest {
	return TimeoutTest{t, f, nil}
}

func (t TimeoutTest) Run() {

	c := make(chan struct{})
	timeout := time.After(time.Second)

	go func() {
		t.f()
		close(c)
	}()
	select {
	case <-c:
	case <-timeout:
		if t.OnTimeout != nil {
			t.OnTimeout()
		} else {
			t.t.Error("Timeout")
		}
	}
}

func TestInsertDomainEvents(t *testing.T) {
	var actual []domain.Event
	tt := withTimeout(t, func() {
		ctx := t.Context()
		repo := initRepository()
		ch, err := couchdb.DefaultConnection.StartListener(ctx)
		assert.NoError(t, err)

		// Insert an entity with two domain events
		acc := auth.UseCaseOfEntity(domaintest.InitPasswordAuthAccount())
		event1 := authdomain.CreateValidationRequestEvent(acc.Entity.Account)
		event2 := authdomain.CreateAccountRegisteredEvent(acc.Entity.Account)
		acc.AddEvent(event1)
		acc.AddEvent(event2)
		assert.NoError(t, repo.Insert(ctx, acc))

		// Wait for the domain events to appear. Ignore other events,
		expected := []domain.Event{event1, event2}
		for e := range ch {
			if reflect.DeepEqual(e, event1) || reflect.DeepEqual(e, event2) {
				actual = append(actual, e)
			}
			if len(actual) == 2 {
				assert.ElementsMatch(t, expected, actual)
				break
			}
		}
	})
	tt.OnTimeout = func() { t.Errorf("Failed finding all expected events. Found %+v", actual) }
	tt.Run()
}
