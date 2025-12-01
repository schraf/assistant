package job

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/schraf/assistant/internal/interfaces"
	"github.com/schraf/assistant/pkg/generators"
	"github.com/schraf/assistant/pkg/models"
)

// Processor handles the job processing workflow.
type Processor struct {
	assistant models.Assistant
	publisher interfaces.Publisher
	notifier  interfaces.Notifier
	logger    *slog.Logger
}

// NewProcessor creates a new Processor with the given dependencies.
func NewProcessor(assistant models.Assistant, publisher interfaces.Publisher, notifier interfaces.Notifier, logger *slog.Logger) *Processor {
	return &Processor{
		assistant: assistant,
		publisher: publisher,
		notifier:  notifier,
		logger:    logger,
	}
}

// Process executes the complete job workflow:
// 1. Get request and config from environment
// 2. Generate content
// 3. Publish document
// 4. Send notification
func (p *Processor) Process(ctx context.Context) error {
	//--========================================================================--
	//--== GET THE REQUEST
	//--========================================================================--

	request, err := getRequest()
	if err != nil {
		p.logger.ErrorContext(ctx, "invalid_request",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("invalid request: %w", err)
	}

	logger := p.logger.With(
		slog.String("request_id", request.Id.String()),
	)

	logger.InfoContext(ctx, "request_body",
		slog.Any("body", request.Body),
	)

	//--========================================================================--
	//--== GET THE CONFIG
	//--========================================================================--

	config, err := getConfig()
	if err != nil {
		logger.ErrorContext(ctx, "failed_getting_config",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed getting config: %w", err)
	}

	//--========================================================================--
	//--== APPLY MODEL SELECTION
	//--========================================================================--

	if model, ok := (*config)["model"].(string); ok {
		var modelName string
		if model == "pro" {
			modelName = "gemini-pro-latest"
		} else if model == "basic" {
			modelName = "gemini-flash-latest"
		} else {
			// For other assistants (like Ollama), use the model name directly
			modelName = model
		}

		ctx = p.assistant.WithModel(ctx, modelName)

		logger.InfoContext(ctx, "using_model",
			slog.String("model", modelName),
		)
	}

	//--========================================================================--
	//--== GENERATE CONTENT
	//--========================================================================--

	logger.InfoContext(ctx, "generating_content")

	contentGenerator, err := getContentGenerator(*config)
	if err != nil {
		logger.ErrorContext(ctx, "failed_creating_generator",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed creating generator: %w", err)
	}

	doc, err := contentGenerator.Generate(ctx, *request, p.assistant)
	if err != nil {
		logger.ErrorContext(ctx, "content_generation_error",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("content generation error: %w", err)
	}

	logger.InfoContext(ctx, "generated_document",
		slog.String("title", doc.Title),
		slog.String("author", doc.Author),
		slog.Int("section_count", len(doc.Sections)),
	)

	//--========================================================================--
	//--== PUBLISH CONTENT
	//--========================================================================--

	logger.InfoContext(ctx, "publishing_document")

	url, err := p.publisher.PublishDocument(ctx, doc)
	if err != nil {
		logger.ErrorContext(ctx, "publish_error",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("publish error: %w", err)
	}

	logger.InfoContext(ctx, "published_document",
		slog.String("url", url.String()),
	)

	//--========================================================================--
	//--== SEND NOTIFICATION
	//--========================================================================--

	logger.InfoContext(ctx, "sending_notification")

	if err := p.notifier.SendPublishedURLNotification(url, doc.Title); err != nil {
		logger.ErrorContext(ctx, "failed_sending_notification",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed sending notification: %w", err)
	}

	logger.InfoContext(ctx, "notification_sent")
	logger.InfoContext(ctx, "job_completed_successfully")

	return nil
}

func getConfig() (*generators.Config, error) {
	encodedConfig := os.Getenv("CONTENT_CONFIG")
	if encodedConfig == "" {
		return &generators.Config{}, nil
	}

	configJson, err := base64.StdEncoding.DecodeString(encodedConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 config found in CONTENT_CONFIG: %w", err)
	}

	var config generators.Config

	if err := json.Unmarshal(configJson, &config); err != nil {
		return nil, fmt.Errorf("failed to parse json body in CONTENT_CONFIG: %w", err)
	}

	return &config, nil
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
		return nil, fmt.Errorf("invalid base64 request found in REQUEST_BODY: %w", err)
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

func getContentGenerator(config generators.Config) (models.ContentGenerator, error) {
	contentType := os.Getenv("CONTENT_TYPE")
	if contentType == "" {
		return nil, fmt.Errorf("no content type found in CONTENT_TYPE environment variable")
	}

	return generators.Create(contentType, config)
}
