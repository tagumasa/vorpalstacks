// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import (
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	trustedProxies   = []string{}
	trustedProxiesMu sync.RWMutex
)

// SetTrustedProxies sets the list of trusted proxy IP addresses.
// This is used by GetHost, GetScheme, and GetClientIP to determine
// whether to trust proxy headers (X-Forwarded-Host, X-Forwarded-Proto, X-Forwarded-For).
//
// Parameters:
//   - proxies: A slice of IP addresses to trust. Use "*" to trust all proxies.
//
// Example:
//
//	SetTrustedProxies([]string{"10.0.0.1", "10.0.0.2"})
//	SetTrustedProxies([]string{"*"})  // Trust all proxies (not recommended for production)
func SetTrustedProxies(proxies []string) {
	trustedProxiesMu.Lock()
	defer trustedProxiesMu.Unlock()
	trustedProxies = proxies
}

// GetTrustedProxies returns the list of trusted proxy IP addresses.
// This is used by GetHost, GetScheme, and GetClientIP to determine
// whether to trust proxy headers.
//
// Returns:
//   - []string: A slice of trusted proxy IP addresses.
//
// Example:
//
//	proxies := GetTrustedProxies()
func GetTrustedProxies() []string {
	trustedProxiesMu.RLock()
	defer trustedProxiesMu.RUnlock()
	return trustedProxies
}

// GetHost extracts the host from an HTTP request.
// It checks the X-Forwarded-Host header if the request comes from a trusted proxy.
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - string: The host value (from X-Forwarded-Host if trusted, otherwise from r.Host)
//
// Example:
//
//	host := GetHost(r)
func GetHost(r *http.Request) string {
	forwardedHost := r.Header.Get("X-Forwarded-Host")
	if forwardedHost != "" && isTrustedProxy(r) {
		return forwardedHost
	}
	return r.Host
}

func isTrustedProxy(r *http.Request) bool {
	trustedProxiesMu.RLock()
	proxies := trustedProxies
	trustedProxiesMu.RUnlock()

	if len(proxies) == 0 {
		return false
	}
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	for _, proxy := range proxies {
		if proxy == remoteIP || proxy == "*" {
			return true
		}
	}
	clientIPHeader := r.Header.Get("X-Forwarded-For")
	if clientIPHeader != "" {
		parts := strings.Split(clientIPHeader, ",")
		if len(parts) > 0 {
			clientIP := strings.TrimSpace(parts[0])
			for _, proxy := range proxies {
				if proxy == clientIP || proxy == "*" {
					return true
				}
			}
		}
	}
	return false
}
