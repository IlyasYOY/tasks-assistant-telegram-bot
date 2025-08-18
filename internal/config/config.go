package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	EnvTelegramToken   = "TASKS_ASSISTANT_TG_BOT_TELEGRAM_TOKEN"
	EnvOpenAPIBasePath = "TASKS_ASSISTANT_TG_BOT_OPEN_API_BASE_PATH"
	EnvOpenAPIKey      = "TASKS_ASSISTANT_TG_BOT_OPEN_API_API_KEY"
	EnvOpenAPIModel    = "TASKS_ASSISTANT_TG_BOT_OPEN_API_MODEL"
	EnvAllowedUserIDs  = "TASKS_ASSISTANT_TG_BOT_ALLOWED_USER_IDS"
	EnvSQLDSN          = "TASKS_ASSISTANT_TG_BOT_SQL_DSN"
)

type Config struct {
	TelegramToken   string
	OpenAPIBasePath string
	OpenAPIKey      string
	OpenAPIModel    string
	SQLDSN          string
	AllowedUserIDs  []int64
}

var (
	ErrMissingTelegramToken = errors.New("missing env variable")
	ErrMissingOpenAPIKey    = errors.New("missing OpenAPI key")
	ErrMissingOpenAPIModel  = errors.New("missing OpenAPI model")
)

func New() (*Config, error) {
	sqlDSN := os.Getenv(EnvSQLDSN)
	if sqlDSN == "" {
		sqlDSN = "file::memory:?cache=shared"
	}

	cfg := &Config{
		TelegramToken:   os.Getenv(EnvTelegramToken),
		OpenAPIBasePath: os.Getenv(EnvOpenAPIBasePath),
		OpenAPIKey:      os.Getenv(EnvOpenAPIKey),
		OpenAPIModel:    os.Getenv(EnvOpenAPIModel),
		AllowedUserIDs:  parseUserIDs(os.Getenv(EnvAllowedUserIDs)),
		SQLDSN:          sqlDSN,
	}

	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf(
			"missing Telegram token (env %q): %w",
			EnvTelegramToken,
			ErrMissingTelegramToken,
		)
	}

	if cfg.OpenAPIKey == "" {
		return nil, fmt.Errorf(
			"missing OpenAPI key (env %q): %w",
			EnvOpenAPIKey,
			ErrMissingOpenAPIKey,
		)
	}

	if cfg.OpenAPIModel == "" {
		return nil, fmt.Errorf(
			"missing OpenAPI model (env %q): %w",
			EnvOpenAPIModel,
			ErrMissingOpenAPIModel,
		)
	}

	return cfg, nil
}

func parseUserIDs(s string) []int64 {
	if s == "" {
		return nil
	}

	var ids []int64
	for part := range strings.SplitSeq(s, ",") {
		if id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}
