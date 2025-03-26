package auth

import (
	"context"
	"encoding/gob"
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var ErrBadCredentials = errors.New("authenticate: Bad credentials")

type Authenticator struct{}

type AccountID string

func NewID() (string, error) { return gonanoid.New(32) }

type Account struct {
	Id    AccountID
	Email string
}

func (a Account) ID() AccountID { return a.Id }

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
	gob.Register(AccountID(""))
}
