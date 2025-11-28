package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/schraf/assistant/internal/gemini"
	"github.com/schraf/assistant/internal/telegraph"
	"github.com/schraf/assistant/internal/utils"
	"github.com/schraf/assistant/pkg/models"
)

func main() {
	ctx := context.Background()
	logger := utils.NewLogger()

	if err := utils.LoadEnv(".env"); err != nil {
		logger.ErrorContext(ctx, "load_env_failed",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	//--========================================================================--
	//--== GET THE REQUEST
	//--========================================================================--

	request, err := getRequest()
	if err != nil {
		logger.ErrorContext(ctx, "invalid_request",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	logger = logger.With(
		slog.String("request_id", request.Id.String()),
	)

	//--========================================================================--
	//--== CREATE THE ASSISTANT
	//--========================================================================--

	assistant, err := gemini.NewClient(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed_creating_assistant",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	//--========================================================================--
	//--== GENERATE CONTENT
	//--========================================================================--

	logger.InfoContext(ctx, "generating_content")

	/** TODO
	doc, err := generator.Generate(ctx, *request, assistant)
	if err != nil {
		logger.ErrorContext(ctx, "content_generation_error",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
	**/

	logger.InfoContext(ctx, "generated_document",
		slog.String("title", doc.Title),
		slog.String("author", doc.Author),
		slog.Int("section_count", len(doc.Sections)),
	)

	//--========================================================================--
	//--== PUBLISH CONTENT
	//--========================================================================--

	logger.InfoContext(ctx, "publishing_document")

	url, err := publishDocument(ctx, doc)
	if err != nil {
		logger.ErrorContext(ctx, "publish_error",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	logger.InfoContext(ctx, "published_document",
		slog.String("url", url.String()),
	)

	logger.InfoContext(ctx, "job_completed_successfully")
	os.Exit(0)
}

func getRequest() (*models.ContentRequest, error) {
	requestId, err := uuid.Parse(os.Getenv("REQUEST_ID"))
	if err != nil {
		return nil, err
	}

	encodedBody := os.Getenv("REQUEST_BODY")
	if encodedBody == "" {
		return nil, fmt.Errorf("no request body found in REQUEST_BODY environment variable")
	}

	bodyJson, err := base64.StdEncoding.DecodeString(encodedBody)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 request found froun in REQUEST_BODY: %w", err)
	}

	var body map[string]any

	if err := json.Unmarshal(bodyJson, &body); err != nil {
		return nil, fmt.Errorf("failed to parse json body in REQUEST_BODY: %w", err)
	}

	return &models.ContentRequest{
		Id:   requestId,
		Body: body,
	}, nil
}

func publishDocument(ctx context.Context, doc *models.Document) (*url.URL, error) {
	apiToken := os.Getenv("TELEGRAPH_API_KEY")
	if apiToken == "" {
		return nil, fmt.Errorf("missing TELEGRAPH_API_TOKEN environment variable")
	}

	publisher := telegraph.NewDefaultClient()

	content := telegraph.Nodes{}

	for _, section := range doc.Sections {
		content = append(content, telegraph.NodeElement{
			Tag: "h3",
			Children: telegraph.Nodes{
				section.Title,
			},
		})

		for _, paragraph := range section.Paragraphs {
			content = append(content, telegraph.NodeElement{
				Tag: "p",
				Children: telegraph.Nodes{
					paragraph,
				},
			})
		}
	}

	returnContent := false

	pageRequest := telegraph.CreatePageRequest{
		AccessToken:   apiToken,
		Title:         doc.Title,
		AuthorName:    &doc.Author,
		Content:       content,
		ReturnContent: &returnContent,
	}

	page, err := publisher.CreatePage(ctx, pageRequest)
	if err != nil {
		return nil, err
	}

	pageURL, err := url.Parse(page.URL)
	if err != nil {
		return nil, err
	}

	return pageURL, nil
}
