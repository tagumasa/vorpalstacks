// Package logs provides logging functionality for vorpalstacks.
package logs

import "time"

// Default constants for log configuration.
const (
	DefaultMaxLogSize       = 100 * 1024 * 1024
	DefaultMaxLogEntries    = 1000000
	DefaultLogMaxAge        = 7 * 24 * time.Hour
	DefaultRotationInterval = time.Hour
)
