package domaintest

import "harmony/internal/features/auth/authdomain"

// InitAccount creates and returns a valid minimal Account for test scenarios
// that requires a valid account, but details are irrelevant.
func InitAccount() authdomain.Account {
	return authdomain.Account{}
}

// InitAuthenticatedAccount creates and returns an AuthenticatedAccount with a
// minimal account for use in test scenarios where an authenticated account is
// required, but the specific user details are irrelevant.
func InitAuthenticatedAccount() authdomain.AuthenticatedAccount {
	acc := InitAccount()
	return authdomain.AuthenticatedAccount{Account: &acc}
}
