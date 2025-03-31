package auth

import (
	"context"
	"errors"
	. "harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
)

var ErrInvalidInput = errors.New("Invalid input")
var ErrAccountEmailNotValidated = errors.New("Email address not validated")

type InsertAccount struct {
	Account
	PasswordAuthentication
}

type AccountUseCaseResult = UseCaseResult[InsertAccount]

type AccountRepository interface {
	Insert(context.Context, AccountUseCaseResult) error
}

type RegistratorInput struct {
	Email       string
	Password    password.Password
	Name        string
	DisplayName string
}

type Registrator struct {
	Repository AccountRepository
}

// Register attempts to create a new user account.
func (r Registrator) Register(ctx context.Context, input RegistratorInput) error {
	id, err1 := NewID()
	hash, err2 := input.Password.Hash()
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	email, err := NewUnvalidatedEmail(input.Email)
	if err != nil {
		return ErrInvalidInput
	}
	account := Account{
		ID:          AccountID(id),
		Email:       email,
		Name:        input.Name,
		DisplayName: input.DisplayName,
	}
	res := *NewResult(InsertAccount{account,
		PasswordAuthentication{
			AccountID:    account.ID,
			PasswordHash: hash,
		}})

	res.AddEvent(AccountRegistered{AccountID: account.ID})
	res.AddEvent(EmailValidationRequest{
		AccountID:  account.ID,
		Code:       email.Challenge.Code,
		ValidUntil: email.Challenge.NotAfter,
	})
	return r.Repository.Insert(ctx, res)
}
