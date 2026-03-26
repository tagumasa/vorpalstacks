// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import "net/http"

// GetFullURL returns the full URL of the request.
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - string: The full URL
//
// Example:
//
//	url := GetFullURL(r)
func GetFullURL(r *http.Request) string {
	scheme := GetScheme(r)
	host := GetHost(r)
	return scheme + "://" + host + r.RequestURI
}
