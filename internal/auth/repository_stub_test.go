package auth_test

import (
	"context"
	"harmony/internal/auth"
	domain "harmony/internal/auth/authdomain"
	"harmony/internal/testing/repotest"
	"testing"
)

type PWAuthTranslator struct{}

func (t PWAuthTranslator) ID(e domain.PasswordAuthentication) domain.AccountID {
	return e.Account.ID
}

type PWAuthRepositoryStub struct {
	repotest.RepositoryStub[domain.PasswordAuthentication, domain.AccountID]
}

func NewPWAuthRepositoryStub(t testing.TB) *PWAuthRepositoryStub {
	return &PWAuthRepositoryStub{repotest.NewRepositoryStub(t, PWAuthTranslator{})}
}

func (i PWAuthRepositoryStub) FindPWAuthByEmail(
	ctx context.Context, email string,
) (domain.PasswordAuthentication, error) {
	for _, v := range i.Entities {
		if v.Email.Equals(email) {
			return *v, nil
		}
	}
	return domain.PasswordAuthentication{}, auth.ErrNotFound
}

type AccountTranslator struct{}

func (t AccountTranslator) ID(e domain.Account) domain.AccountID {
	return e.ID
}

type AccountRepositoryStub struct {
	repotest.RepositoryStub[domain.Account, domain.AccountID]
}

func NewAccountRepositoryStub(t testing.TB, acc ...*domain.Account) *AccountRepositoryStub {
	return &AccountRepositoryStub{repotest.NewRepositoryStub(t, AccountTranslator{}, acc...)}
}

func (i AccountRepositoryStub) FindByEmail(
	ctx context.Context, email string,
) (domain.Account, error) {
	for _, v := range i.Entities {
		if v.Email.Equals(email) {
			return *v, nil
		}
	}
	return domain.Account{}, auth.ErrNotFound
}
