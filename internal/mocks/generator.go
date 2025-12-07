package mocks

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
)

// MockContentGenerator is a mock implementation of models.ContentGenerator.
type MockContentGenerator struct {
	GenerateFunc func(ctx context.Context, request models.ContentRequest, assistant models.Assistant) (*models.Document, error)
}

// Generate calls GenerateFunc if set, otherwise returns a mock document.
func (m *MockContentGenerator) Generate(ctx context.Context, request models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, request, assistant)
	}

	return &models.Document{
		Title:  "Mock Document",
		Author: "Mock Author",
		Sections: []models.DocumentSection{
			{
				Title:      "Mock Section",
				Paragraphs: []string{"Mock paragraph 1", "Mock paragraph 2"},
			},
		},
	}, nil
}
