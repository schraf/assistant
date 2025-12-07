package mocks

import (
	"context"
	"encoding/json"
)

// MockAssistant is a mock implementation of models.Assistant.
type MockAssistant struct {
	AskFunc           func(ctx context.Context, persona string, request string) (*string, error)
	StructuredAskFunc func(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error)
	WithModelFunc     func(ctx context.Context, model string) context.Context
}

// Ask calls AskFunc if set, otherwise returns a mock response.
func (m *MockAssistant) Ask(ctx context.Context, persona string, request string) (*string, error) {
	if m.AskFunc != nil {
		return m.AskFunc(ctx, persona, request)
	}

	response := "Mock response"
	return &response, nil
}

// StructuredAsk calls StructuredAskFunc if set, otherwise returns mock JSON.
func (m *MockAssistant) StructuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error) {
	if m.StructuredAskFunc != nil {
		return m.StructuredAskFunc(ctx, persona, request, schema)
	}

	return json.RawMessage(`{"mock": "data"}`), nil
}

// WithModel calls WithModelFunc if set, otherwise returns the context unchanged.
func (m *MockAssistant) WithModel(ctx context.Context, model string) context.Context {
	if m.WithModelFunc != nil {
		return m.WithModelFunc(ctx, model)
	}

	return ctx
}
