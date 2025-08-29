package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownHandler struct {
	sender MessageSender
}

func NewUnknownHandler(sender MessageSender) *UnknownHandler {
	return &UnknownHandler{sender: sender}
}

func (h *UnknownHandler) Handle(msg *tgbotapi.Message) error {
	reply := fmt.Sprintf(
		"❓ I don't understand the command %q. Use /help to see available commands.",
		msg.Text,
	)
	if err := h.sender.Send(msg.Chat.ID, reply); err != nil {
		return fmt.Errorf("sending unknown command reply: %w", err)
	}
	return nil
}
