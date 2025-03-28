package auth

import (
	"context"
	"errors"
	. "harmony/internal/features/auth/authdomain"
)

type AccountUseCaseResult = UseCaseResult[Account, AccountID]

type AccountRepository interface {
	Insert(context.Context, AccountUseCaseResult) error
}

type RegistratorInput struct {
	Email    string
	Password Password
}

type Registrator struct {
	Repository AccountRepository
}

func (r Registrator) Register(ctx context.Context, input RegistratorInput) error {
	id, err1 := NewID()
	hash, err2 := NewHash(input.Password)
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	account := Account{
		Id:       AccountID(id),
		Email:    input.Email,
		Password: hash,
	}
	res := NewResult(account)
	res.AddEvent(AccountRegistered{AccountID: account.Id})
	return r.Repository.Insert(ctx, *res)
}
