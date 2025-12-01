package test

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/schraf/assistant/pkg/models"
)

// MockJobScheduler is a mock implementation of interfaces.JobScheduler.
type MockJobScheduler struct {
	ScheduleJobFunc func(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error
}

// ScheduleJob calls ScheduleJobFunc if set, otherwise returns nil.
func (m *MockJobScheduler) ScheduleJob(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error {
	if m.ScheduleJobFunc != nil {
		return m.ScheduleJobFunc(ctx, contentType, config, request)
	}
	return nil
}

// MockPublisher is a mock implementation of interfaces.Publisher.
type MockPublisher struct {
	PublishDocumentFunc func(ctx context.Context, doc *models.Document) (*url.URL, error)
}

// PublishDocument calls PublishDocumentFunc if set, otherwise returns a mock URL.
func (m *MockPublisher) PublishDocument(ctx context.Context, doc *models.Document) (*url.URL, error) {
	if m.PublishDocumentFunc != nil {
		return m.PublishDocumentFunc(ctx, doc)
	}
	return url.Parse("https://telegra.ph/mock-page")
}

// MockNotifier is a mock implementation of interfaces.Notifier.
type MockNotifier struct {
	SendPublishedURLNotificationFunc func(publishedURL *url.URL, title string) error
}

// SendPublishedURLNotification calls SendPublishedURLNotificationFunc if set, otherwise returns nil.
func (m *MockNotifier) SendPublishedURLNotification(publishedURL *url.URL, title string) error {
	if m.SendPublishedURLNotificationFunc != nil {
		return m.SendPublishedURLNotificationFunc(publishedURL, title)
	}
	return nil
}

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
