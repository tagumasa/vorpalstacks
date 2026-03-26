// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

// IsPreflightRequest checks if a request is a CORS preflight request.
//
// Parameters:
//   - method: The HTTP method to check
//
// Returns:
//   - bool: True if the method is OPTIONS
func IsPreflightRequest(method string) bool {
	return method == "OPTIONS"
}
