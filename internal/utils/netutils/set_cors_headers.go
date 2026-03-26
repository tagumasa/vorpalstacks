// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import (
	"fmt"
	"net/http"
)

// SetCORSHeaders sets CORS headers on the response writer.
//
// Parameters:
//   - w: The HTTP response writer
//   - allowedOrigins: The allowed origins (comma-separated)
//   - allowedMethods: The allowed methods (comma-separated)
//   - allowedHeaders: The allowed headers (comma-separated)
//   - allowCredentials: Whether to allow credentials (default: false)
//
// Returns:
//   - error: An error if credentials are enabled with wildcard origin
//
// Example:
//
//	SetCORSHeaders(w, "*", "GET,POST", "Content-Type", false)
func SetCORSHeaders(w http.ResponseWriter, allowedOrigins, allowedMethods, allowedHeaders string, allowCredentials bool) error {
	if allowCredentials && allowedOrigins == "*" {
		return fmt.Errorf("cannot use wildcard origin with credentials: specify exact origin")
	}
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
	w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	if allowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	return nil
}
