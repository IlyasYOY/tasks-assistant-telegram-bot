package msgsender

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type messageSender struct {
	bot *tgbotapi.BotAPI
}

func NewMessageSender(bot *tgbotapi.BotAPI) *messageSender {
	return &messageSender{bot: bot}
}

func (h *messageSender) Send(chatID int64, text string) error {
	m := tgbotapi.NewMessage(chatID, text)
	m.ParseMode = tgbotapi.ModeMarkdown
	if _, err := h.bot.Send(m); err != nil {
		return fmt.Errorf("failed to send message to chat %d: %w",
			chatID, err)
	}
	return nil
}
