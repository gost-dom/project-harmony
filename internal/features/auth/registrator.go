package auth

import (
	"context"
	"errors"
	"harmony/internal/features/auth/authdomain"
	domain "harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"net/mail"
)

var ErrInvalidInput = errors.New("Invalid input")

type AccountUseCaseResult = UseCaseResult[domain.PasswordAuthentication]

type AccountRepository interface {
	Insert(context.Context, AccountUseCaseResult) error
}

type RegistratorInput struct {
	Email       mail.Address
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
	account := domain.PasswordAuthentication{
		Account: domain.Account{
			ID:          domain.AccountID(authdomain.NewID()),
			Email:       domain.NewUnvalidatedEmail(input.Email),
			Name:        input.Name,
			DisplayName: input.DisplayName,
		},
		PasswordHash: hash,
	}

	res := *NewResult(account)
	res.AddEvent(domain.AccountRegistered{AccountID: account.ID})
	res.AddEvent(domain.EmailValidationRequest{
		AccountID:  account.ID,
		Code:       account.Email.Challenge.Code,
		ValidUntil: account.Email.Challenge.NotAfter,
	})
	return r.Repository.Insert(ctx, res)
}
