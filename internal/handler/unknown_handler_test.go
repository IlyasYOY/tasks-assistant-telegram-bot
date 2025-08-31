package handler_test

import (
	"testing"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestUnknownHandler_Handle_Success(t *testing.T) {
	mockSender, unknownHandler := NewUnknownHandlerWithMocks(t)
	mockSender.
		EXPECT().
		Send(int64(123), "❓ I don't understand the command \"/unknown\". Use /help to see available commands.").
		Return(nil)

	gotErr := unknownHandler.Handle(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "/unknown",
	})

	require.NoError(t, gotErr)
}

func TestUnknownHandler_Handle_SendError(t *testing.T) {
	mockSender, unknownHandler := NewUnknownHandlerWithMocks(t)
	mockSender.
		EXPECT().
		Send(int64(456), "❓ I don't understand the command \"/oops\". Use /help to see available commands.").
		Return(assert.AnError)

	gotErr := unknownHandler.Handle(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 456},
		Text: "/oops",
	})

	require.ErrorIs(t, gotErr, assert.AnError)
}

func NewUnknownHandlerWithMocks(
	t *testing.T,
) (*MockMessageSender, *handler.UnknownHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSender := NewMockMessageSender(ctrl)

	return mockSender, handler.NewUnknownHandler(mockSender)
}
