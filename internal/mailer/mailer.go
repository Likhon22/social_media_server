package mailer

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
type Mailer struct {
	Client
}

func NewClient(apiKey, from string) *Mailer {
	return &Mailer{
		Client: &MailGunMailer{
			FromEmail: from,
			APIKey:    apiKey,
		},
	}
}
