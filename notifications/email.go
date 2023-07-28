package notifications

import (
	"fmt"
	"net/smtp"
	"strconv"

	certwatch "github.com/invisiblelab-dev/certwatch"
)

func SendEmail(subject string, emailConfig certwatch.Email) error {
	username := emailConfig.Username
	password := emailConfig.Password
	smtpHost := emailConfig.SMTPHost
	port := emailConfig.Port

	auth := smtp.PlainAuth("", username, password, smtpHost)

	from := emailConfig.From
	to := []string{emailConfig.To}
	email := "To: " + to[0] + "\n\n" +
		"From: " + from + "\n\n" +
		"Subject: " + subject

	message := []byte(email)
	smtpURL := smtpHost + ":" + strconv.Itoa(port)

	err := smtp.SendMail(smtpURL, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
