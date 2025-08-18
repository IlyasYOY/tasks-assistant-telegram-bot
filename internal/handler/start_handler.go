package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	bh *BotHelper
}

func NewStartHandler(bh *BotHelper) *StartHandler {
	return &StartHandler{bh: bh}
}

func (h *StartHandler) Handle(msg *tgbotapi.Message) error {
	reply := "👋 Hello! I'm *Tasks Assistant* – I can help you manage your tasks using AI.\n\n" +
		"Just send me any text and I’ll treat it as a new task. I’ll always reply with the current task list."
	if err := h.bh.Send(msg.Chat.ID, reply); err != nil {
		return fmt.Errorf("sending start reply: %w", err)
	}
	return nil
}
