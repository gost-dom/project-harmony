package mailhog

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

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

func GetAll() ([]MailhogMessage, error) {
	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var msgResp MailhogGetMessagesResp
	err = json.Unmarshal(data, &msgResp)
	return msgResp.Messages, err
}
