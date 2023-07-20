package runners

import (
	"time"
)

type DomainDeadline struct {
	Domain           string
	DaysTillDeadline float64
}

type DomainsDeadlines struct {
	Deadlines []DomainDeadline
}

func CalculateDeadline(certificates SSLInfoArray) DomainsDeadlines {
	domainsDeadlines := DomainsDeadlines{}
	for i := 0; i < len(certificates.DomainsSSLs); i++ {
		timeHours := time.Until(certificates.DomainsSSLs[i].SSL.PeerCertificates[0].NotAfter)
		timeDays := timeHours.Hours() / 24
		deadline := DomainDeadline{Domain: certificates.DomainsSSLs[i].Domain, DaysTillDeadline: timeDays}
		domainsDeadlines.Deadlines = append(domainsDeadlines.Deadlines, deadline)
	}
	return domainsDeadlines
}
