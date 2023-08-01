package notifications

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"strings"
	"text/template"
	"time"

	"gopkg.in/mail.v2"
)

//go:embed "templates"
var templateFS embed.FS

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

func (e *EmailNotifier) Notify(title string, data MessageData, recipients ...string) error {
	tmpl, err := template.New("email").ParseFS(templateFS, path.Join("templates", "notification.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	subject := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return fmt.Errorf("failed to hydrate subject email template: %w", err)
	}

	plaintext := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(plaintext, "plain", data); err != nil {
		return fmt.Errorf("failed to hydrate plaintext email template: %w", err)
	}

	html := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(html, "html", data); err != nil {
		return fmt.Errorf("failed to hydrate html email template: %w", err)
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", strings.Join(recipients, ","))
	msg.SetHeader("From", e.cfg.From)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plaintext.String())
	msg.AddAlternative("text/html", html.String())

	dialer := mail.NewDialer(e.cfg.Host, e.cfg.Port, e.cfg.Login, e.cfg.Password)
	dialer.Timeout = 5 * time.Second

	if err = dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (e *EmailNotifier) Recipient() string {
	return e.cfg.To
}
