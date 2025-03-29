package auth

import (
	"context"
	"errors"
	. "harmony/internal/features/auth/authdomain"
)

type AccountUseCaseResult struct {
	UseCaseResult[Account]
	PasswordAuthentication
}

type AccountRepository interface {
	Insert(context.Context, AccountUseCaseResult) error
}

type RegistratorInput struct {
	Email       string
	Password    Password
	Name        string
	DisplayName string
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
		ID:          AccountID(id),
		Email:       input.Email,
		Name:        input.Name,
		DisplayName: input.DisplayName,
	}
	res := AccountUseCaseResult{
		*NewResult(account),
		PasswordAuthentication{
			AccountID:    account.ID,
			PasswordHash: hash,
		},
	}
	res.AddEvent(AccountRegistered{AccountID: account.ID})
	return r.Repository.Insert(ctx, res)
}
