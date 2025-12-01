package models

import (
	"context"
	"net/url"

	pkgmodels "github.com/schraf/assistant/pkg/models"
)

// JobScheduler defines the interface for scheduling Cloud Run jobs.
type JobScheduler interface {
	ScheduleJob(ctx context.Context, contentType string, config map[string]any, request pkgmodels.ContentRequest) error
}

// Publisher defines the interface for publishing documents.
type Publisher interface {
	PublishDocument(ctx context.Context, doc *pkgmodels.Document) (*url.URL, error)
}

// Notifier defines the interface for sending notifications.
type Notifier interface {
	SendPublishedURLNotification(publishedURL *url.URL, title string) error
}
