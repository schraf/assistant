package log

import (
	"log/slog"
	"net/url"
)

type LogNotifier struct{}

func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

func (n *LogNotifier) SendPublishedURLNotification(publishedURL *url.URL, title string) error {
	slog.Info("notification",
		slog.String("title", title),
		slog.String("url", publishedURL.String()),
	)

	return nil
}
