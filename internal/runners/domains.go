package runners

import (
	"fmt"
	"net/url"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
)

func RunAddDomainCommand(opts certwatch.AddDomainOptions) {
	_, err := url.ParseRequestURI(opts.Domain)
	if err != nil {
		fmt.Printf("failed to read url: %v\n", err)
		return
	}

	fmt.Println("domain", opts.Domain, opts.DaysBefore)
}
