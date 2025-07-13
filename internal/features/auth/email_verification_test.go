package auth_test

import (
	"context"
	"errors"
	"fmt"
	"harmony/internal/core/corerepo"
	"harmony/internal/domain"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/messaging"
	"harmony/internal/messaging/ioc"
	"harmony/internal/testing/domaintest"
	"harmony/internal/testing/mailhog"
	"net/mail"
	"reflect"
	"testing"
	"time"

	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
	"github.com/stretchr/testify/assert"
)

type domainEvt map[domain.EventID]domain.Event

func (e domainEvt) Update(ctx context.Context, event domain.Event) (domain.Event, error) {
	e[event.ID] = event
	return event, nil
}

func TestEmailValidatorValidate(t *testing.T) {
	addr, err := mail.ParseAddress("jd@example.com")
	assert.NoError(t, err, "error parsing email in test")

	t.Run("Passing an invalid code", func(t *testing.T) {
		acc := domaintest.InitAccount(domaintest.WithEmailAddress(addr))
		acc.StartEmailValidationChallenge()
		repo := NewAccountRepositoryStub(t, &acc)
		validator := auth.EmailChallengeValidator{Repository: repo}

		got, err := validator.Validate(t.Context(), auth.ValidateEmailInput{
			Email: addr,
			Code:  authdomain.EmailValidationCode("invalid-code"),
		})

		if assert.ErrorIs(t, err, auth.ErrBadChallengeResponse, "Validate error result") {
			assert.Zero(t, got, "Failed response should result in zero value")
		}
	})

	t.Run("Passing an invalid email", func(t *testing.T) {
		acc := domaintest.InitAccount(domaintest.WithEmailAddress(addr))
		acc.StartEmailValidationChallenge()
		repo := NewAccountRepositoryStub(t, &acc)
		validator := auth.EmailChallengeValidator{Repository: repo}

		em := domaintest.InitEmail()
		em.NewChallenge()
		got, err := validator.Validate(t.Context(), auth.ValidateEmailInput{
			Email: &em.Address,
			Code:  em.Challenge.Code,
		})

		if assert.ErrorIs(t, err, auth.ErrNotFound, "Validate error result") {
			assert.Zero(t, got, "Failed response should result in zero value")
		}
	})

	t.Run("Passing the valid code", func(t *testing.T) {
		acc := domaintest.InitAccount(domaintest.WithEmailAddress(addr))
		acc.StartEmailValidationChallenge()
		repo := NewAccountRepositoryStub(t, &acc)
		validator := auth.EmailChallengeValidator{Repository: repo}

		got, err := validator.Validate(t.Context(), auth.ValidateEmailInput{
			Email: &acc.Email.Address,
			Code:  acc.Email.Challenge.Code,
		})

		assert.NoError(t, err)
		assert.Equal(
			t, acc.ID, got.ID,
			"The authenticated account should be returned",
		)
		assert.True(t, acc.Email.Validated, "Account is validated")

		assert.Equal(t, *got.Account, acc, "Account was updated in repository")
	})
}

func TestSendEmailValidationChallenge(t *testing.T) {
	assert.NoError(t, mailhog.DeleteAll())

	acc := domaintest.InitAccount(func(acc *authdomain.Account) {
		acc.DisplayName = "John"
		acc.Name = "John Smith"
	})
	event := acc.StartEmailValidationChallenge()
	assert.False(t, acc.Validated(), "guard: account should be an invalidated account")

	domainEvents := domainEvt{}
	graph := surgeon.Replace[auth.AccountLoader](ioc.Graph, NewAccountRepositoryStub(t, &acc))
	graph = surgeon.Replace[messaging.DomainEventUpdater](graph, domainEvents)
	v := graph.Instance()

	assert.NoError(t, v.ProcessDomainEvent(t.Context(), event))

	g := gomega.NewWithT(t)
	g.Expect(
		mailhog.GetAll(),
	).To(gomega.ContainElement(HaveHeader("To", MatchEmailAddress(acc.Email.Address.Address))))

	assert.NotNil(t, domainEvents[event.ID].PublishedAt, "Domain event marked as published")
}

func TestIntegrationSendEmailValidationChallenge(t *testing.T) {
	// This test should be more general. This verifies that domain events marked
	// as published are not returned when starting a channel of unpublished
	// events after the event was processed.
	if testing.Short() {
		t.SkipNow()
	}
	ctx := t.Context()

	acc1 := domaintest.InitAccount()
	acc2 := domaintest.InitAccount()
	event1 := acc1.StartEmailValidationChallenge()
	event2 := acc2.StartEmailValidationChallenge()

	event1, err1 := corerepo.DefaultDomainEventRepo.Insert(ctx, event1)
	event2, err2 := corerepo.DefaultDomainEventRepo.Insert(ctx, event2)
	assert.NoError(t, errors.Join(err1, err2))

	graph := surgeon.Replace[auth.AccountLoader](ioc.Graph, NewAccountRepositoryStub(t, &acc1))
	v := graph.Instance()

	assert.NoError(t, v.ProcessDomainEvent(t.Context(), event1))

	// This channel should not receive the published domain event, as we start listening
	// after it was published
	ch, err := corerepo.DefaultDomainEventRepo.StreamOfEvents(ctx)
	assert.NoError(t, err)

	timeout := time.After(1000 * time.Millisecond)
	func() {
		for {
			select {
			case <-timeout:
				t.Error("Timeout waiting for event")
				return
			case e := <-ch:
				// This test relies on change events occurring in the same order
				// they were inserted, i.e., if event1 is sent to the channel,
				// it would happen _before_ event2, so when we see event2,
				// without having seen event1, it is not an error.
				if reflect.DeepEqual(event1, e) {
					t.Errorf("Processed event should not be in event stream: %+v", e)
					return
				}
				if reflect.DeepEqual(event2, e) {
					return
				}
			}
		}
	}()
}

func MatchEmailAddress(expected string) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(actual string) (bool, error) {
		addr, err := mail.ParseAddress(actual)
		return addr.Address == expected, err
	})
}

func HaveHeader(key string, expected any) types.GomegaMatcher {
	matcher, ok := expected.(types.GomegaMatcher)
	if !ok {
		matcher = gomega.Equal(expected)
	}
	return gomega.WithTransform(func(m mailhog.MailhogMessage) ([]string, error) {
		res, ok := m.Content.Headers[key]
		if !ok {
			return nil, fmt.Errorf("Message did not contain header: %s", key)
		}
		return res, nil
	}, gomega.ContainElement(matcher))
}
