package notifications

import (
	"fmt"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
)

func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) (string, error) {
	var subject string
	for i := 0; i < len(domainDeadlines); i++ {
		if domainDeadlines[i].OnDeadline {
			subject = subject + "\n\n" + domainDeadlines[i].Domain + "certificate expires in " + fmt.Sprintf("%f", domainDeadlines[i].DaysTillDeadline) + " days."
		} else {
			continue
		}
	}
	return subject, nil
}
