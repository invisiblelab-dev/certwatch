package notifications

import (
	"log"
	"net/smtp"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
)

func SendEmail(subject string, config certwatch.ConfigFile) (bool, error) {
	username := config.Notifications.Email.Username
	password := config.Notifications.Email.Password
	smtpHost := config.Notifications.Email.SmtpHost

	auth := smtp.PlainAuth("", username, password, smtpHost)

	from := config.Notifications.Email.From
	to := []string{config.Notifications.Email.To}
	email := "To: " + to[0] + "\n\n" +
		"From: " + from + "\n\n" +
		"Subject: " + subject

	message := []byte(email)
	smtpUrl := smtpHost + ":465"

	err := smtp.SendMail(smtpUrl, auth, from, to, message)

	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}
