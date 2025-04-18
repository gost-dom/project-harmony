package auth_test

import (
	"encoding/json"
	"fmt"
	"harmony/internal/domain"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/testing/domaintest"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/stretchr/testify/assert"
)

const host = "harmony.example.com"

type mailhogMessageContent struct {
	Headers url.Values
}

type mailhogMessage struct {
	ID      string `json:"id"`
	Content mailhogMessageContent
}

type mailhogGetMessagesResp struct {
	Messages []mailhogMessage `json:"items"`
}

func TestSendEmailValidationChallenge(t *testing.T) {
	acc := domaintest.InitAccount()
	acc.StartEmailValidationChallenge()
	assert.False(
		t,
		acc.Validated(),
		"guard: account should be an invalidated account for this test",
	)

	req, err := http.NewRequest("DELETE", "http://localhost:8025/api/v1/messages", nil)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()

	id := domain.NewID()
	messageID := fmt.Sprintf("<%s@%s>", id, host)
	sendMessage(messageID, acc)
	resp, err = http.Get("http://localhost:8025/api/v2/messages")
	assert.NoError(t, err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var msgResp mailhogGetMessagesResp
	json.Unmarshal(data, &msgResp)
	t.Log("\n\nDATA:\n", string(data))
	t.Logf("\n\nBody: %+v\n", msgResp)
	expect := gomega.NewWithT(t).Expect
	expect(
		msgResp.Messages,
	).To(gomega.ContainElement(HaveHeader("To", acc.Email.Address.String())))
}

func HaveHeader(key, value string) types.GomegaMatcher {
	return gomega.WithTransform(func(m mailhogMessage) ([]string, error) {
		res, ok := m.Content.Headers[key]
		if !ok {
			return nil, fmt.Errorf("Message did not contain header: %s", key)
		}
		return res, nil
	}, gomega.ContainElement(gomega.Equal(value)))
}

func sendMessage(messageID string, acc authdomain.Account) {
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
	msg := []byte("To: " + receiver + "\r\n" +
		"Subject: Welcome to Harmony. Please validate your email address.\r\n" +
		"From: info@harmony.example.com\r\n" +
		fmt.Sprintf("To: %s\r\n", receiver.String()) +
		fmt.Sprintf("Message-ID: %s\r\n", messageID) +
		"\r\n" +
		body +
		"\r\n")

	// Send the email
	err := smtp.SendMail(
		"localhost:1025",
		nil,
		"info@harmony.example.com",
		[]string{"user@harmony.example.com"},
		msg,
	)
	if err != nil {
		log.Fatal(err)
	}
}
