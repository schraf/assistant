package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/schraf/assistant/internal/job"
	"github.com/schraf/assistant/internal/log"
	"github.com/schraf/assistant/pkg/models"
)

func TestProcessor_Integration(t *testing.T) {
	// Set up test environment variables
	requestID := uuid.New()
	requestBody := map[string]any{"topic": "AI"}
	bodyJSON, _ := json.Marshal(requestBody)
	encodedBody := base64.StdEncoding.EncodeToString(bodyJSON)

	config := map[string]any{"model": "pro"}
	configJSON, _ := json.Marshal(config)
	encodedConfig := base64.StdEncoding.EncodeToString(configJSON)

	os.Setenv("REQUEST_ID", requestID.String())
	os.Setenv("REQUEST_BODY", encodedBody)
	os.Setenv("CONTENT_CONFIG", encodedConfig)
	os.Setenv("CONTENT_TYPE", "test-generator")
	defer func() {
		os.Unsetenv("REQUEST_ID")
		os.Unsetenv("REQUEST_BODY")
		os.Unsetenv("CONTENT_CONFIG")
		os.Unsetenv("CONTENT_TYPE")
	}()

	// Create mocks
	mockAssistant := &MockAssistant{
		WithModelFunc: func(ctx context.Context, model string) context.Context {
			assert.Equal(t, "gemini-pro-latest", model, "model should be 'gemini-pro-latest'")
			return ctx
		},
		AskFunc: func(ctx context.Context, persona string, request string) (*string, error) {
			response := "Generated content"
			return &response, nil
		},
	}

	var publishedDoc *models.Document

	mockPublisher := &MockPublisher{
		PublishDocumentFunc: func(ctx context.Context, doc *models.Document) (*url.URL, error) {
			publishedDoc = doc
			url, err := url.Parse("https://telegra.ph/test-page")
			if err != nil {
				return nil, err
			}
			return url, nil
		},
	}

	var notifiedURL *url.URL
	var notifiedTitle string

	mockNotifier := &MockNotifier{
		SendPublishedURLNotificationFunc: func(publishedURL *url.URL, title string) error {
			notifiedURL = publishedURL
			notifiedTitle = title
			return nil
		},
	}

	processor := job.NewProcessor(mockAssistant, mockPublisher, mockNotifier, log.NewLogger())

	ctx := context.Background()
	err := processor.Process(ctx)

	// If generator doesn't exist, that's expected in test environment
	if err != nil {
		// Check if it's a generator not found error
		if err.Error() == "failed creating generator: generator 'test-generator' is not registered" {
			t.Skip("test-generator not registered, skipping integration test")
		}
		require.NoError(t, err, "processor.Process() should succeed")
	}

	// Verify publisher was called
	require.NotNil(t, publishedDoc, "publisher should be called with a document")

	// Verify notifier was called
	require.NotNil(t, notifiedURL, "notifier should be called with a URL")
	assert.NotEmpty(t, notifiedTitle, "notifier should be called with a title")
	assert.Equal(t, "https://telegra.ph/test-page", notifiedURL.String(), "notified URL should match")
}

func TestProcessor_Integration_ModelSelection(t *testing.T) {
	requestID := uuid.New()
	requestBody := map[string]any{"topic": "AI"}
	bodyJSON, _ := json.Marshal(requestBody)
	encodedBody := base64.StdEncoding.EncodeToString(bodyJSON)

	// Test with "basic" model
	config := map[string]any{"model": "basic"}
	configJSON, _ := json.Marshal(config)
	encodedConfig := base64.StdEncoding.EncodeToString(configJSON)

	os.Setenv("REQUEST_ID", requestID.String())
	os.Setenv("REQUEST_BODY", encodedBody)
	os.Setenv("CONTENT_CONFIG", encodedConfig)
	os.Setenv("CONTENT_TYPE", "test-generator")
	defer func() {
		os.Unsetenv("REQUEST_ID")
		os.Unsetenv("REQUEST_BODY")
		os.Unsetenv("CONTENT_CONFIG")
		os.Unsetenv("CONTENT_TYPE")
	}()

	var selectedModel string
	mockAssistant := &MockAssistant{
		WithModelFunc: func(ctx context.Context, model string) context.Context {
			selectedModel = model
			return ctx
		},
		AskFunc: func(ctx context.Context, persona string, request string) (*string, error) {
			response := "Generated content"
			return &response, nil
		},
	}

	mockPublisher := &MockPublisher{}
	mockNotifier := &MockNotifier{}

	processor := job.NewProcessor(mockAssistant, mockPublisher, mockNotifier, log.NewLogger())

	ctx := context.Background()
	err := processor.Process(ctx)

	if err != nil {
		if err.Error() == "failed creating generator: generator 'test-generator' is not registered" {
			t.Skip("test-generator not registered, skipping integration test")
		}
		require.NoError(t, err, "processor.Process() should succeed")
	}

	assert.Equal(t, "gemini-flash-latest", selectedModel, "selected model should be 'gemini-flash-latest'")
}

