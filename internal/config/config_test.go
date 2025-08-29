package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_MissingTelegramToken(t *testing.T) {
	withEnv(t, map[string]string{
		EnvOpenAPIKey:   "openai-key",
		EnvOpenAPIModel: "gpt-4o-mini",
	})

	_, err := New()

	require.ErrorIs(t, err, ErrMissingTelegramToken)
}

func TestNew_MissingOpenAPIKey(t *testing.T) {
	withEnv(t, map[string]string{
		EnvTelegramToken: "tg-token",
		EnvOpenAPIModel:  "gpt-4o-mini",
	})

	_, err := New()

	require.ErrorIs(t, err, ErrMissingOpenAPIKey)
}

func TestNew_MissingOpenAPIModel(t *testing.T) {
	withEnv(t, map[string]string{
		EnvTelegramToken: "tg-token",
		EnvOpenAPIKey:    "openai-key",
	})

	_, err := New()

	require.ErrorIs(t, err, ErrMissingOpenAPIModel)
}

func withEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v)
	}
}
