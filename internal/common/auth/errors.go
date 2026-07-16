// Package auth provides AWS authentication utilities for vorpalstacks.
package auth

import "errors"

var (
	// ErrMissingAuthHeader is returned when the authorization header is missing
	// from the request.
	ErrMissingAuthHeader = errors.New("missing authorisation header")

	// ErrInvalidAlgorithm is returned when the authorization algorithm is not valid.
	ErrInvalidAlgorithm = errors.New("invalid authorization algorithm")

	// ErrMissingAuthFields is returned when required fields are missing from
	// the authorization header.
	ErrMissingAuthFields = errors.New("missing required fields in authorisation header")

	// ErrMissingDateHeader is returned when the X-Amz-Date header is missing
	// from the request.
	ErrMissingDateHeader = errors.New("missing X-Amz-Date header")

	// ErrReadBodyFailed is returned when the request body cannot be read.
	ErrReadBodyFailed = errors.New("failed to read request body")

	// ErrGetCredentialsFailed is returned when credentials cannot be obtained.
	ErrGetCredentialsFailed = errors.New("failed to get credentials")

	// ErrCalculateSignatureFailed is returned when the signature calculation fails.
	ErrCalculateSignatureFailed = errors.New("failed to calculate signature")

	// ErrSignatureMismatch is returned when the calculated signature does not
	// match the provided signature.
	ErrSignatureMismatch = errors.New("signature mismatch")

	// ErrInvalidCredentialFormat is returned when the credential format is not valid.
	ErrInvalidCredentialFormat = errors.New("invalid credential format")

	// ErrParseDateFailed is returned when the date cannot be parsed.
	ErrParseDateFailed = errors.New("failed to parse date")
)
