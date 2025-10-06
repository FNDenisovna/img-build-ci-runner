package telegram

const telegramAPIURL = "https://api.telegram.org/bot"

type Message struct {
	Text   string `json:"text"`
	ChatID string `json:"chat_id"`
}
