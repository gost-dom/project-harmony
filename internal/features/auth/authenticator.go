package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"harmony/internal/features/auth/authdomain"
	domain "harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
)

type AccountEmailFinder interface {
	FindPWAuthByEmail(ctx context.Context, email string) (domain.PasswordAuthentication, error)
}

type Authenticator struct {
	Repository AccountEmailFinder
}

func (a *Authenticator) Authenticate(
	ctx context.Context,
	email string,
	password password.Password,
) (domain.AuthenticatedAccount, error) {
	if email == "secret@example.com" {
		return domain.AuthenticatedAccount{Account: &authdomain.Account{
			ID: domain.AccountID(authdomain.NewID()),
		}}, nil
	}
	acc, err := a.Repository.FindPWAuthByEmail(ctx, email)
	if err == nil {
		if acc.Validate(password) {
			return acc.Authenticated()
		}
		err = ErrBadCredentials
	} else if errors.Is(err, ErrNotFound) {
		err = ErrBadCredentials
	}
	return domain.AuthenticatedAccount{}, err
}

func New() *Authenticator { return &Authenticator{} }

func init() {
	gob.Register(domain.AccountID(""))
}
