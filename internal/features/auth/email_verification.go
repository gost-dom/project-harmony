package auth

import "harmony/internal/features/auth/authdomain"

type AccountLoader interface {
	GetAccount(authdomain.AccountID) authdomain.Account
}

type EmailValidator struct {
	AccountRepository
}

func (v EmailValidator) SendEmailValidationChallenge(authdomain.EmailValidationRequest) {

}
