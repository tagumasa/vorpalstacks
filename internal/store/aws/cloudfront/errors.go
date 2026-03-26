// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrNotFound is returned when the specified CloudFront entity does not exist.
	ErrNotFound = errors.New("entity not found")

	// ErrDistributionNotFound is returned when the specified distribution does not exist.
	ErrDistributionNotFound = errors.New("distribution not found")

	// ErrAlreadyExists is returned when attempting to create a CloudFront entity
	// that already exists.
	ErrAlreadyExists = errors.New("entity already exists")

	// ErrInUse is returned when the entity is in use and cannot be modified
	// or deleted.
	ErrInUse = errors.New("entity is in use")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = errors.New("internal error")
)

// StoreError represents an error that occurs during CloudFront store operations.
type StoreError struct {
	Operation string
	Err       error
}

// Error returns a string representation of the StoreError.
func (e *StoreError) Error() string {
	return e.Operation + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *StoreError) Unwrap() error {
	return e.Err
}

// NewStoreError creates a new CloudFront store error with the given operation and error.
func NewStoreError(operation string, err error) *StoreError {
	return &StoreError{Operation: operation, Err: err}
}

// IsNotFound checks if the error indicates that a CloudFront entity was not found.
func IsNotFound(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrNotFound) || common.IsNotFound(se.Err)
	}
	return common.IsNotFound(err)
}

// IsAlreadyExists checks if the error indicates that a CloudFront entity
// already exists.
func IsAlreadyExists(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrAlreadyExists) || common.IsAlreadyExists(se.Err)
	}
	return common.IsAlreadyExists(err)
}

// IsInUse checks if the error indicates that a CloudFront entity is in use.
func IsInUse(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrInUse)
	}
	return errors.Is(err, ErrInUse)
}
