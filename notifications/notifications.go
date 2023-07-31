package notifications

import (
	"fmt"
	"math"

	certwatch "github.com/invisiblelab-dev/certwatch"
)

// TODO: later review logic
func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) (string, error) {
	var subject string
	for _, domainDeadline := range domainDeadlines {
		if domainDeadline.OnDeadline {
			domain, err := certwatch.RemoveHTTPS(domainDeadline.Domain)
			if err != nil {
				return "", fmt.Errorf("failed to remove HTTPS form URL: %w", err)
			}
			var message string
			if domainDeadline.DaysTillDeadline <= 0 {
				message =
					fmt.Sprintf("\n\n %s certificate has expired %d days ago.",
						domain,
						int64(math.Abs(domainDeadline.DaysTillDeadline)))
			} else {
				message =
					fmt.Sprintf("\n\n %s certificate expires in %d days.",
						domain,
						int64(domainDeadline.DaysTillDeadline))
			}
			subject += message
		}
	}

	return subject, nil
}
