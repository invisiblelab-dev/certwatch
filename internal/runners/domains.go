package runners

import (
	"errors"
	"fmt"
	"net/url"

	certwatch "github.com/invisiblelab-dev/certwatch/internal"
	"github.com/invisiblelab-dev/certwatch/internal/config"
	"gopkg.in/yaml.v3"
)

func AddDomain(domain string, daysToNotify int) error {
	domains := config.ReadYaml()
	newDomain := certwatch.Domain{Name: domain, NotificationDays: daysToNotify}

	for _, listedDomain := range domains.Domains {
		if listedDomain.Name == domain {
			return errors.New("domain already added")
		} else {
			continue
		}
	}

	domains.Domains = append(domains.Domains, newDomain)
	marshalData, err := yaml.Marshal(&domains)
	if err != nil {
		fmt.Println("not marshalling file, error: ", err)
		return err
	}
	return config.WriteYaml(marshalData)
}

func RunAddDomainCommand(opts certwatch.AddDomainOptions) {
	url, err := url.Parse(opts.Domain)
	if err != nil {
		fmt.Printf("failed to parse url: %v\n", err)
		return
	}
	if url.Scheme != "https" {
		fmt.Printf("url is not https: %v\n", opts.Domain)
		return
	}
	if opts.DaysBefore <= 0 {
		fmt.Printf("Days cant be <=0: %v\n", opts.DaysBefore)
		return
	}

	err = AddDomain(opts.Domain, int(opts.DaysBefore))
	if err != nil {
		fmt.Printf("failed to add url: %v\n", err)
		return
	}
}
