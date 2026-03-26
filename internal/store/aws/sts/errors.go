package sts

// Package sts provides STS (Security Token Service) data store implementations
// for vorpalstacks.

import (
	"errors"
)

var (
	// ErrSessionNotFound is returned when the specified STS session
	// does not exist.
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired is returned when the STS session has expired.
	ErrSessionExpired = errors.New("session expired")
)
