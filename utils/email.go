package utils

import (
	"fmt"

	"chatapp/env"

	"github.com/wneessen/go-mail"
)

func SendEmail(toEmail, subject, body string) error {
	message := mail.NewMsg()
	if err := message.From(env.EmailFrom); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}
	if err := message.To(toEmail); err != nil {
		return fmt.Errorf("failed to set TO address: %w", err)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, body)

	opts := []mail.Option{mail.WithTLSPortPolicy(mail.NoTLS)} // dev by default
	if env.EmailSmtpUser != "" && env.EmailSmtpPassword != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
			mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithUsername(env.EmailSmtpUser),
			mail.WithPassword(env.EmailSmtpPassword),
		)
	}

	client, err := mail.NewClient(env.EmailSmtpHost, opts...)
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %w", err)
	}
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to deliver mail: %w", err)
	}

	return nil
}
