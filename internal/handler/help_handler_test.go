package handler_test

import (
	"testing"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestHelpHandler_Handle_Success(t *testing.T) {
	mockSender, helpHandler := NewHelpHandlerWithMocks(t)

	mockSender.
		EXPECT().
		Send(int64(123), "*How to use the bot:*\n"+
			"• Send any plain text – it will be saved as a new task.\n"+
			"• The bot will immediately reply with the complete list of your tasks.\n"+
			"• /start – greeting message\n"+
			"• /help – this help text").
		Return(nil)

	err := helpHandler.Handle(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
	})

	require.NoError(t, err)
}

func TestHelpHandler_Handle_SendError(t *testing.T) {
	mockSender, helpHandler := NewHelpHandlerWithMocks(t)

	mockSender.
		EXPECT().
		Send(int64(456), "*How to use the bot:*\n"+
			"• Send any plain text – it will be saved as a new task.\n"+
			"• The bot will immediately reply with the complete list of your tasks.\n"+
			"• /start – greeting message\n"+
			"• /help – this help text").
		Return(assert.AnError)

	gotErr := helpHandler.Handle(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 456},
	})

	require.ErrorIs(t, gotErr, assert.AnError)
}

func NewHelpHandlerWithMocks(
	t *testing.T,
) (*MockMessageSender, *handler.HelpHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSender := NewMockMessageSender(ctrl)

	return mockSender, handler.NewHelpHandler(mockSender)
}
