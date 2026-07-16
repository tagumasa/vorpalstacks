// Package timeutils provides time-related utility functions for vorpalstacks.
package timeutils

import "time"

// FormatEpochSeconds converts a time.Time to Unix epoch seconds.
//
// Parameters:
//   - t: The time to convert
//
// Returns:
//   - float64: The Unix epoch time in seconds (with milliseconds as decimal)
func FormatEpochSeconds(t time.Time) float64 {
	if t.IsZero() {
		return 0
	}
	return float64(t.UnixMilli()) / 1000.0
}
