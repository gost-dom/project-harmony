package domaintest

import (
	"fmt"
	"harmony/internal/features/auth/authdomain"

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
func InitAuthenticatedAccount() authdomain.AuthenticatedAccount {
	acc := InitAccount()
	return authdomain.AuthenticatedAccount{Account: &acc}
}
