package gemini

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
)

type contextKey struct{}

var modelKey contextKey

const defaultLiteModel = "gemini-flash-latest"
const defaultDeepModel = "gemini-pro-latest"

func modelFromContext(ctx context.Context) string {
	if modelType, ok := ctx.Value(modelKey).(models.ModelType); ok {
		if modelType == models.ModelDeep {
			return defaultDeepModel
		}
	}

	return defaultLiteModel
}
