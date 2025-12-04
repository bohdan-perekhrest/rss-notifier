package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type SendMessageRequest struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
}

var (
	botToken string
	chatID   string
)

func init() {
	godotenv.Load()
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID = os.Getenv("TELEGRAM_CHAT_ID")
}

func SendMessage(text string) error {
	data := &SendMessageRequest{
		ChatID: chatID,
		Text: text,
		ParseMode: "HTML",
		DisableWebPagePreview: false,
		DisableNotification: true,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.telegram.org/bot" + botToken + "/sendMessage",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: status %d", resp.StatusCode)
	}
	return nil
}
