package gemini

import "context"

type contextKey struct{}

var modelKey contextKey

const defaultModel = "gemini-flash-latest"

func WithModel(ctx context.Context, model string) context.Context {
	return context.WithValue(ctx, modelKey, model)
}

func WithProModel(ctx context.Context) context.Context {
	return WithModel(ctx, "gemini-pro-latest")
}

func WithFlashModel(ctx context.Context) context.Context {
	return WithModel(ctx, "gemini-flash-latest")
}

func WithFlashLiteModel(ctx context.Context) context.Context {
	return WithModel(ctx, "gemini-flash-lite-latest")
}

func modelFromContext(ctx context.Context) string {
	if model, ok := ctx.Value(modelKey).(string); ok {
		return model
	}

	return defaultModel
}
