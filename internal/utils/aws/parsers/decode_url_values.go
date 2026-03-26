// Package parsers provides AWS-specific utility functions for vorpalstacks.
package parsers

import (
	"net/url"
)

// DecodeURLValues decodes URL-encoded values in condition context.
//
// Parameters:
//   - values: The map of URL-encoded values to decode
//
// Returns:
//   - map[string]string: The decoded values
func DecodeURLValues(values map[string]string) map[string]string {
	if values == nil {
		return make(map[string]string)
	}
	result := make(map[string]string, len(values))
	for k, v := range values {
		if decoded, err := url.QueryUnescape(v); err == nil {
			result[k] = decoded
		} else {
			result[k] = v
		}
	}
	return result
}
