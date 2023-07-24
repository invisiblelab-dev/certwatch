package notifications

import (
	"errors"
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

func EmailNotification() error {
	certificates := runners.GetCertificates()
	deadlines := runners.CalculateDaysToDeadline(certificates)
	domains := runners.ReadYaml()
	for i, domain := range deadlines.Deadlines {
		if domain.Domain != domains.Domains[i].Name {
			return errors.New("domains don't match")
		}

		if domain.DaysTillDeadline <= float64(domains.Domains[i].NotificationDays) {
			subject := domain.Domain + " certificate expires in " + fmt.Sprintf("%f", domain.DaysTillDeadline) + " days."
			sendEmail(subject, domains)
			fmt.Println("Email sent for domain " + domain.Domain)
		} else {
			continue
		}
	}
	return nil
}
