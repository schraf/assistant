package mocks

import "net/url"

// MockNotifier is a mock implementation of models.Notifier.
type MockNotifier struct {
	SendPublishedURLNotificationFunc func(publishedURL *url.URL, title string) error
}

// SendPublishedURLNotification calls SendPublishedURLNotificationFunc if set, otherwise returns nil.
func (m *MockNotifier) SendPublishedURLNotification(publishedURL *url.URL, title string) error {
	if m.SendPublishedURLNotificationFunc != nil {
		return m.SendPublishedURLNotificationFunc(publishedURL, title)
	}

	return nil
}
