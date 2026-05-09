package service

import (
	"fmt"
	"strings"
)

const openAIImplicitStickyPrefix = "implicit:"

func buildOpenAIImplicitStickyKey(apiKeyID int64, requestedModel string) string {
	if apiKeyID <= 0 {
		return ""
	}
	model := strings.ToLower(strings.TrimSpace(requestedModel))
	if model == "" {
		model = "unknown"
	}
	return fmt.Sprintf("%s%d:%s", openAIImplicitStickyPrefix, apiKeyID, model)
}
