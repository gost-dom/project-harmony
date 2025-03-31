package authdomain

import (
	"errors"
	"harmony/internal/features/auth/authdomain/password"
	"net/mail"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
)

var ErrBadEmailValidationCode = errors.New("Bad email validation code")

type ValidationCode string

func NewValidationCode() ValidationCode {
	return ValidationCode(nanoid.MustGenerate("0123456789", 6))
}

// An email "challenge", i.e., a randomly generated code sent to an email
// address that the owner must provide as a "challenge response" to prove
// ownership of the email address.
type EmailChallenge struct {
	Code     ValidationCode
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

func (e Email) Equals(address string) bool {
	return e.address == address && address != ""
}

func (e Email) Validate(code ValidationCode) (res Email, err error) {
	res = e
	if e.Challenge.Code != code || e.Challenge.Expired() {
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
	EmailValidationCode ValidationCode
}

// ValidateEmail is the email "challenge response" for the email validation
// code.
func (a *Account) ValidateEmail(code ValidationCode) (err error) {
	a.Email, err = a.Email.Validate(code)
	return
}

type PasswordAuthentication struct {
	AccountID
	password.PasswordHash
}

type AccountRegistered struct {
	AccountID
}

type EmailValidationRequest struct {
	AccountID
	Code       ValidationCode
	ValidUntil time.Time
}
