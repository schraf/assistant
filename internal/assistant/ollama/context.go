package ollama

import "context"

type contextKey struct{}

var modelKey contextKey

const defaultModel = "llama3.2"

func modelFromContext(ctx context.Context) string {
	if model, ok := ctx.Value(modelKey).(string); ok {
		return model
	}

	return defaultModel
}
