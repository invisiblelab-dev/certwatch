package notifications

import (
	"fmt"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
)

func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) string {
	var subject string
	for i := 0; i < len(domainDeadlines); i++ {
		if domainDeadlines[i].OnDeadline {
			subject = subject +
				fmt.Sprintf("\n\n %s certificate expires in %d days.",
					domainDeadlines[i].Domain,
					int64(domainDeadlines[i].DaysTillDeadline))
		}
	}
	return subject
}
