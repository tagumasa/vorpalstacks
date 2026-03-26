// Package timeutils provides time utility functions for vorpalstacks.
package timeutils

import (
	"fmt"
	"strings"
	"time"
)

// DurationString returns a human-readable string representation of a duration.
//
// Parameters:
//   - d: The duration to format
//
// Returns:
//   - string: Human-readable duration string
//
// Example:
//
//	DurationString(90 * time.Second) // "1m30s"
//	DurationString(2*time.Hour + 5*time.Minute) // "2h5m"
func DurationString(d time.Duration) string {
	if d < time.Second {
		return d.String()
	}

	var parts []string

	hours := int(d.Hours())
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}

	minutes := int(d.Minutes()) % 60
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	seconds := int(d.Seconds()) % 60
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, "")
}
