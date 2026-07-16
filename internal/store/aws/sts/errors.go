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

	// ErrDelegationTokenNotFound is returned when the specified trade-in
	// token does not exist.
	ErrDelegationTokenNotFound = errors.New("delegated token not found")

	// ErrDelegationTokenExpired is returned when the trade-in token has expired.
	ErrDelegationTokenExpired = errors.New("delegated token expired")
)
