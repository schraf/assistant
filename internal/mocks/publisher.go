package mocks

import (
	"context"
	"net/url"

	"github.com/schraf/assistant/pkg/models"
)

// MockPublisher is a mock implementation of models.Publisher.
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
