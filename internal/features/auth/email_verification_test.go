package auth_test

import (
	"encoding/json"
	"fmt"
	"harmony/internal/domain"
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
	req, err := http.NewRequest("DELETE", "http://localhost:8025/api/v1/messages", nil)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()

	id := domain.NewID()
	messageID := fmt.Sprintf("<%s@%s>", id, host)
	sendMessage(messageID)
	sendMessage(messageID)
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
	expect(msgResp.Messages).To(gomega.ContainElement(HaveHeader("Message-ID", messageID)))
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

func sendMessage(messageID string) {
	receiver := "user@harmony.example.com"
	firstName := "John"
	code := "123456"
	bodyLines := []string{
		fmt.Sprintf(`Hi %s, Welcome to Harmony`, firstName),
		"",
		"Before you can use the system, you need to verify that you own this email",
		fmt.Sprintf("address. Use the following validation code: %s", code),
		"",
		"You browser you used for creating your account should already be ready to accept",
		"the code. You can also navigate to the following address:",
		"",
		fmt.Sprintf("http://localhost:7331/auth/validate-email?email=%s", receiver),
		"",
		"The Harmony Team.",
	}
	body := strings.Join(bodyLines, "\r\n")
	msg := []byte("To: " + receiver + "\r\n" +
		"Subject: Welcome to Harmony. Please validate your email address.\r\n" +
		"From: info@harmony.example.com\r\n" +
		"To: user@harmony.example.com\r\n" +
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
