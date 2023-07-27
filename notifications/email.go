package notifications

import (
	"log"
	"net/smtp"
	"strconv"

	certwatch "github.com/invisiblelab-dev/certwatch"
)

func SendEmail(subject string, emailConfig certwatch.Email) error {
	username := emailConfig.Username
	password := emailConfig.Password
	smtpHost := emailConfig.SmtpHost
	port := emailConfig.Port

	auth := smtp.PlainAuth("", username, password, smtpHost)

	from := emailConfig.From
	to := []string{emailConfig.To}
	email := "To: " + to[0] + "\n\n" +
		"From: " + from + "\n\n" +
		"Subject: " + subject

	message := []byte(email)
	smtpUrl := smtpHost + ":" + strconv.Itoa(port)
	err := smtp.SendMail(smtpUrl, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
