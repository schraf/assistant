package interfaces

import (
	"context"
	"net/url"

	"github.com/schraf/assistant/pkg/models"
)

type JobScheduler interface {
	ScheduleJob(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error
}

type Publisher interface {
	PublishDocument(ctx context.Context, doc *models.Document) (*url.URL, error)
}

type Notifier interface {
	SendPublishedURLNotification(publishedURL *url.URL, title string) error
}
