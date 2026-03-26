// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import "net"

// IsLocalhost checks if an IP address is a localhost/loopback address.
//
// Parameters:
//   - ip: The IP address to check
//
// Returns:
//   - bool: True if the IP is a loopback address (127.0.0.0/8 or ::1)
//
// Example:
//
//	IsLocalhost("127.0.0.1")    // true
//	IsLocalhost("::1")          // true
//	IsLocalhost("192.168.1.1") // false
func IsLocalhost(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.IsLoopback()
}
