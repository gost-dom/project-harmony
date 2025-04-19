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
	Email            *mail.Address
	Password         password.Password
	Name             string
	DisplayName      string
	NewsletterSignup bool
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
	if r.Repository == nil {
		return errors.New("TODO: Get repo working")
	}
	if err != nil {
		return err
	}
	account := domain.PasswordAuthentication{
		Account: domain.Account{
			ID:          domain.AccountID(authdomain.NewID()),
			Email:       domain.NewUnvalidatedEmail(*input.Email),
			Name:        input.Name,
			DisplayName: input.DisplayName,
		},
		PasswordHash: hash,
	}

	res := UseCaseOfEntity(account)
	res.AddEvent(domain.CreateAccountRegisteredEvent(account.Account))
	res.AddEvent(res.Entity.StartEmailValidationChallenge())
	return r.Repository.Insert(ctx, res)
}
