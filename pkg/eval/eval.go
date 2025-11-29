package eval

import (
	"context"
	"fmt"
	"time"

	"github.com/schraf/assistant/internal/gemini"
	"github.com/schraf/assistant/pkg/models"
)

func Evaluate(ctx context.Context, generator models.ContentGenerator, request models.ContentRequest) error {
	assistant, err := gemini.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed creating assistant client: %w", err)
	}

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
