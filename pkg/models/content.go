package models

import (
	"context"

	"github.com/google/uuid"
)

type ContentRequest struct {
	Id   uuid.UUID
	Body map[string]any
}

type ContentGenerator interface {
	Generate(ctx context.Context, request ContentRequest, assistant Assistant) (*Document, error)
}
