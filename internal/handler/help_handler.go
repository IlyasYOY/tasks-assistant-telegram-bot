package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct {
	bh *BotHelper
}

func NewHelpHandler(bh *BotHelper) *HelpHandler {
	return &HelpHandler{bh: bh}
}

func (h *HelpHandler) Handle(msg *tgbotapi.Message) error {
	helpText := "*How to use the bot:*\n" +
		"• Send any plain text – it will be saved as a new task.\n" +
		"• The bot will immediately reply with the complete list of your tasks.\n" +
		"• /start – greeting message\n" +
		"• /help – this help text"
	if err := h.bh.Send(msg.Chat.ID, helpText); err != nil {
		return fmt.Errorf("sending help reply: %w", err)
	}
	return nil
}
