package authenticator

import (
	"context"
	"encoding/gob"
	"errors"
)

var ErrBadCredentials = errors.New("authenticate: Bad credentials")

type Authenticator struct{}

type AccountId string

type Account struct{ Id AccountId }

func (a *Authenticator) Authenticate(
	ctx context.Context,
	username string,
	password string,
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
	gob.Register(AccountId(""))
}
