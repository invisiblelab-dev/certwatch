package certwatch

import (
	"fmt"
	"net/url"
)

type Domain struct {
	Name             string `yaml:"name"`
	NotificationDays int    `yaml:"days"`
}

func FormatDomain(domain string) (string, error) {
	url, err := url.Parse(domain)
	if err != nil {
		fmt.Printf("failed to parse url: %v\n", err)
		return "", err
	}

	if url.Scheme == "" {
		url.Scheme = "https"
	}

	return url.String(), nil
}
