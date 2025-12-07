package mocks

import (
	"context"
	"encoding/json"
	"fmt"
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

// StructuredAsk calls StructuredAskFunc if set, otherwise returns mock JSON that matches the schema.
func (m *MockAssistant) StructuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error) {
	if m.StructuredAskFunc != nil {
		return m.StructuredAskFunc(ctx, persona, request, schema)
	}

	// Generate mock data based on the schema
	mockData := generateMockFromSchema(schema)
	mockJSON, err := json.Marshal(mockData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mock data: %w", err)
	}

	return json.RawMessage(mockJSON), nil
}

// generateMockFromSchema generates mock data that matches the provided JSON schema.
func generateMockFromSchema(schema map[string]any) any {
	if schema == nil {
		return map[string]any{"mock": "data"}
	}

	// Check if this is a JSON schema with a "type" field
	schemaType, hasType := schema["type"].(string)
	if !hasType {
		// If no type, check for "properties" (object schema)
		if properties, ok := schema["properties"].(map[string]any); ok {
			return generateObjectFromProperties(properties, schema)
		}
		// Fallback to generic object
		return map[string]any{"mock": "data"}
	}

	switch schemaType {
	case "object":
		return generateObjectFromSchema(schema)
	case "array":
		return generateArrayFromSchema(schema)
	case "string":
		return "mock_string"
	case "number", "integer":
		return 42
	case "boolean":
		return true
	default:
		return map[string]any{"mock": "data"}
	}
}

// generateObjectFromSchema generates a mock object from a JSON schema.
func generateObjectFromSchema(schema map[string]any) map[string]any {
	result := make(map[string]any)

	// Get properties
	if properties, ok := schema["properties"].(map[string]any); ok {
		for key, propSchema := range properties {
			if propMap, ok := propSchema.(map[string]any); ok {
				result[key] = generateMockFromSchema(propMap)
			} else {
				result[key] = "mock_value"
			}
		}
	}

	// If no properties found, return a simple mock object
	if len(result) == 0 {
		result["mock"] = "data"
	}

	return result
}

// generateObjectFromProperties generates a mock object from properties map.
func generateObjectFromProperties(properties map[string]any, schema map[string]any) map[string]any {
	result := make(map[string]any)

	for key, propSchema := range properties {
		if propMap, ok := propSchema.(map[string]any); ok {
			result[key] = generateMockFromSchema(propMap)
		} else {
			result[key] = "mock_value"
		}
	}

	// If no properties found, return a simple mock object
	if len(result) == 0 {
		result["mock"] = "data"
	}

	return result
}

// generateArrayFromSchema generates a mock array from a JSON schema.
func generateArrayFromSchema(schema map[string]any) []any {
	result := []any{}

	// Get items schema
	if items, ok := schema["items"].(map[string]any); ok {
		// Generate a single item based on the items schema
		result = append(result, generateMockFromSchema(items))
	} else {
		// Default to a single string item
		result = append(result, "mock_item")
	}

	return result
}

// WithModel calls WithModelFunc if set, otherwise returns the context unchanged.
func (m *MockAssistant) WithModel(ctx context.Context, model string) context.Context {
	if m.WithModelFunc != nil {
		return m.WithModelFunc(ctx, model)
	}

	return ctx
}
