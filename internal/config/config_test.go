package config_test

import (
	"testing"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNew_MissingTelegramToken(t *testing.T) {
	withEnv(t, map[string]string{
		config.EnvOpenAPIKey:   "openai-key",
		config.EnvOpenAPIModel: "gpt-4o-mini",
	})

	_, err := config.New()

	require.ErrorIs(t, err, config.ErrMissingTelegramToken)
}

func TestNew_MissingOpenAPIKey(t *testing.T) {
	withEnv(t, map[string]string{
		config.EnvTelegramToken: "tg-token",
		config.EnvOpenAPIModel:  "gpt-4o-mini",
	})

	_, err := config.New()

	require.ErrorIs(t, err, config.ErrMissingOpenAPIKey)
}

func TestNew_MissingOpenAPIModel(t *testing.T) {
	withEnv(t, map[string]string{
		config.EnvTelegramToken: "tg-token",
		config.EnvOpenAPIKey:    "openai-key",
	})

	_, err := config.New()

	require.ErrorIs(t, err, config.ErrMissingOpenAPIModel)
}

func TestNew_DefaultSQLDSN(t *testing.T) {
	withMinimalEnv(t)

	cfg, err := config.New()

	require.NoError(t, err)
	require.Equal(t, "file::memory:?cache=shared", cfg.SQLDSN)
}

func TestNew_SingleAllowedUserID(t *testing.T) {
	withMinimalEnv(t)
	withEnv(t, map[string]string{
		config.EnvAllowedUserIDs: "123456789",
	})

	cfg, err := config.New()

	require.NoError(t, err)
	require.Equal(t, []int64{123456789}, cfg.AllowedUserIDs)
}

func TestNew_MultipleAllowedUserID(t *testing.T) {
	withMinimalEnv(t)
	withEnv(t, map[string]string{
		config.EnvAllowedUserIDs: "123456789,987654321",
	})

	cfg, err := config.New()

	require.NoError(t, err)
	require.Equal(t, []int64{123456789, 987654321}, cfg.AllowedUserIDs)
}

func TestNew_AllFieldsFilled(t *testing.T) {
	withEnv(t, map[string]string{
		config.EnvTelegramToken:   "tg-token",
		config.EnvOpenAPIKey:      "openai-key",
		config.EnvOpenAPIModel:    "gpt-4o-mini",
		config.EnvOpenAPIBasePath: "https://api.openai.com/v1",
		config.EnvAllowedUserIDs:  "42,  1001,3002",
		config.EnvSQLDSN:          "user:pass@/dbname",
	})

	cfg, err := config.New()
	require.NoError(t, err)

	require.Equal(t, "tg-token", cfg.TelegramToken)
	require.Equal(t, "openai-key", cfg.OpenAPIKey)
	require.Equal(t, "gpt-4o-mini", cfg.OpenAPIModel)
	require.Equal(t, "https://api.openai.com/v1", cfg.OpenAPIBasePath)
	require.Equal(t, "user:pass@/dbname", cfg.SQLDSN)
	require.Equal(t, []int64{42, 1001, 3002}, cfg.AllowedUserIDs)
}

func withMinimalEnv(t *testing.T) {
	t.Helper()
	withEnv(t, map[string]string{
		config.EnvTelegramToken: "tg-token",
		config.EnvOpenAPIKey:    "openai-key",
		config.EnvOpenAPIModel:  "gpt-4o-mini",
	})
}

func withEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v)
	}
}
