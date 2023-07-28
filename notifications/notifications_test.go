package notifications

import (
	"fmt"
	"testing"

	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/test/helpers"
)

func TestComposeMessage(t *testing.T) {
	domainDeadlines := []certwatch.DomainDeadline{
		{Domain: "https://test1.com", DaysTillDeadline: 1, OnDeadline: true},
		{Domain: "https://test2.com", DaysTillDeadline: 1, OnDeadline: false},
		{Domain: "https://test3.com", DaysTillDeadline: -1, OnDeadline: true},
	}
	subject, err := ComposeMessage(domainDeadlines)
	if err != nil {
		fmt.Println("Not able to Compose, err: ", err)
		return
	}

	helpers.Equal(t, subject, "\n\n test1.com certificate expires in 1 days.\n\n test3.com certificate has expired 1 days ago.")
}
