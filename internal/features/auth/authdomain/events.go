package authdomain

import (
	"encoding/json"
	"harmony/internal/domain"
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

func CreateValidationRequestEvent(account Account) domain.Event {
	return domain.NewDomainEvent(EmailValidationRequest{
		AccountID:  account.ID,
		Code:       account.Email.Challenge.Code,
		ValidUntil: account.Email.Challenge.NotAfter,
	})
}

func CreateAccountRegisteredEvent(account Account) domain.Event {
	return domain.NewDomainEvent(AccountRegistered{AccountID: account.ID})
}

func UnmarshalAuthEvent(data []byte) (domain.EventBody, error) {
	var res EmailValidationRequest
	return res, json.Unmarshal(data, &res)
}

func init() {
	domain.RegisterUnmarshaller(domain.UnmarshallerFunc(UnmarshalAuthEvent))
}
