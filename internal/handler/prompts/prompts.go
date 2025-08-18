package prompts

import (
	_ "embed"
	"strings"
)

//go:embed system_prompt.txt
var systemPrompt string

func SystemPrompt() string {
	return strings.TrimSpace(systemPrompt)
}
