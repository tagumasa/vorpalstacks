// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

import "time"

// ParseTime parses a time string using specified format.
// If format is empty, it tries multiple common formats.
//
// Parameters:
//   - s: The time string to parse
//   - format: The format string (optional)
//
// Returns:
//   - time.Time: The parsed time
//   - error: An error if parsing fails
//
// Example:
//
//	ParseTime("2024-01-01T12:00:00.000Z", ISO8601UTCFormat)
func ParseTime(s string, format string) (time.Time, error) {
	if format == "" {
		return ParseAutoTime(s)
	}
	return time.Parse(format, s)
}
