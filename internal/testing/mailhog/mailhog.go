package mailhog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DeleteAll removes all messages from the local MailHog server.
// Returns an error if the HTTP request fails or the server responds with a status code of 400 or higher.
func DeleteAll() error {
	req, err := http.NewRequest("DELETE", "http://localhost:8025/api/v1/messages", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("mailhog: DeleteAll: bad status: %d", resp.StatusCode)
	}
	return nil
}

type MailhogMessageContent struct {
	Headers url.Values
}

type MailhogMessage struct {
	ID      string `json:"id"`
	Content MailhogMessageContent
}

type MailhogGetMessagesResp struct {
	Messages []MailhogMessage `json:"items"`
}

// GetAll retrieves all email messages from the local MailHog server.
// Returns a slice of MailhogMessage and an error if the request fails, the response status is not 200, or the response cannot be parsed.
func GetAll() ([]MailhogMessage, error) {
	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("mailhog: GetAll: unexpected status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var msgResp MailhogGetMessagesResp
	err = json.Unmarshal(data, &msgResp)
	return msgResp.Messages, err
}
