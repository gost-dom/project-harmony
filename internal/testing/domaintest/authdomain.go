package domaintest

import "harmony/internal/features/auth/authdomain"

func InitAuthenticatedAccount() authdomain.AuthenticatedAccount {
	return authdomain.AuthenticatedAccount{Account: &authdomain.Account{}}
}
