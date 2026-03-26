// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

import (
	"fmt"
	"time"
)

// ParseAutoTime attempts to parse a time string using multiple common formats.
//
// Parameters:
//   - s: The time string to parse
//
// Returns:
//   - time.Time: The parsed time
//   - error: An error if parsing fails
//
// Example:
//
//	ParseAutoTime("2024-01-01T12:00:00Z")
//	ParseAutoTime("2024-01-01 12:00:00")
func ParseAutoTime(s string) (time.Time, error) {
	formats := []string{
		ISO8601Format,
		ISO8601UTCFormat,
		RFC3339Format,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		ISO8601SimpleFormat,
		"2006-01-02",
	}

	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}
