package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/schraf/assistant/internal/service"
	"github.com/schraf/assistant/pkg/models"
)

func TestHandler_Integration(t *testing.T) {
	// Set up test environment
	apiToken := "test-api-token"
	os.Setenv("API_TOKEN", apiToken)
	defer os.Unsetenv("API_TOKEN")

	// Create mock scheduler
	mockScheduler := &MockJobScheduler{
		ScheduleJobFunc: func(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error {
			// Verify the job was called with correct parameters
			assert.Equal(t, "article", contentType, "contentType should be 'article'")
			// Header keys are case-sensitive, check both possible cases
			modelValue := config["model"]
			if modelValue == nil {
				modelValue = config["Model"] // HTTP headers can be canonicalized
			}
			assert.Equal(t, "pro", modelValue, "config model should be 'pro'")
			assert.Equal(t, "AI", request.Body["topic"], "request body topic should be 'AI'")
			return nil
		},
	}

	// Create handler
	handler := service.NewHandler(mockScheduler)

	// Create test request
	requestBody := map[string]any{"topic": "AI"}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/content", bytes.NewReader(bodyBytes))
	req.Header.Set("X-API-Token", apiToken)
	req.Header.Set("X-Content-Type", "article")
	req.Header.Set("X-Config-model", "pro")

	w := httptest.NewRecorder()

	// Call handler
	handler.HandleRequest(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "status code should be 200")

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err, "should decode response JSON")

	assert.True(t, response["success"].(bool), "success should be true")

	_, err = uuid.Parse(response["request_id"].(string))
	assert.NoError(t, err, "request_id should be a valid UUID")
}

func TestHandler_Integration_MissingToken(t *testing.T) {
	os.Setenv("API_TOKEN", "test-token")
	defer os.Unsetenv("API_TOKEN")

	mockScheduler := &MockJobScheduler{}
	handler := service.NewHandler(mockScheduler)

	req := httptest.NewRequest("POST", "/content", bytes.NewReader([]byte("{}")))
	// Missing X-API-Token header

	w := httptest.NewRecorder()
	handler.HandleRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "status code should be 400")
}

func TestHandler_Integration_InvalidToken(t *testing.T) {
	os.Setenv("API_TOKEN", "correct-token")
	defer os.Unsetenv("API_TOKEN")

	mockScheduler := &MockJobScheduler{}
	handler := service.NewHandler(mockScheduler)

	req := httptest.NewRequest("POST", "/content", bytes.NewReader([]byte("{}")))
	req.Header.Set("X-API-Token", "wrong-token")
	req.Header.Set("X-Content-Type", "article")

	w := httptest.NewRecorder()
	handler.HandleRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "status code should be 401")
}

func TestHandler_Integration_MissingContentType(t *testing.T) {
	os.Setenv("API_TOKEN", "test-token")
	defer os.Unsetenv("API_TOKEN")

	mockScheduler := &MockJobScheduler{}
	handler := service.NewHandler(mockScheduler)

	req := httptest.NewRequest("POST", "/content", bytes.NewReader([]byte("{}")))
	req.Header.Set("X-API-Token", "test-token")
	// Missing X-Content-Type header

	w := httptest.NewRecorder()
	handler.HandleRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "status code should be 400")
}

func TestHandler_Integration_InvalidJSON(t *testing.T) {
	os.Setenv("API_TOKEN", "test-token")
	defer os.Unsetenv("API_TOKEN")

	mockScheduler := &MockJobScheduler{}
	handler := service.NewHandler(mockScheduler)

	req := httptest.NewRequest("POST", "/content", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("X-API-Token", "test-token")
	req.Header.Set("X-Content-Type", "article")

	w := httptest.NewRecorder()
	handler.HandleRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "status code should be 400")
}

func TestHandler_Integration_SchedulerError(t *testing.T) {
	os.Setenv("API_TOKEN", "test-token")
	defer os.Unsetenv("API_TOKEN")

	mockScheduler := &MockJobScheduler{
		ScheduleJobFunc: func(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error {
			return http.ErrBodyReadAfterClose
		},
	}

	handler := service.NewHandler(mockScheduler)

	requestBody := map[string]any{"topic": "AI"}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/content", bytes.NewReader(bodyBytes))
	req.Header.Set("X-API-Token", "test-token")
	req.Header.Set("X-Content-Type", "article")

	w := httptest.NewRecorder()
	handler.HandleRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "status code should be 500")
}
