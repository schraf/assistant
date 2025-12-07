package telegraph

import (
	"context"
	"fmt"
	"net/url"
	"os"

	internal_models "github.com/schraf/assistant/internal/models"
	"github.com/schraf/assistant/pkg/models"
)

// Publisher implements internal_models.Publisher using the Telegraph API.
type Publisher struct {
	client Client
}

// NewPublisher creates a new Publisher.
func NewPublisher() internal_models.Publisher {
	return &Publisher{
		client: NewDefaultClient(),
	}
}

// PublishDocument publishes a document to Telegraph and returns its URL.
func (p *Publisher) PublishDocument(ctx context.Context, doc *models.Document) (*url.URL, error) {
	apiToken := os.Getenv("TELEGRAPH_API_KEY")
	if apiToken == "" {
		return nil, fmt.Errorf("missing TELEGRAPH_API_KEY environment variable")
	}

	content := Nodes{}

	for _, section := range doc.Sections {
		content = append(content, NodeElement{
			Tag: "h3",
			Children: Nodes{
				section.Title,
			},
		})

		for _, paragraph := range section.Paragraphs {
			content = append(content, NodeElement{
				Tag: "p",
				Children: Nodes{
					paragraph,
				},
			})
		}
	}

	returnContent := false

	pageRequest := CreatePageRequest{
		AccessToken:   apiToken,
		Title:         doc.Title,
		AuthorName:    &doc.Author,
		Content:       content,
		ReturnContent: &returnContent,
	}

	page, err := p.client.CreatePage(ctx, pageRequest)
	if err != nil {
		return nil, err
	}

	pageURL, err := url.Parse(page.URL)
	if err != nil {
		return nil, err
	}

	return pageURL, nil
}
