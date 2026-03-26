// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

import "time"

// TimeBetween checks if a time is between two other times (inclusive).
//
// Parameters:
//   - t: The time to check
//   - start: The start time
//   - end: The end time
//
// Returns:
//   - bool: True if t is between start and end (inclusive)
func TimeBetween(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}
