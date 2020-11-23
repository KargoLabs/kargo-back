package urlValidator

import (
	"net/url"
	"strings"
)

func IsValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return true
}

func IsValidURLOfGivenDomain(rawURL, domain string) bool {
	return strings.HasPrefix(rawURL, domain) && IsValidURL(rawURL)
}
