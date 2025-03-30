package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"harmony/internal/features/auth/authdomain"
	. "harmony/internal/features/auth/authdomain"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var ErrBadCredentials = errors.New("authenticate: Bad credentials")

type Authenticator struct{}

func NewID() (string, error) { return gonanoid.New(32) }

func (a *Authenticator) Authenticate(
	ctx context.Context,
	username string,
	password authdomain.Password,
) (account Account, err error) {
	if username == "valid-user@example.com" {
		account = Account{}
	} else {
		err = ErrBadCredentials
	}
	return
}

func New() *Authenticator { return &Authenticator{} }

func init() {
	gob.Register(AccountID(""))
}
