package notifications

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/invisiblelab-dev/certwatch/internal/runners"
)

func sendEmail(subject string, config runners.ConfigFile) (bool, error) {
	// Mailtrap account config

	username := config.Notifications.Email.Mailtrap.Username
	password := config.Notifications.Email.Mailtrap.Password
	smtpHost := "sandbox.smtp.mailtrap.io" // TEST email mailtrap host

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

func CheckAndNotify() error {
	message, err := composeMessage()
	if err != nil {
		return err
	}

	_, err = sendEmail(message, runners.ReadYaml())
	if err != nil {
		return err
	}
	return nil
}

func composeMessage() (string, error) {
	certificates := runners.GetCertificates()
	deadlines, err := runners.CalculateDaysToDeadline(certificates)
	if err != nil {
		return "", err
	}

	var subject string
	for i := 0; i < len(deadlines.Deadlines); i++ {
		if deadlines.Deadlines[i].OnDeadline {
			subject = subject + "\n\n" + deadlines.Deadlines[i].Domain + "certificate expires in " + fmt.Sprintf("%f", deadlines.Deadlines[i].DaysTillDeadline) + " days."
		} else {
			continue
		}
	}
	return subject, nil
}
