package auth

import (
	"context"
	"errors"
	"fmt"
	"harmony/internal/core"
	"harmony/internal/auth/authdomain"
	"net/smtp"
	"strings"
)

const host = "harmony.example.com"

type EmailChallengeRepository interface {
	FindByEmail(context.Context, string) (authdomain.Account, error)
	Update(context.Context, authdomain.Account) (authdomain.Account, error)
}

type EmailChallengeValidator struct {
	Repository EmailChallengeRepository
}

func (a EmailChallengeValidator) Validate(
	ctx context.Context,
	input ValidateEmailInput,
) (res authdomain.AuthenticatedAccount, err error) {
	defer func() {
		if errors.Is(err, authdomain.ErrBadEmailChallengeResponse) {
			err = ErrBadChallengeResponse
		}
	}()

	acc, err := a.Repository.FindByEmail(ctx, input.Email.Address)
	if err == nil {
		err = acc.ValidateEmail(input.Code)
	}
	if err == nil {
		acc, err = a.Repository.Update(ctx, acc)
	}
	if err == nil {
		return acc.Authenticated()
	}
	return
}

type AccountLoader interface {
	Get(context.Context, authdomain.AccountID) (authdomain.Account, error)
}

type EmailValidator struct {
	Repository AccountLoader
}

func NewEmailValidator() *EmailValidator { return &EmailValidator{nil} }

func (v EmailValidator) ProcessDomainEvent(ctx context.Context, event core.DomainEvent) error {
	req, ok := event.Body.(authdomain.EmailValidationRequest)
	if !ok { // Not an event we want to handle
		return nil
	}
	acc, err := v.Repository.Get(ctx, req.AccountID)
	if err == nil {
		err = sendChallengeEmail(string(event.ID), acc)
	}
	if err != nil {
		err = fmt.Errorf("auth: ProcessDomainEvent: %w", err)
	}
	return err
}

func sendChallengeEmail(eventID string, acc authdomain.Account) error {
	messageID := fmt.Sprintf("<%s@%s>", eventID, host)
	receiver := acc.Email.Address // Yeah, net/mail.Address has an Address field
	receiver.Name = acc.Name
	firstName := acc.DisplayName
	code := string(acc.Email.Challenge.Code)

	bodyLines := []string{
		fmt.Sprintf(`Hi %s, Welcome to Harmony`, firstName),
		"",
		"Before you can use the system, you need to verify that you own this email",
		"address. Use the following validation code",
		"",
		"    " + code,
		"",
		"",
		"The browser you used when registering should already be ready to accept",
		"the code. If not, you can also navigate to the following address:",
		"",
		fmt.Sprintf("http://localhost:7331/auth/validate-email?email=%s", receiver.Address),
		"",
		"The Harmony Team.",
	}
	body := strings.Join(bodyLines, "\r\n")
	msg := []byte("To: " + receiver.Address + "\r\n" +
		"Subject: Welcome to Harmony. Please validate your email address.\r\n" +
		"From: info@harmony.example.com\r\n" +
		fmt.Sprintf("To: %s\r\n", receiver.String()) +
		fmt.Sprintf("Message-ID: %s\r\n", messageID) +
		"\r\n" +
		body +
		"\r\n")

	// Send the email
	return smtp.SendMail(
		"localhost:1025",
		nil,
		"info@harmony.example.com",
		[]string{receiver.Address},
		msg,
	)

}
