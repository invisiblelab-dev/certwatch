package notifications

import (
	"fmt"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
)

func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) (string, error) {
	var subject string
	for i := 0; i < len(domainDeadlines); i++ {
		if domainDeadlines[i].OnDeadline {
			domain, err := certwatch.RemoveHttps(domainDeadlines[i].Domain)
			if err != nil {
				return "", err
			}
			subject = subject +
				fmt.Sprintf("\n\n %s certificate expires in %d days.",
					domain,
					int64(domainDeadlines[i].DaysTillDeadline))
		}
	}
	return subject, nil
}
