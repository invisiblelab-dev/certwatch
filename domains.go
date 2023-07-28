package certwatch

import (
	"fmt"
	"net/url"
)

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

	if scheme == "https" {
		fmtURL = fmtURL[2:]
	}

	return fmtURL, nil
}
