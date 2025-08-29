package handler

import (
	"errors"
	"fmt"
	"slices"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	sender MessageSender

	cfg *config.Config

	cmdMap   map[string]CommandHandler
	plain    CommandHandler
	fallback CommandHandler
}

type CommandHandler interface {
	Handle(msg *tgbotapi.Message) error
}

type MessageSender interface {
	Send(chatID int64, text string) error
}

func New(
	cfg *config.Config,
	sender MessageSender,
	startHandler CommandHandler,
	helpHandler CommandHandler,
	newTaskHandler CommandHandler,
	unknownHandler CommandHandler,
) *Handler {
	return &Handler{
		cfg:    cfg,
		sender: sender,
		cmdMap: map[string]CommandHandler{
			"start": startHandler,
			"help":  helpHandler,
		},
		plain:    newTaskHandler,
		fallback: unknownHandler,
	}
}

var ErrUserIsNotAuthorized = errors.New("user is not authorized")

func (h *Handler) HandleUpdate(update *tgbotapi.Update) error {
	msg := update.Message
	if msg.From == nil {
		return nil
	}
	if !slices.Contains(h.cfg.AllowedUserIDs, msg.From.ID) {
		return fmt.Errorf("%w: %d", ErrUserIsNotAuthorized, msg.From.ID)
	}
	if !msg.IsCommand() {
		if err := h.plain.Handle(msg); err != nil {
			return fmt.Errorf(
				"plain‑text handler error for user %d: %w",
				msg.From.ID,
				err,
			)
		}
		return nil
	}

	cmd, ok := h.cmdMap[msg.Command()]
	if !ok {
		if err := h.fallback.Handle(msg); err != nil {
			return fmt.Errorf(
				"fallback handler error for unknown command %q (user %d): %w",
				msg.Text,
				msg.From.ID,
				err,
			)
		}
		return nil
	}

	if err := cmd.Handle(msg); err != nil {
		return fmt.Errorf(
			"command %q handling error for user %d: %w",
			msg.Command(),
			msg.From.ID,
			err,
		)
	}

	return nil
}
