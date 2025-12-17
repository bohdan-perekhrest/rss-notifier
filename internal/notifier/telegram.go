package notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
)

type Telegram struct {
	chatID string
	token string
	client *http.Client
}

func NewTelegram(client *http.Client, chatID, token string) *Telegram {
	return &Telegram{
		chatID: chatID,
		token:  token,
		client: client,
	}
}

func (t *Telegram) Send(message string) error {
	payload := map[string]interface{}{
		"chat_id":                  t.chatID,
		"text":                     message,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
		"disable_notification":     true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send a request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}
