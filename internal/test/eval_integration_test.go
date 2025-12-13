package test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/schraf/assistant/internal/mocks"
	"github.com/schraf/assistant/pkg/eval"
	"github.com/schraf/assistant/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluate_Integration(t *testing.T) {
	// Set up test environment - use gemini as default provider
	originalProvider := os.Getenv("ASSISTANT_PROVIDER")
	os.Setenv("ASSISTANT_PROVIDER", "mock")
	defer func() {
		if originalProvider == "" {
			os.Unsetenv("ASSISTANT_PROVIDER")
		} else {
			os.Setenv("ASSISTANT_PROVIDER", originalProvider)
		}
	}()

	// Create test request
	requestID := uuid.New()
	request := models.ContentRequest{
		Id:   requestID,
		Body: map[string]any{"topic": "AI"},
	}

	// Track generator calls
	var capturedRequest models.ContentRequest
	var capturedAssistant models.Assistant

	// Create mock generator
	generator := &mocks.MockContentGenerator{
		GenerateFunc: func(ctx context.Context, req models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
			capturedRequest = req
			capturedAssistant = assistant
			return &models.Document{
				Title:  "Test Document",
				Author: "Test Author",
				Sections: []models.DocumentSection{
					{
						Title:      "Introduction",
						Paragraphs: []string{"This is paragraph 1.", "This is paragraph 2."},
					},
					{
						Title:      "Conclusion",
						Paragraphs: []string{"This is the conclusion."},
					},
				},
			}, nil
		},
	}

	ctx := context.Background()
	model := "mock-model"

	err := eval.Evaluate(ctx, generator, request, &model)
	require.NoError(t, err, "Evaluate should succeed")

	// Verify generator was called correctly
	assert.Equal(t, requestID, capturedRequest.Id, "generator should be called with correct request ID")
	assert.Equal(t, "AI", capturedRequest.Body["topic"], "generator should be called with correct request body")
	assert.NotNil(t, capturedAssistant, "generator should be called with assistant")
}

func TestEvaluate_Integration_UnknownProvider(t *testing.T) {
	// Test error handling for unknown provider
	originalProvider := os.Getenv("ASSISTANT_PROVIDER")
	os.Setenv("ASSISTANT_PROVIDER", "unknown-provider")
	defer func() {
		if originalProvider == "" {
			os.Unsetenv("ASSISTANT_PROVIDER")
		} else {
			os.Setenv("ASSISTANT_PROVIDER", originalProvider)
		}
	}()

	requestID := uuid.New()
	request := models.ContentRequest{
		Id:   requestID,
		Body: map[string]any{"topic": "Test"},
	}

	generator := &mocks.MockContentGenerator{}

	ctx := context.Background()
	model := "unknown-model"

	err := eval.Evaluate(ctx, generator, request, &model)

	require.Error(t, err, "Evaluate should return error for unknown provider")
	assert.Contains(t, err.Error(), "failed creating assistant client", "error should mention assistant client creation")
	assert.Contains(t, err.Error(), "unknown assistant provider", "error should mention unknown provider")
}

func TestEvaluate_Integration_GeneratorError(t *testing.T) {
	// Test error handling when generator fails
	originalProvider := os.Getenv("ASSISTANT_PROVIDER")
	os.Setenv("ASSISTANT_PROVIDER", "mock")
	defer func() {
		if originalProvider == "" {
			os.Unsetenv("ASSISTANT_PROVIDER")
		} else {
			os.Setenv("ASSISTANT_PROVIDER", originalProvider)
		}
	}()

	requestID := uuid.New()
	request := models.ContentRequest{
		Id:   requestID,
		Body: map[string]any{"topic": "Test"},
	}

	generator := &mocks.MockContentGenerator{
		GenerateFunc: func(ctx context.Context, req models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
			return nil, assert.AnError
		},
	}

	ctx := context.Background()
	model := "mock-model"

	err := eval.Evaluate(ctx, generator, request, &model)
	require.Error(t, err, "Evaluate should return error when generator fails")
	assert.Contains(t, err.Error(), "failed generating content", "error should mention content generation failure")
}
