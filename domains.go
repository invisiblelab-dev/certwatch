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
