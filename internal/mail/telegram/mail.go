package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Mail struct {
	botToken  string
	channelId string
}

func New(botToken, channelId string) *Mail {
	return &Mail{
		botToken:  botToken,
		channelId: channelId,
	}
}

func (m *Mail) SendMessage(text string) error {
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIURL, m.botToken)
	message := Message{ChatID: m.channelId, Text: text}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}
	return nil
}
