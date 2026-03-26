// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

// NowString returns current time as a string in ISO8601UTCFormat.
//
// Returns:
//   - string: Current time as ISO8601 UTC string
func NowString() string {
	return NowUTC().Format(ISO8601UTCFormat)
}
