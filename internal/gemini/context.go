package gemini

import "context"

type contextKey struct{}

var modelKey contextKey

const defaultModel = "gemini-flash-latest"

func modelFromContext(ctx context.Context) string {
	if model, ok := ctx.Value(modelKey).(string); ok {
		return model
	}

	return defaultModel
}
