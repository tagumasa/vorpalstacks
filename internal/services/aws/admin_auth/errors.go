// Package admin_auth provides administrative authentication operations.
package admin_auth

import "errors"

// Error variables for admin authentication.
var (
	// ErrMissingCredentials is returned when username or password is not provided.
	ErrMissingCredentials = errors.New("username and password are required")

	// ErrInvalidCredentials is returned when the username or password is incorrect.
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrMissingRefreshToken is returned when the refresh token is not provided.
	ErrMissingRefreshToken = errors.New("refresh token is required")

	// ErrRefreshNotImplemented is returned when the refresh token flow is not
	// implemented.
	ErrRefreshNotImplemented = errors.New("refresh token flow not implemented")

	// ErrMissingAccessToken is returned when the access token is not provided.
	ErrMissingAccessToken = errors.New("access token is required")

	// ErrInvalidToken is returned when the token is invalid or has expired.
	ErrInvalidToken = errors.New("invalid or expired token")

	// ErrUserNotFound is returned when the specified user does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidRefreshToken is returned when the refresh token is invalid
	// or has expired.
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)
