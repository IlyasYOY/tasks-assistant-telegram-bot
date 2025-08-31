package handler_test

import (
	"testing"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestStartHandler_Handle_Success(t *testing.T) {
	mockSender, startHandler := NewStartHandlerWithMocks(t)
	mockSender.
		EXPECT().
		Send(int64(123), "👋 Hello! I'm *Tasks Assistant* – I can help you manage your tasks using AI.\n\n"+
			"Just send me any text and I’ll treat it as a new task. I’ll always reply with the current task list.").
		Return(nil)

	err := startHandler.Handle(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
	})

	require.NoError(t, err)
}

func TestStartHandler_Handle_SendError(t *testing.T) {
	mockSender, startHandler := NewStartHandlerWithMocks(t)
	mockSender.
		EXPECT().
		Send(int64(115), "👋 Hello! I'm *Tasks Assistant* – I can help you manage your tasks using AI.\n\n"+
			"Just send me any text and I’ll treat it as a new task. I’ll always reply with the current task list.").
		Return(assert.AnError)

	gotErr := startHandler.Handle(&tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 115},
	})

	require.ErrorIs(t, gotErr, assert.AnError)
}

func NewStartHandlerWithMocks(
	t *testing.T,
) (*MockMessageSender, *handler.StartHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSender := NewMockMessageSender(ctrl)

	return mockSender, handler.NewStartHandler(mockSender)
}
