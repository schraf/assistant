package models

import (
	"context"
	"encoding/json"
)

type Assistant interface {
	Ask(ctx context.Context, persona string, request string) (*string, error)
	StructuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error)
}
