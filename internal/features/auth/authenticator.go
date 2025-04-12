package auth

import (
	"context"
	"encoding/gob"
	"errors"
	domain "harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type AccountEmailFinder interface {
	FindByEmail(ctx context.Context, email string) (domain.PasswordAuthentication, error)
}

type Authenticator struct {
	Repository AccountEmailFinder
}

func NewID() (string, error) { return gonanoid.New(32) }

func (a *Authenticator) Authenticate(
	ctx context.Context,
	email string,
	password password.Password,
) (domain.AuthenticatedAccount, error) {
	acc, err := a.Repository.FindByEmail(ctx, email)
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
