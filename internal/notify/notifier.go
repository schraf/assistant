package notify

import "net/url"

type Notifier interface {
	SendPublishedURLNotification(publishedURL *url.URL, title string) error
}
