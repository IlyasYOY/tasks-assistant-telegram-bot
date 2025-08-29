package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/config"
	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/handler/prompts"
	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/task"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/openai/openai-go"
)

type NewTaskHandler struct {
	sender   MessageSender
	store    Store
	cfg      *config.Config
	aiClient *openai.Client
}

type Store interface {
	Get(userID int64) (task.TasksText, error)
	Set(userID int64, t task.TasksText) error
}

func NewNewTaskHandler(
	sender MessageSender,
	store Store,
	cfg *config.Config,
	aiClient *openai.Client,
) *NewTaskHandler {
	return &NewTaskHandler{
		sender:   sender,
		store:    store,
		cfg:      cfg,
		aiClient: aiClient,
	}
}

var ErrNoChoices = errors.New("no choices")

func (h *NewTaskHandler) Handle(msg *tgbotapi.Message) error {
	userInput, err := h.validateUserInput(msg)
	if err != nil {
		return err
	}

	userPromptString, err := h.buildUserPromptString(msg, userInput)
	if err != nil {
		return err
	}

	newTasks, err := h.runAI(userPromptString, msg)
	if err != nil {
		return err
	}

	setErr := h.store.Set(msg.From.ID, task.TasksText(newTasks))
	if setErr != nil {
		return fmt.Errorf("failed to set task: %w", setErr)
	}

	if err := h.sender.Send(msg.Chat.ID, fmt.Sprintf("# Tasks\n\n%s ", newTasks)); err != nil {
		return fmt.Errorf("failed to send tasks message: %w", err)
	}

	return nil
}

func (h *NewTaskHandler) runAI(
	userPromptString string,
	msg *tgbotapi.Message,
) (string, error) {
	resp, err := h.aiClient.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: h.cfg.OpenAPIModel,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(prompts.SystemPrompt()),
				openai.UserMessage(userPromptString),
			},
		},
	)
	if err != nil {
		sendErr := h.sender.Send(
			msg.Chat.ID,
			fmt.Sprintf("❌ Failed to send message: %q", err.Error()),
		)
		if sendErr != nil {
			return "", fmt.Errorf("failed to send message: %w", sendErr)
		}

		return "", fmt.Errorf("AI request failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		if sendErr := h.sender.Send(msg.Chat.ID, "❌ AI returned no choices."); sendErr != nil {
			return "", fmt.Errorf("failed to send message: %w", sendErr)
		}
		return "", ErrNoChoices
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func (h *NewTaskHandler) buildUserPromptString(
	msg *tgbotapi.Message,
	userInput string,
) (string, error) {
	var userPrompt strings.Builder
	currentTasks, err := h.store.Get(msg.From.ID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve from store: %w", err)
	}
	if currentTasks != "" {
		userPrompt.WriteString(string(currentTasks))
		userPrompt.WriteString("\n---\n")
	}
	userPrompt.WriteString(userInput)
	userPromptString := userPrompt.String()
	return userPromptString, nil
}

func (h *NewTaskHandler) validateUserInput(
	msg *tgbotapi.Message,
) (string, error) {
	userInput := strings.TrimSpace(msg.Text)
	if userInput == "" {
		if err := h.sender.Send(
			msg.Chat.ID,
			"❗️ I received an empty message – please type something.",
		); err != nil {
			return "", fmt.Errorf("sending empty input warning: %w", err)
		}
		return "", nil
	}
	return userInput, nil
}
