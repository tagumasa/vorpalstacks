// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import "net/http"

// GetScheme extracts the URL scheme (http or https) from an HTTP request.
// It checks the X-Forwarded-Proto header if the request comes from a trusted proxy.
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - string: The scheme ("http" or "https")
//
// Example:
//
//	scheme := GetScheme(r)
func GetScheme(r *http.Request) string {
	forwardedProto := r.Header.Get("X-Forwarded-Proto")
	if forwardedProto != "" && isTrustedProxy(r) {
		return forwardedProto
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
