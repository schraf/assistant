package models

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrContentBlocked = errors.New("content blocked")
)

type Assistant interface {
	Ask(ctx context.Context, persona string, request string) (*string, error)
	StructuredAsk(ctx context.Context, persona string, request string, schema map[string]any) (json.RawMessage, error)
	WithModel(ctx context.Context, model string) context.Context
}
