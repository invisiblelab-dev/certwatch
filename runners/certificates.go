package runners

import (
	"fmt"
	"os"
	"time"

	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/config"
	"github.com/invisiblelab-dev/certwatch/notifications"
)

func getCertificates(domains []certwatch.Domain, refresh int) (map[string]certwatch.DomainQuery, error) {
	queries, err := config.ReadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to get queries: %w", err)
	}

	for _, domain := range domains {
		domainLastCheck := queries[domain.Name].LastCheck
		timeSinceLastCheck := time.Since(domainLastCheck).Seconds()

		if int(timeSinceLastCheck) >= refresh {
			certificate, err := certwatch.Certificate(domain.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch certificate: %w", err)
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
		return nil, fmt.Errorf("failed to write queries: %w", err)
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

// nolint: forbidigo
func RunCheckCertificatesCommand(opts certwatch.CheckCertificatesOptions) error {
	for _, domain := range opts.Domains {
		fmt.Println("Domain:", domain)
		certificate, err := certwatch.Certificate(domain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch domain [%s]: %v", domain, err)
			continue
		}
		peerCertificate := certificate.PeerCertificates[0]
		fmt.Println(peerCertificate.String())
	}

	return nil
}

func RunCheckAllCertificatesCommand(opts certwatch.CheckAllCertificatesOptions) error {
	configData, err := config.ReadYaml(opts.Path)
	if err != nil {
		return fmt.Errorf("could not read config file [%s]: %w", opts.Path, err)
	}

	certificates, err := getCertificates(configData.Domains, configData.Refresh)
	if err != nil {
		return fmt.Errorf("failed to get certificates: %w", err)
	}

	domainDeadlines := calculateDaysToDeadline(certificates, configData)
	message, err := notifications.ComposeMessage(domainDeadlines)
	if err != nil {
		return fmt.Errorf("failed to compose message: %w", err)
	}

	err = notifications.SendEmail(message, configData.Notifications.Email)
	if err != nil {
		return fmt.Errorf("email error: %w", err)
	}

	err = notifications.SendSlack(message, configData.Notifications.Slack)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}

	return nil
}
