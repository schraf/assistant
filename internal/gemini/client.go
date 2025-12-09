package gemini

import (
	"context"
	"encoding/json"
	"time"

	"github.com/schraf/assistant/internal/retry"
	"google.golang.org/genai"
)

type Client struct {
	genaiClient *genai.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	genaiClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		genaiClient: genaiClient,
	}, nil
}

func (c *Client) Ask(ctx context.Context, persona string, request string) (*string, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(persona, genai.RoleModel),
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
				URLContext:   &genai.URLContext{},
			},
		},
	}

	result, err := c.generateContext(ctx, request, config)
	if err != nil {
		return nil, err
	}

	responseText := result.Text()
	return &responseText, nil
}

func (c *Client) StructuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: schema,
		SystemInstruction:  genai.NewContentFromText(persona, genai.RoleModel),
	}

	result, err := c.generateContext(ctx, request, config)
	if err != nil {
		return nil, err
	}

	responseText := result.Text()

	var responseJson json.RawMessage

	err = json.Unmarshal([]byte(responseText), &responseJson)
	if err != nil {
		return nil, err
	}

	return responseJson, nil
}

func (c *Client) generateContext(ctx context.Context, request string, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	var result *genai.GenerateContentResponse

	model := modelFromContext(ctx)
	prompt := genai.Text(request)

	retryable := retry.Retryer{
		MaxRetries:       3,
		InitialBackoff:   1 * time.Second,
		MaxBackoff:       30 * time.Second,
		IsRetryableError: func(error) bool { return true },
		Attempt: func(ctx context.Context) error {
			var err error
			result, err = c.genaiClient.Models.GenerateContent(ctx, model, prompt, cfg)
			return err
		},
	}

	if err := retryable.Try(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

// WithModel returns a context with the specified model set.
func (c *Client) WithModel(ctx context.Context, model string) context.Context {
	return context.WithValue(ctx, modelKey, model)
}
