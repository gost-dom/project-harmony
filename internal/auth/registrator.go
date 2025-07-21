package auth

import (
	"context"
	"errors"
	"harmony/internal/auth/authdomain"
	domain "harmony/internal/auth/authdomain"
	"harmony/internal/auth/authdomain/password"
	"harmony/internal/core"
	"net/mail"
)

var ErrInvalidInput = errors.New("Invalid input")

type AccountUseCaseResult = core.UseCaseResult[domain.PasswordAuthentication]

type AccountInserter interface {
	Insert(context.Context, AccountUseCaseResult) (domain.PasswordAuthentication, error)
}

type ValidateEmailInput struct {
	Email *mail.Address
	Code  authdomain.EmailValidationCode
}

type RegistratorInput struct {
	Email            *mail.Address
	Password         password.Password
	Name             string
	DisplayName      string
	NewsletterSignup bool
}

type Registrator struct {
	Repository AccountInserter
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
	account := domain.PasswordAuthentication{
		Account: domain.Account{
			ID:          domain.AccountID(authdomain.NewID()),
			Email:       domain.NewUnvalidatedEmail(*input.Email),
			Name:        input.Name,
			DisplayName: input.DisplayName,
		},
		PasswordHash: hash,
	}

	res := core.UseCaseOfEntity(account)
	res.AddEvent(domain.CreateAccountRegisteredEvent(account.Account))
	res.AddEvent(res.Entity.StartEmailValidationChallenge())
	_, err = r.Repository.Insert(ctx, res)
	return err
}
