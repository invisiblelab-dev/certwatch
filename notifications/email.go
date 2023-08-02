package notifications

import (
	"bytes"
	"embed"
	"fmt"
	"net/smtp"
	"path"
	"strings"
	"text/template"

	"github.com/invisiblelab-dev/certwatch"
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

	html := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(html, "html", data); err != nil {
		return fmt.Errorf("failed to hydrate html email template: %w", err)
	}

	hex, err := certwatch.RandomHex(16)
	if err != nil {
		return fmt.Errorf("failed to generate boundary: %w", err)
	}

	boundary := fmt.Sprintf("_%s", hex)
	subject := title
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("From: %s\n", e.cfg.From))
	builder.WriteString(fmt.Sprintf("To: %s\n", strings.Join(recipients, ",")))
	builder.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	builder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\n\n", boundary))
	builder.WriteString(fmt.Sprintf("--%s\n", boundary))
	builder.WriteString(fmt.Sprintf("Content-Type: text/plain; charset=%q\n", "utf-8"))
	builder.WriteString("Content-Transfer-Encoding: quoted-printable\n")
	builder.WriteString("Content-Disposition: inline\n\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", strings.Join(data.Messages, "\n")))
	builder.WriteString(fmt.Sprintf("--%s\n", boundary))
	builder.WriteString(fmt.Sprintf("Content-Type: text/html; charset=%q\n", "utf-8"))
	builder.WriteString("Content-Transfer-Encoding: quoted-printable\n")
	builder.WriteString("Content-Disposition: inline\n\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", html.String()))
	builder.WriteString(fmt.Sprintf("--%s--\n", boundary))

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.Login, e.cfg.Password, e.cfg.Host)
	err = smtp.SendMail(addr, auth, e.cfg.From, recipients, []byte(builder.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailNotifier) Recipient() string {
	return e.cfg.To
}
