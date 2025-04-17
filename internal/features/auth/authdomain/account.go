package authdomain

import (
	"errors"
	"harmony/internal/domain"
	"harmony/internal/features/auth/authdomain/password"
)

// ErrAccountEmailNotValidated is returned when an action requires the account
// to be valid. E.g, email address ownership must be verified before the user
// can successfully authenticate.
var ErrAccountEmailNotValidated = errors.New("Email address not validated")

type AccountID string

var NewID = domain.NewID

type Account struct {
	ID          AccountID
	Email       Email
	Name        string
	DisplayName string
}

// ValidateEmail is the email "challenge response" for the email validation
// code.
func (a *Account) ValidateEmail(code EmailValidationCode) (err error) {
	a.Email, err = a.Email.ChallengeResponse(code)
	return
}

// PasswordAuthentication represents an account and it's associated password.
// This type is introduced for two purposes
//
// - Decouple authentication from user account.
// - Security
//
// While password authentication is the only supported method in the first
// prototype, other types could be supported, e.g., google, facebook, github as
// external IDPs; as well as passkey (which could have multiple instances for
// the same user account).
//
// This also reduces the risk of security related issues in code, as passwords
// are only processed during login, authentication, and changing passwords. Once
// the user is logged in, the types being used don't contain password
// information anymore.
type PasswordAuthentication struct {
	Account
	password.PasswordHash
}

// Authenticated tells the account that authentication has been successful.
//
// It has been left for the Account itself to verify that the account itself is
// in a valid state. While different authentication mechanisms can only verify
// that the user has succeeded specific challenges, that doesn't prove that the
// account permits being logged into at all.
func (a *Account) Authenticated() (AuthenticatedAccount, error) {
	var res AuthenticatedAccount
	if !a.Email.Validated {
		return res, ErrAccountEmailNotValidated
	}
	res.Account = a
	return res, nil
}

/* -------- AuthenticatedAccount -------- */

// AuthenticatedAccount represents an Account that has succeded an
// authentication flow. Code that needs to check who is performing an operation
// can depend on this type.
//
// At the moment this type merely indicates that an authentication check has
// succeeded. But it could hold information regarding which kind of
// authentication mechanism was used, e.g., password, passkey. Was 2FA used,
// etc. Or the a returning user with "remember me" enabled.
type AuthenticatedAccount struct{ *Account }
