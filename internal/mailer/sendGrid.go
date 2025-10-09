package mailer

import (
	"context"
	"time"

	mailgun "github.com/mailgun/mailgun-go/v4"
)

type MailGunMailer struct {
	FromEmail string
	APIKey    string
	Client    *mailgun.MailgunImpl
}

func NewSendGrind(apiKey, from string) *MailGunMailer {
	client := mailgun.NewMailgun("sandbox4b7c75e350f94c55b3e2b4d065bb126b.mailgun.org", apiKey)
	return &MailGunMailer{
		FromEmail: "test@sandbox4b7c75e350f94c55b3e2b4d065bb126b.mailgun.org",
		APIKey:    apiKey,
		Client:    client,
	}
}

func (s *MailGunMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	message := mailgun.NewMessage(s.FromEmail, "", "", email)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err := s.Client.Send(ctx, message)
	if err != nil {
		return err

	}
	return nil
}
