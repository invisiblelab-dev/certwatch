package notifications

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"

	"github.com/invisiblelab-dev/certwatch/internal/runners"
)

func sendEmail(subject string) (bool, error) {
	// Mailtrap account config
	username := "411d906e45a1ed"
	smtpHost := "sandbox.smtp.mailtrap.io" // TEST email mailtrap host
	password := "72cb4b151d2841"

	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Message data

	from := "john.doe@your.domain"

	to := []string{"kate.doe@example.com"}
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
	domains := runners.ReadDomains().Domains
	for i, domain := range deadlines.Deadlines {
		if domain.Domain != domains[i].Name {
			return errors.New("domains don't match")
		}

		if domain.DaysTillDeadline <= float64(domains[i].NotificationDays) {
			subject := domain.Domain + " certificate expires in " + fmt.Sprintf("%f", domain.DaysTillDeadline) + " days."
			sendEmail(subject)
			fmt.Println("Email sent for domain " + domain.Domain)
		} else {
			continue
		}
	}
	return nil
}
