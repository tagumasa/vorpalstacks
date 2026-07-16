// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

import "time"

// ISO8601Format is standard ISO 8601 format string.
const ISO8601Format = "2006-01-02T15:04:05.999Z07:00"

// ISO8601UTCFormat is ISO 8601 format in UTC with milliseconds.
const ISO8601UTCFormat = "2006-01-02T15:04:05.000Z"

// ISO8601SimpleFormat is ISO 8601 format without milliseconds (AWS standard).
const ISO8601SimpleFormat = "2006-01-02T15:04:05Z"

// ISO8601NoZFormat is ISO 8601 format without Z suffix (used for parsing).
const ISO8601NoZFormat = "2006-01-02T15:04:05"

// ISO8601WithTimezoneFormat is ISO 8601 format with timezone offset.
const ISO8601WithTimezoneFormat = "2006-01-02T15:04:05-07:00"

// RFC3339Format is RFC3339 format string.
const RFC3339Format = time.RFC3339

// FormatTime formats a time.Time to specified format.
// If format is empty, it defaults to ISO8601UTCFormat.
//
// Parameters:
//   - t: The time to format
//   - format: The format string (optional)
//
// Returns:
//   - string: The formatted time string
//
// Example:
//
//	t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
//	FormatTime(t, ISO8601UTCFormat) // "2024-01-01T12:00:00.000Z"
func FormatTime(t time.Time, format string) string {
	if format == "" {
		format = ISO8601UTCFormat
	}
	return t.Format(format)
}
