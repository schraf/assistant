package notify

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	internal_models "github.com/schraf/assistant/internal/models"
)

// EmailNotifier implements internal_models.Notifier using SMTP email.
type EmailNotifier struct{}

// NewEmailNotifier creates a new EmailNotifier.
func NewEmailNotifier() internal_models.Notifier {
	return &EmailNotifier{}
}

// SendPublishedURLNotification sends an email notification with the published URL
// using environment variables defined in terraform/job.tf
func (n *EmailNotifier) SendPublishedURLNotification(publishedURL *url.URL, title string) error {
	host := os.Getenv("MAIL_SMTP_SERVER")
	if host == "" {
		return fmt.Errorf("missing MAIL_SMTP_SERVER environment variable")
	}

	portStr := os.Getenv("MAIL_SMTP_PORT")
	if portStr == "" {
		return fmt.Errorf("missing MAIL_SMTP_PORT environment variable")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid MAIL_SMTP_PORT: %w", err)
	}

	from := os.Getenv("MAIL_SENDER_EMAIL")
	if from == "" {
		return fmt.Errorf("missing MAIL_SENDER_EMAIL environment variable")
	}

	password := os.Getenv("MAIL_SENDER_PASSWORD")
	if password == "" {
		return fmt.Errorf("missing MAIL_SENDER_PASSWORD environment variable")
	}

	to := os.Getenv("MAIL_RECIPIENT_EMAIL")
	if to == "" {
		return fmt.Errorf("missing MAIL_RECIPIENT_EMAIL environment variable")
	}

	subject := title
	body := publishedURL.String()

	return SendEmail(host, port, from, password, to, from, subject, body)
}
