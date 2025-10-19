package domaintest

import (
	"fmt"
	"harmony/internal/auth/domain"
	"harmony/internal/auth/domain/password"
	"net/mail"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func NewAddress() string {
	return fmt.Sprintf("%s@example.com", gonanoid.MustGenerate("abcdefghjklmnipqrstuvwxyz", 20))
}

func InitEmail() domain.Email {
	addr, err := mail.ParseAddress(NewAddress())
	must("domaintest: InitEmail", err)
	res := domain.NewUnvalidatedEmail(*addr)
	return res
}

func must(prefix string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", prefix, err))
	}
}

func WithDisplayName(name string) InitAccountOption {
	return func(acc *domain.Account) {
		acc.DisplayName = name
	}
}

func WithEmail(email string) InitAccountOption {
	addr, err := mail.ParseAddress(email)
	must("domaintest: WithEmail", err)
	return WithEmailAddress(addr)
}

func WithEmailAddress(addr *mail.Address) InitAccountOption {
	return func(acc *domain.Account) {
		em := domain.NewUnvalidatedEmail(*addr)
		acc.Email = em
	}
}

func WithName(name string) InitAccountOption {
	return func(acc *domain.Account) {
		acc.Name = name
	}
}

// WithEmailValidation makes sure that the email address has been validated, so
// the account can be treated as "authenticated".
//
// See also: [domain.Account.ValidateEmail]
func WithEmailValidation() InitAccountOption {
	return func(acc *domain.Account) {
		if !acc.Email.Validated {
			c := acc.Email.NewChallenge()
			if err := acc.ValidateEmail(c.Code); err != nil {
				panic("WithValidatedEmail: error validating email: " + err.Error())
			}
		}
	}
}

type InitAccountOption = func(*domain.Account)

// InitAccount creates and returns a valid minimal Account for test scenarios
// that requires a valid account, but details are irrelevant.
func InitAccount(opts ...InitAccountOption) domain.Account {
	result := domain.Account{
		ID:    domain.AccountID(domain.NewID()),
		Email: InitEmail(),
	}
	for _, o := range opts {
		o(&result)
	}
	return result
}

// InitAuthenticatedAccount creates and returns an AuthenticatedAccount with a
// minimal account for use in test scenarios where an authenticated account is
// required, but the specific user details are irrelevant.
func InitAuthenticatedAccount(opts ...InitAccountOption) domain.AuthenticatedAccount {
	acc := InitAccount(append(opts, WithEmailValidation())...)
	res, err := acc.Authenticated()
	if err != nil {
		panic("InitAuthenticatedAccount: failed creating authenticated account: " + err.Error())
	}
	return res
}

type InitPasswordOption = func(*domain.PasswordAuthentication)

func WithPassword(pw string) InitPasswordOption {
	return func(ac *domain.PasswordAuthentication) {
		var err error
		if ac.PasswordHash, err = password.Parse(pw).Hash(); err != nil {
			panic(fmt.Sprintf("WithPassword: hashing failed: %v", err))
		}
	}
}

// InitPasswordAuthAccount creates a new [domain.PasswordAuthentication]
// entity. The options must be either [InitPasswordOption] or
// [InitAccountOption]. The function panics if any of the options are not a
// compatible type.
func InitPasswordAuthAccount(opts ...any) domain.PasswordAuthentication {
	res := domain.PasswordAuthentication{Account: InitAccount()}
	for _, o := range opts {
		switch t := o.(type) {
		case InitPasswordOption:
			t(&res)
		case InitAccountOption:
			t(&res.Account)
		default:
			panic(fmt.Sprintf("InitPasswordAuthAccount: invalid option type: %T", o))
		}
	}
	return res
}
