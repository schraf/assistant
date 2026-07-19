package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Format   string    `json:"format,omitempty"`
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func NewClient(ctx context.Context) (*Client, error) {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) Ask(ctx context.Context, persona string, request string) (*string, error) {
	model := modelFromContext(ctx)

	messages := []message{
		{
			Role:    "system",
			Content: persona,
		},
		{
			Role:    "user",
			Content: request,
		},
	}

	chatReq := chatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	responseText := strings.TrimSpace(chatResp.Message.Content)
	return &responseText, nil
}

func (c *Client) StructuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error) {
	model := modelFromContext(ctx)

	// Convert schema to JSON string for the format parameter
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Enhance persona with JSON format instruction
	enhancedPersona := persona + "\n\nYou must respond with valid JSON that matches the following schema:\n" + string(schemaJSON)

	messages := []message{
		{
			Role:    "system",
			Content: enhancedPersona,
		},
		{
			Role:    "user",
			Content: request,
		},
	}

	chatReq := chatRequest{
		Model:    model,
		Messages: messages,
		Format:   "json",
		Stream:   false,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	responseText := strings.TrimSpace(chatResp.Message.Content)

	// Remove markdown code blocks if present
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	var responseJSON json.RawMessage
	if err := json.Unmarshal([]byte(responseText), &responseJSON); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return responseJSON, nil
}

// WithModel returns a context with the specified model set.
func (c *Client) WithModel(ctx context.Context, model string) context.Context {
	return context.WithValue(ctx, modelKey, model)
}
