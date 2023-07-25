package notifications

import (
	"log"
	"net/smtp"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
)

func SendEmail(subject string, config certwatch.ConfigFile) (bool, error) {
	// Mailtrap account config
	username := config.Notifications.Email.Mailtrap.Username
	password := config.Notifications.Email.Mailtrap.Password
	smtpHost := config.Notifications.Email.Mailtrap.SmtpHost // TEST email mailtrap host

	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Message data
	from := config.Notifications.Email.From
	to := []string{config.Notifications.Email.To}
	email := "To: " + to[0] + "\n\n" +
		"From: " + from + "\n\n" +
		"Subject: " + subject

	message := []byte(email)

	// Connect to the server and send message
	smtpUrl := smtpHost + ":465"

	err := smtp.SendMail(smtpUrl, auth, from, to, message)

	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}
