package gateway

import (
	"net/http"
	"strings"
)

// Extract the subdomain from request.
func extractSubdomainFromRequest(r *http.Request) string {
	host := strings.TrimSpace(r.Host)
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}

	var subdomain string
	host_parts := strings.Split(host, ".")
	if len(host_parts) > 2 {
		subdomain = host_parts[0]
	}
	return subdomain
}
