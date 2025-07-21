package authdomain

import (
	"harmony/internal/core"
	"reflect"
	"time"
)

// EmailValidationRequest is a domain event published when an email has been
// registered, and the owner needs to provide a challenge response to prove
// ownership of the email address.
type EmailValidationRequest struct {
	AccountID  `json:"account_id"`
	Code       EmailValidationCode `json:"validation_code"`
	ValidUntil time.Time           `json:"valid_until"`
}

// AccountRegistered is a domain event published when a new account has been
// created.
type AccountRegistered struct {
	AccountID
}

func CreateAccountRegisteredEvent(account Account) core.DomainEvent {
	return core.NewDomainEvent(AccountRegistered{AccountID: account.ID})
}

func init() {
	core.RegisterEventType(
		reflect.TypeFor[EmailValidationRequest](),
		"auth.EmailValidationRequest",
	)
	core.RegisterEventType(reflect.TypeFor[AccountRegistered](), "auth.AccountRegistered")
}
