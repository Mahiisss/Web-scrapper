package services

import (
	"net/url"
	"strings"
)

func Clean(host string, href string) string {
	// Make sure host ends with slash
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}

	// Parse href URL
	parsedHref, err := url.Parse(href)
	if err == nil && (parsedHref.Scheme == "http" || parsedHref.Scheme == "https") {
		// If href is an absolute URL (http/https), return empty string
		return ""
	}

	// use it with host if the href is relative URL
	if parsedHref != nil {
		host = strings.TrimSuffix(host, "/top10/")
		host += "/"
		return host + strings.TrimPrefix(parsedHref.Path, "/")
	}
	return host
}
