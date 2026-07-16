package admin_auth

import "errors"

var (
	// ErrTokenNotFound is returned when a refresh token is not found.
	ErrTokenNotFound = errors.New("refresh token not found")
	// ErrTokenExpired is returned when a refresh token has expired.
	ErrTokenExpired = errors.New("refresh token has expired")
)
