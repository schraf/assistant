package test

import (
	"context"

	"github.com/schraf/assistant/pkg/generators"
	"github.com/schraf/assistant/pkg/models"
)

func init() {
	generators.MustRegister("test-generator", factory)
}

func factory(generators.Config) (models.ContentGenerator, error) {
	return &generator{}, nil
}

type generator struct{}

func (g *generator) Generate(ctx context.Context, request models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
	return &models.Document{
		Title:  "Test Document",
		Author: "Some person",
		Sections: []models.DocumentSection{
			{
				Title: "Test Section",
				Paragraphs: []string{
					"Test Paragraph",
				},
			},
		},
	}, nil
}
