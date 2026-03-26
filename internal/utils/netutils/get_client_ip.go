// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP extracts the client IP address from an HTTP request.
// It checks the X-Forwarded-For, X-Real-IP, and RemoteAddr headers in that order.
// Only trusts X-Forwarded-For and X-Real-IP headers if the request comes from
// a trusted proxy (as configured via SetTrustedProxies).
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - string: The client IP address
//
// Example:
//
//	ip := GetClientIP(r)
func GetClientIP(r *http.Request) string {
	// Only trust proxy headers if from trusted proxy
	trusted := isTrustedProxy(r)

	// Check X-Forwarded-For header (only if from trusted proxy)
	if trusted {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// X-Forwarded-For can contain multiple IPs, take the first one
			if idx := strings.Index(xff, ","); idx != -1 {
				xff = xff[:idx]
			}
			xff = strings.TrimSpace(xff)
			if xff != "" {
				return xff
			}
		}

		// Check X-Real-IP header (only if from trusted proxy)
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			return xri
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If parsing fails, return the whole RemoteAddr
		return r.RemoteAddr
	}
	return ip
}
