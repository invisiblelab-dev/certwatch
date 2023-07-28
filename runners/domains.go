package runners

import (
	"errors"
	"fmt"

	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/config"
	"gopkg.in/yaml.v3"
)

func AddDomain(domain string, daysToNotify int, path string) error {
	domains, err := config.ReadYaml(path)
	if err != nil {
		return fmt.Errorf("error while reading yaml file %s: %w", path, err)
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
		return fmt.Errorf("error while marshalling yaml file %w", err)
	}

	if err := config.WriteYaml(marshalData, path); err != nil {
		return fmt.Errorf("failed to add domain: %w", err)
	}

	return nil
}

func RunAddDomainCommand(opts certwatch.AddDomainOptions) error {
	url, err := certwatch.AddHTTPS(opts.Domain)
	if err != nil {
		return fmt.Errorf("error parsing domain: %w", err)
	}

	if opts.DaysBefore <= 0 {
		return fmt.Errorf("days can't be <= 0: %d", opts.DaysBefore)
	}

	err = AddDomain(url, int(opts.DaysBefore), opts.Path)
	if err != nil {
		return fmt.Errorf("failed to add url: %w", err)
	}

	return nil
}
