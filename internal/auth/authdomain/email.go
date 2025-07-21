package authdomain

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// ErrBadEmailChallengeResponse is returned when an incorrect email validation
// challenge response was provided when proving ownership of an email address.
var ErrBadEmailChallengeResponse = errors.New("authdomain: bad email challenge response")

var ErrEmailChallengeExpired = errors.New(
	"authdomain: email challenge response has expired",
)

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
	Address   mail.Address
	Validated bool
	Challenge *EmailChallenge
}

// Equals returns true of the two emails have the same address.
func (e Email) Equals(address string) bool {
	return strings.EqualFold(e.Address.Address, address) && address != ""
}

func (e Email) String() string { return e.Address.Address }

// ChallengeResponse processes a challenge response and returns a validated Email if
// the challenge succeeds; otherwise the same email is returned with one of
// these errors:
//
//   - [ErrBadEmailValidationCode] if the validation code was wrong
//   - [ErrEmailChallengeExpired] if the validation code has expired
func (e Email) ChallengeResponse(response EmailValidationCode) (Email, error) {
	if e.Validated {
		return e, nil
	}
	if e.Challenge != nil && e.Challenge.Code == response {
		if e.Challenge.Expired() {
			return e, ErrEmailChallengeExpired
		}
		res := e
		res.Validated = true
		res.Challenge = nil
		return res, nil
	}
	return e, ErrBadEmailChallengeResponse
}

func (e *Email) NewChallenge() EmailChallenge {
	challenge := EmailChallenge{
		Code:     NewValidationCode(),
		NotAfter: time.Now().Add(15 * time.Minute).UTC(),
	}
	e.Challenge = &challenge
	return challenge
}

func NewUnvalidatedEmail(address mail.Address) Email {
	return Email{Address: address}
}
