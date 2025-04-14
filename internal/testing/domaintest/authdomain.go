package domaintest

import (
	"fmt"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func NewAddress() string { return fmt.Sprintf("%s@example.com", gonanoid.Must()) }

func InitEmail() authdomain.Email {
	res, err := authdomain.NewUnvalidatedEmail(NewAddress())
	if err != nil {
		panic(err)
	}
	return res
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func WithEmail(email string) InitAccountOption {
	return func(acc *authdomain.Account) {
		em, err := authdomain.NewUnvalidatedEmail(email)
		must(err)
		acc.Email = em
	}
}

type InitAccountOption = func(*authdomain.Account)

// InitAccount creates and returns a valid minimal Account for test scenarios
// that requires a valid account, but details are irrelevant.
func InitAccount(opts ...InitAccountOption) authdomain.Account {
	result := authdomain.Account{
		ID:    authdomain.AccountID(authdomain.NewID()),
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
func InitAuthenticatedAccount(opts ...InitAccountOption) authdomain.AuthenticatedAccount {
	acc := InitAccount(opts...)
	return authdomain.AuthenticatedAccount{Account: &acc}
}

type InitPasswordOption = func(*authdomain.PasswordAuthentication)

func WithPassword(pw string) InitPasswordOption {
	return func(ac *authdomain.PasswordAuthentication) {
		var err error
		if ac.PasswordHash, err = password.Parse(pw).Hash(); err != nil {
			panic(fmt.Sprintf("WithPassword: hashing failed: %v", err))
		}
	}
}

// InitPasswordAuthAccount creates a new [authdomain.PasswordAuthentication]
// entity. The options must be either [InitPasswordOption] or
// [InitAccountOption]. The function panics if any of the options are not a
// compatible type.
func InitPasswordAuthAccount(opts ...any) authdomain.PasswordAuthentication {
	res := authdomain.PasswordAuthentication{Account: InitAccount()}
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
