package notifications

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailNotifierConfig struct {
	Login    string
	Password string
	Host     string
	Port     int
	From     string
	To       string
}

type EmailNotifier struct {
	cfg EmailNotifierConfig
}

func NewEmailNotifier(cfg EmailNotifierConfig) *EmailNotifier {
	return &EmailNotifier{cfg}
}

func (e *EmailNotifier) Notify(title string, message string, recipients ...string) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ",")))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n\r\n", title))
	builder.WriteString(fmt.Sprintf("%s\r\n", message))

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.Login, e.cfg.Password, e.cfg.Host)
	err := smtp.SendMail(addr, auth, e.cfg.From, recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailNotifier) Recipient() string {
	return e.cfg.To
}
