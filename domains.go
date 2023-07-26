package certwatch

import (
	"fmt"
	"net/url"
)

type Domain struct {
	Name             string `yaml:"name"`
	NotificationDays int    `yaml:"days"`
}

func AddHttps(domain string) (string, error) {
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

func RemoveHttps(domain string) (string, error) {
	url, err := url.Parse(domain)
	if err != nil {
		fmt.Printf("failed to parse url: %v\n", err)
		return "", err
	}
	scheme := url.Scheme
	url.Scheme = ""
	fmtUrl := url.String()

	if scheme == "https" {
		fmtUrl = fmtUrl[2:]
	}
	return fmtUrl, nil
}
