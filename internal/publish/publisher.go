package publish

import (
	"context"
	"net/url"

	"github.com/schraf/assistant/pkg/models"
)

type Publisher interface {
	PublishDocument(ctx context.Context, doc *models.Document) (*url.URL, error)
}
