package certwatch

import (
	"fmt"
	"net/url"
)

type Domain struct {
	Name             string `yaml:"name"`
	NotificationDays int    `yaml:"days"`
}

// TODO: review logic later
func AddHTTPS(domain string) (string, error) {
	url, err := url.Parse(domain)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	if url.Scheme == "" {
		url.Scheme = "https"
	}

	return url.String(), nil
}

func RemoveHTTPS(domain string) (string, error) {
	url, err := url.Parse(domain)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	scheme := url.Scheme
	url.Scheme = ""
	fmtURL := url.String()

	if scheme[0:4] == "http" {
		fmtURL = fmtURL[2:]
	}

	return fmtURL, nil
}
