package runners

import (
	"errors"
	"fmt"

	certwatch "github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/config"
	"gopkg.in/yaml.v3"
)

func AddDomain(domain string, daysToNotify int, path string) error {
	domains, err := config.ReadYaml(path)
	if err != nil {
		return err
	}
	newDomain := certwatch.Domain{Name: domain, NotificationDays: daysToNotify}

	for _, listedDomain := range domains.Domains {
		if listedDomain.Name == domain {
			return errors.New("domain already added")
		}
	}

	domains.Domains = append(domains.Domains, newDomain)
	marshalData, err := yaml.Marshal(&domains)
	if err != nil {
		fmt.Println("error marshalling file: ", err)
		return err
	}
	return config.WriteYaml(marshalData)
}

func RunAddDomainCommand(opts certwatch.AddDomainOptions) {
	url, err := certwatch.AddHttps(opts.Domain)
	if err != nil {
		fmt.Println("error parsing domain:", err)
	}

	if opts.DaysBefore <= 0 {
		fmt.Printf("Days cant be <=0: %v\n", opts.DaysBefore)
		return
	}

	err = AddDomain(url, int(opts.DaysBefore), opts.Path)
	if err != nil {
		fmt.Printf("failed to add url: %v\n", err)
		return
	}
}
