package eval

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/schraf/assistant/internal/gemini"
	"github.com/schraf/assistant/internal/mocks"
	"github.com/schraf/assistant/internal/ollama"
	"github.com/schraf/assistant/pkg/models"
)

func newAssistant(ctx context.Context) (models.Assistant, error) {
	provider := strings.ToLower(os.Getenv("ASSISTANT_PROVIDER"))
	switch provider {
	case "mock":
		return &mocks.MockAssistant{}, nil
	case "ollama":
		return ollama.NewClient(ctx)
	case "gemini", "":
		return gemini.NewClient(ctx)
	default:
		return nil, fmt.Errorf("unknown assistant provider: %s", provider)
	}
}

func Evaluate(ctx context.Context, generator models.ContentGenerator, request models.ContentRequest, model string) error {
	assistant, err := newAssistant(ctx)
	if err != nil {
		return fmt.Errorf("failed creating assistant client: %w", err)
	}

	ctx = assistant.WithModel(ctx, model)
	doc, err := generator.Generate(ctx, request, assistant)
	if err != nil {
		return fmt.Errorf("failed generating content: %w", err)
	}

	now := time.Now()
	layout := "January 2, 2006"
	dateString := now.Format(layout)

	fmt.Println("# " + doc.Title)
	fmt.Println(dateString + " · " + doc.Author)

	for _, section := range doc.Sections {
		fmt.Println("## " + section.Title)

		for _, paragraph := range section.Paragraphs {
			fmt.Println(paragraph + "\n")
		}
	}

	return nil
}
