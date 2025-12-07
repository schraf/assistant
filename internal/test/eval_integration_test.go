package test

import (
	"context"
	"os"
	"strings"
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
	os.Setenv("ASSISTANT_PROVIDER", "gemini")
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

	// Note: Since newAssistant is private and creates real clients,
	// we'll test the flow but may need to handle cases where real clients
	// aren't available. For now, we'll test with gemini provider.
	// If API keys aren't available, the test will fail at assistant creation,
	// which is acceptable for an integration test.

	ctx := context.Background()
	model := "gemini-pro-latest"

	err := eval.Evaluate(ctx, generator, request, model)

	// If assistant creation fails (e.g., no API key), skip the test
	if err != nil && strings.Contains(err.Error(), "failed creating assistant client") {
		t.Skip("Assistant client creation failed (likely missing API keys), skipping integration test")
	}

	require.NoError(t, err, "Evaluate should succeed")

	// Verify generator was called correctly
	assert.Equal(t, requestID, capturedRequest.Id, "generator should be called with correct request ID")
	assert.Equal(t, "AI", capturedRequest.Body["topic"], "generator should be called with correct request body")
	assert.NotNil(t, capturedAssistant, "generator should be called with assistant")
}

func TestEvaluate_Integration_ModelParameter(t *testing.T) {
	// Test that different models are passed correctly
	originalProvider := os.Getenv("ASSISTANT_PROVIDER")
	os.Setenv("ASSISTANT_PROVIDER", "gemini")
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
		Body: map[string]any{"topic": "Testing"},
	}

	generator := &mocks.MockContentGenerator{
		GenerateFunc: func(ctx context.Context, req models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
			return &models.Document{
				Title:  "Model Test",
				Author: "Test",
				Sections: []models.DocumentSection{
					{
						Title:      "Test",
						Paragraphs: []string{"Test paragraph"},
					},
				},
			}, nil
		},
	}

	ctx := context.Background()
	testModel := "custom-model-name"

	err := eval.Evaluate(ctx, generator, request, testModel)

	if err != nil && strings.Contains(err.Error(), "failed creating assistant client") {
		t.Skip("Assistant client creation failed (likely missing API keys), skipping integration test")
	}

	require.NoError(t, err, "Evaluate should succeed")
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
	model := "test-model"

	err := eval.Evaluate(ctx, generator, request, model)

	require.Error(t, err, "Evaluate should return error for unknown provider")
	assert.Contains(t, err.Error(), "failed creating assistant client", "error should mention assistant client creation")
	assert.Contains(t, err.Error(), "unknown assistant provider", "error should mention unknown provider")
}

func TestEvaluate_Integration_GeneratorError(t *testing.T) {
	// Test error handling when generator fails
	originalProvider := os.Getenv("ASSISTANT_PROVIDER")
	os.Setenv("ASSISTANT_PROVIDER", "gemini")
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
	model := "test-model"

	err := eval.Evaluate(ctx, generator, request, model)

	// If assistant creation fails, skip
	if err != nil && strings.Contains(err.Error(), "failed creating assistant client") {
		t.Skip("Assistant client creation failed (likely missing API keys), skipping integration test")
	}

	require.Error(t, err, "Evaluate should return error when generator fails")
	assert.Contains(t, err.Error(), "failed generating content", "error should mention content generation failure")
}
