package auth

import (
	"context"
	"errors"
	"harmony/internal/features/auth/authdomain"
	domain "harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
)

var ErrInvalidInput = errors.New("Invalid input")

type AccountUseCaseResult = UseCaseResult[domain.PasswordAuthentication]

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

// Register attempts to create a new user account with password-based
// authentication.
func (r Registrator) Register(ctx context.Context, input RegistratorInput) error {
	hash, err := input.Password.Hash()
	if err != nil {
		return err
	}
	email, err := domain.NewUnvalidatedEmail(input.Email)
	if err != nil {
		return ErrInvalidInput
	}
	account := domain.PasswordAuthentication{
		Account: domain.Account{
			ID:          domain.AccountID(authdomain.NewID()),
			Email:       email,
			Name:        input.Name,
			DisplayName: input.DisplayName,
		},
		PasswordHash: hash,
	}

	res := *NewResult(account)
	res.AddEvent(domain.AccountRegistered{AccountID: account.ID})
	res.AddEvent(domain.EmailValidationRequest{
		AccountID:  account.ID,
		Code:       email.Challenge.Code,
		ValidUntil: email.Challenge.NotAfter,
	})
	return r.Repository.Insert(ctx, res)
}
