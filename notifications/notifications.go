package notifications

import (
	"fmt"
	"math"

	certwatch "github.com/invisiblelab-dev/certwatch"
)

func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) (string, error) {
	var subject string
	for _, domainDeadline := range domainDeadlines {
		if domainDeadline.OnDeadline {
			domain, err := certwatch.RemoveHttps(domainDeadline.Domain)
			if err != nil {
				return "", err
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
			subject = subject + message
		}
	}
	return subject, nil
}
