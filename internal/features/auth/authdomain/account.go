package authdomain

import (
	"errors"
	"harmony/internal/features/auth/authdomain/password"
	"net/mail"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
)

// ErrBadEmailValidationCode is returned when an incorrect email validation
// challenge response was provided when proving ownership of an email address.
var ErrBadEmailValidationCode = errors.New("Bad email validation code")

// ErrAccountEmailNotValidated is returned when an action requires the account
// to be valid. E.g, email address ownership must be verified before the user
// can successfully authenticate.
var ErrAccountEmailNotValidated = errors.New("Email address not validated")

type EmailValidationCode string

func NewValidationCode() EmailValidationCode {
	return EmailValidationCode(nanoid.MustGenerate("0123456789", 6))
}

// An email "challenge", i.e., a randomly generated code sent to an email
// address that the owner must provide as a "challenge response" to prove
// ownership of the email address.
type EmailChallenge struct {
	Code     EmailValidationCode
	NotAfter time.Time // A deadline for completing the challenge
}

func (c EmailChallenge) Expired() bool { return time.Now().After(c.NotAfter) }

// Email is a value object encapsulating the complexities of email address
// validation through a challenge.
type Email struct {
	address   string
	Validated bool
	Challenge *EmailChallenge
}

// Equals returns true of the two emails have the same address.
func (e Email) Equals(address string) bool {
	return e.address == address && address != ""
}

// ChallengeResponse processes a challenge response and returns a validated Email if
// the response is correct. Returns a zero-value Email and
// ErrBadEmailValidationCode err value if the challenge response is wrong.
func (e Email) ChallengeResponse(response EmailValidationCode) (res Email, err error) {
	res = e
	if e.Challenge.Code != response || e.Challenge.Expired() {
		err = ErrBadEmailValidationCode
	} else {
		res.Validated = true
		res.Challenge = nil
	}
	return
}

func (e Email) String() string { return e.address }

func NewUnvalidatedEmail(address string) (Email, error) {
	_, err := mail.ParseAddress(address)
	email := Email{
		address: address,
		Challenge: &EmailChallenge{
			Code:     NewValidationCode(),
			NotAfter: time.Now().Add(15 * time.Minute),
		},
	}
	return email, err
}

type AccountID string

type Account struct {
	ID                  AccountID
	Email               Email
	Name                string
	DisplayName         string
	EmailValidationCode EmailValidationCode
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

// AccountRegistered is a domain event published when a new account has been
// created.
type AccountRegistered struct {
	AccountID
}

// EmailValidationRequest is a domain event published when an email has been
// registered, and the owner needs to provide a challenge response to prove
// ownership of the email address.
type EmailValidationRequest struct {
	AccountID
	Code       EmailValidationCode
	ValidUntil time.Time
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
// At the moment this type merely indicatest that an authentication chack has
// succeeded. But it could hold information regarding which kind of
// authentication mechanism was used, e.g., password, passkey. Was 2FA used,
// etc. It this a revisit from a user with "remember me" enabled.
type AuthenticatedAccount struct{ *Account }
