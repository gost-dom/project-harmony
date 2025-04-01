package auth

import (
	"context"
	"encoding/gob"
	"errors"
	. "harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var ErrBadCredentials = errors.New("authenticate: Bad credentials")
var ErrNotFound = errors.New("Not found")

type AccountEmailFinder interface {
	FindByEmail(ctx context.Context, id string) (PasswordAuthentication, error)
}

type Authenticator struct {
	Repository AccountEmailFinder
}

func NewID() (string, error) { return gonanoid.New(32) }

func (a *Authenticator) Authenticate(
	ctx context.Context,
	username string,
	password password.Password,
) (account AuthenticatedAccount, err error) {
	var tmp PasswordAuthentication
	if tmp, err = a.Repository.FindByEmail(ctx, username); errors.Is(err, ErrNotFound) {
		err = ErrBadCredentials
	}
	account.Account = tmp.Account
	if err == nil && (!account.Email.Validated) {
		err = ErrAccountEmailNotValidated
	}
	if err == nil {
		if !tmp.Validate(password) {
			err = ErrBadCredentials
		}
	}
	return
}

func New() *Authenticator { return &Authenticator{} }

func init() {
	gob.Register(AccountID(""))
}
