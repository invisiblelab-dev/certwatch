package runners

import (
	"fmt"
	"time"

	certwatch "github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/config"

	"github.com/invisiblelab-dev/certwatch/notifications"
)

func getCertificates(domains []certwatch.Domain, refresh int) (map[string]certwatch.DomainQuery, error) {
	queries, err := config.ReadQueries()
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		domainLastCheck := queries[domain.Name].LastCheck
		timeSinceLastCheck := time.Since(domainLastCheck).Seconds()

		if int(timeSinceLastCheck) >= refresh {
			certificate, err := certwatch.Certificate(domain.Name)
			if err != nil {
				return nil, err
			}

			peerCertificate := certificate.PeerCertificates[0]
			queries[domain.Name] = certwatch.DomainQuery{
				Issuer:    peerCertificate.Issuer,
				LastCheck: time.Now(),
				NotAfter:  peerCertificate.NotAfter,
			}
		}
	}

	err = config.WriteQueries(queries)
	if err != nil {
		return nil, err
	}
	return queries, nil
}

func calculateDaysToDeadline(certificates map[string]certwatch.DomainQuery, configData certwatch.ConfigFile) []certwatch.DomainDeadline {
	domainsDeadlines := []certwatch.DomainDeadline{}
	for _, domain := range configData.Domains {
		timeHours := time.Until(certificates[domain.Name].NotAfter)
		timeDays := timeHours.Hours() / 24
		onDeadline := timeDays <= float64(domain.NotificationDays)

		deadline := certwatch.DomainDeadline{
			Domain:           domain.Name,
			DaysTillDeadline: timeDays,
			OnDeadline:       onDeadline,
		}
		domainsDeadlines = append(domainsDeadlines, deadline)
	}
	return domainsDeadlines
}

func RunCheckCertificatesCommand(opts certwatch.CheckCertificatesOptions) {
	for _, domain := range opts.Domains {
		fmt.Println("Domain:", domain)
		certificate, err := certwatch.Certificate(domain)
		if err != nil {
			continue
		}
		peerCertificate := certificate.PeerCertificates[0]
		fmt.Println(peerCertificate.String())
	}
}

func RunCheckAllCertificatesCommand(opts certwatch.CheckAllCertificatesOptions) {
	configData, err := config.ReadYaml(opts.Path)
	if err != nil {
		fmt.Println("could not read config file", err)
		return
	}

	certificates, err := getCertificates(configData.Domains, configData.Refresh)
	if err != nil {
		fmt.Println("failed to get certificates", err)
		return
	}

	domainDeadlines := calculateDaysToDeadline(certificates, configData)
	message, err := notifications.ComposeMessage(domainDeadlines)
	if err != nil {
		fmt.Println("failed to compose message:", err)
		return
	}

	err = notifications.SendEmail(message, configData.Notifications.Email)
	if err != nil {
		fmt.Println("failed to send email:", err)
		return
	}

	err = notifications.SendSlack(message, configData.Notifications.Slack)
	if err != nil {
		fmt.Println("failed to send slack:", err)
		return
	}
}