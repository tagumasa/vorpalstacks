// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

import "time"

// NowUTC returns the current time in UTC.
//
// Returns:
//   - time.Time: Current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}
