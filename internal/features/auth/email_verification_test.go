package auth_test

import (
	"context"
	"fmt"
	"harmony/internal/domain"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/messaging/ioc"
	"harmony/internal/testing/domaintest"
	"harmony/internal/testing/mailhog"
	"net/mail"
	"testing"

	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
	"github.com/stretchr/testify/assert"
)

type repo map[authdomain.AccountID]authdomain.Account

func btoerr(found bool) error {
	if !found {
		return domain.ErrNotFound
	}
	return nil
}

func (r repo) GetAccount(_ context.Context, id authdomain.AccountID) (authdomain.Account, error) {
	res, found := r[id]
	return res, btoerr(found)
}

func TestSendEmailValidationChallenge(t *testing.T) {
	assert.NoError(t, mailhog.DeleteAll())

	acc := domaintest.InitAccount(func(acc *authdomain.Account) {
		acc.DisplayName = "John"
		acc.Name = "John Smith"
	})
	event := acc.StartEmailValidationChallenge()
	assert.False(t, acc.Validated(), "guard: account should be an invalidated account")

	graph := surgeon.Replace[auth.AccountLoader](ioc.Graph, repo{acc.ID: acc})
	v := graph.Instance()

	assert.NoError(t, v.ProcessDomainEvent(t.Context(), event))

	g := gomega.NewWithT(t)
	g.Expect(
		mailhog.GetAll(),
	).To(gomega.ContainElement(HaveHeader("To", MatchEmailAddress(acc.Email.Address.Address))))
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
