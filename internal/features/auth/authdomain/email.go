package authdomain

import (
	"errors"
	"net/mail"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// ErrBadEmailValidationCode is returned when an incorrect email validation
// challenge response was provided when proving ownership of an email address.
var ErrBadEmailValidationCode = errors.New("Email.ChallenreResponse: Bad email validation code")

type EmailValidationCode string

func NewValidationCode() EmailValidationCode {
	return EmailValidationCode(gonanoid.MustGenerate("0123456789", 6))
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

func (e Email) String() string { return e.address }

// ChallengeResponse processes a challenge response and returns a validated Email if
// the response is correct. Returns a zero-value Email and
// ErrBadEmailValidationCode err value if the challenge response is wrong.
func (e Email) ChallengeResponse(response EmailValidationCode) (Email, error) {
	if e.Validated {
		return e, nil
	}
	if e.Challenge == nil || e.Challenge.Code != response || e.Challenge.Expired() {
		return e, ErrBadEmailValidationCode
	}
	res := e
	res.Validated = true
	res.Challenge = nil
	return res, nil
}

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
