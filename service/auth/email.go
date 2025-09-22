package auth

import (
	"chatapp/config"
	"fmt"

	"github.com/wneessen/go-mail"
)

// TODO: create utils.Mailer interface to use in different services across the app.
func sendEmail(toEmail, subject, body string) error {
	message := mail.NewMsg()
	if err := message.From(config.EmailFrom); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}
	if err := message.To(toEmail); err != nil {
		return fmt.Errorf("failed to set TO address: %w", err)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, body)

	client, err := mail.NewClient(config.PapercutSmtpHost, mail.WithTLSPortPolicy(mail.NoTLS))
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %w", err)
	}
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to deliver mail: %w", err)
	}

	return nil
}
