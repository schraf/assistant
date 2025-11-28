package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(host string, port int, username, password, to, from, subject, body string) error {
	message := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" + // This blank line separates headers from the body
			body + "\n",
	)

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	err := smtp.SendMail(
		addr,
		auth,
		from,
		[]string{to},
		message,
	)

	if err != nil {
		return err
	}

	return nil
}
