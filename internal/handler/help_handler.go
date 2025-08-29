package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct {
	sender MessageSender
}

func NewHelpHandler(sender MessageSender) *HelpHandler {
	return &HelpHandler{sender: sender}
}

func (h *HelpHandler) Handle(msg *tgbotapi.Message) error {
	helpText := "*How to use the bot:*\n" +
		"• Send any plain text – it will be saved as a new task.\n" +
		"• The bot will immediately reply with the complete list of your tasks.\n" +
		"• /start – greeting message\n" +
		"• /help – this help text"
	if err := h.sender.Send(msg.Chat.ID, helpText); err != nil {
		return fmt.Errorf("sending help reply: %w", err)
	}
	return nil
}
