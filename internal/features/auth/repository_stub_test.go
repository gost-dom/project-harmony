package auth_test

import (
	"context"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/testing/repotest"
	"testing"
)

type InsertAccountTranslator struct{}

func (t InsertAccountTranslator) ID(e authdomain.PasswordAuthentication) string {
	return string(e.Account.ID)
}

type AccountRepositoryStub struct {
	repotest.RepositoryStub[authdomain.PasswordAuthentication]
}

func NewAccountRepoStub(t testing.TB) *AccountRepositoryStub {
	return &AccountRepositoryStub{repotest.NewRepositoryStub(t, InsertAccountTranslator{})}
}

func (i AccountRepositoryStub) FindByEmail(
	ctx context.Context,
	email string,
) (authdomain.PasswordAuthentication, error) {
	for _, v := range i.Entities {
		if v.Email.Equals(email) {
			return *v, nil
		}
	}
	return authdomain.PasswordAuthentication{}, auth.ErrNotFound
}

func (s AccountRepositoryStub) Update(
	_ context.Context,
	acc authdomain.Account,
) (authdomain.Account, error) {
	existing, ok := s.Entities[string(acc.ID)]
	if !ok {
		return authdomain.Account{}, auth.ErrNotFound
	}
	existing.Account = acc
	return acc, nil
}
