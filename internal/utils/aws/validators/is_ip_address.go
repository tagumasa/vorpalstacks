// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"net"
)

// IsIPAddress checks if a string is formatted as an IP address.
//
// Parameters:
//   - s: The string to check
//
// Returns:
//   - bool: True if the string is a valid IP address
func IsIPAddress(s string) bool {
	// Check if string is a valid IPv4 or IPv6 address
	return net.ParseIP(s) != nil
}