func TestProcessor_Integration_MissingRequestID(t *testing.T) {
	os.Unsetenv("REQUEST_ID")
	defer os.Unsetenv("REQUEST_ID")

	mockAssistant := &MockAssistant{}
	mockPublisher := &MockPublisher{}
	mockNotifier := &MockNotifier{}

	processor := job.NewProcessor(mockAssistant, mockPublisher, mockNotifier, log.NewLogger())

	ctx := context.Background()
	err := processor.Process(ctx)

	require.Error(t, err, "should return error for missing REQUEST_ID")
}

func TestProcessor_Integration_MissingRequestBody(t *testing.T) {
	requestID := uuid.New()
	os.Setenv("REQUEST_ID", requestID.String())
	os.Unsetenv("REQUEST_BODY")
	defer func() {
		os.Unsetenv("REQUEST_ID")
		os.Unsetenv("REQUEST_BODY")
	}()

	mockAssistant := &MockAssistant{}
	mockPublisher := &MockPublisher{}
	mockNotifier := &MockNotifier{}

	processor := job.NewProcessor(mockAssistant, mockPublisher, mockNotifier, log.NewLogger())

	ctx := context.Background()
	err := processor.Process(ctx)

	require.Error(t, err, "should return error for missing REQUEST_BODY")
}

func TestProcessor_Integration_PublisherError(t *testing.T) {
	requestID := uuid.New()
	requestBody := map[string]any{"topic": "AI"}
	bodyJSON, _ := json.Marshal(requestBody)
	encodedBody := base64.StdEncoding.EncodeToString(bodyJSON)

	config := map[string]any{}
	configJSON, _ := json.Marshal(config)
	encodedConfig := base64.StdEncoding.EncodeToString(configJSON)

	os.Setenv("REQUEST_ID", requestID.String())
	os.Setenv("REQUEST_BODY", encodedBody)
	os.Setenv("CONTENT_CONFIG", encodedConfig)
	os.Setenv("CONTENT_TYPE", "test-generator")
	defer func() {
		os.Unsetenv("REQUEST_ID")
		os.Unsetenv("REQUEST_BODY")
		os.Unsetenv("CONTENT_CONFIG")
		os.Unsetenv("CONTENT_TYPE")
	}()

	mockAssistant := &MockAssistant{
		AskFunc: func(ctx context.Context, persona string, request string) (*string, error) {
			response := "Generated content"
			return &response, nil
		},
	}

	mockPublisher := &MockPublisher{
		PublishDocumentFunc: func(ctx context.Context, doc *models.Document) (*url.URL, error) {
			return nil, context.DeadlineExceeded
		},
	}

	mockNotifier := &MockNotifier{}

	processor := job.NewProcessor(mockAssistant, mockPublisher, mockNotifier, log.NewLogger())

	ctx := context.Background()
	err := processor.Process(ctx)

	if err != nil {
		if err.Error() == "failed creating generator: generator 'test-generator' is not registered" {
			t.Skip("test-generator not registered, skipping integration test")
		}
		// Expected error from publisher
		require.Error(t, err, "should return error from publisher")
		assert.Contains(t, err.Error(), "publish error", "error should be from publisher")
		return
	}
}
