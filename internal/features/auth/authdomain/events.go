package authdomain

import (
	"encoding/json"
	"harmony/internal/domain"
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

func CreateAccountRegisteredEvent(account Account) domain.Event {
	return domain.NewDomainEvent(AccountRegistered{AccountID: account.ID})
}

func UnmarshalAuthEvent(data []byte) (domain.EventBody, error) {
	var res EmailValidationRequest
	return res, json.Unmarshal(data, &res)
}

func init() {
	domain.RegisterEventType(
		reflect.TypeFor[EmailValidationRequest](),
		"auth.EmailValidationRequest",
	)
	domain.RegisterEventType(reflect.TypeFor[AccountRegistered](), "auth.AccountRegistered")
}
