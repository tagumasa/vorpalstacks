package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrNotFound is returned when the specified WAF entity does not exist.
	ErrNotFound = errors.New("entity not found")

	// ErrAlreadyExists is returned when attempting to create a WAF entity
	// that already exists.
	ErrAlreadyExists = errors.New("entity already exists")

	// ErrInUse is returned when the entity is in use and cannot be deleted.
	ErrInUse = errors.New("entity is in use")

	// ErrLockTokenMismatch is returned when the lock token does not match.
	ErrLockTokenMismatch = errors.New("lock token mismatch")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = errors.New("internal error")
)

// StoreError represents an error that occurs during WAF store operations.
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

// NewStoreError creates a new WAF store error with the given operation and error.
func NewStoreError(operation string, err error) *StoreError {
	return &StoreError{Operation: operation, Err: err}
}

// IsNotFound checks if the error indicates that a WAF entity was not found.
func IsNotFound(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrNotFound) || errors.Is(se.Err, common.ErrNotFound)
	}
	return errors.Is(err, ErrNotFound) || errors.Is(err, common.ErrNotFound)
}

// IsAlreadyExists checks if the error indicates that a WAF entity already exists.
func IsAlreadyExists(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrAlreadyExists) || errors.Is(se.Err, common.ErrAlreadyExists)
	}
	return errors.Is(err, ErrAlreadyExists) || errors.Is(err, common.ErrAlreadyExists)
}

// IsInUse checks if the error indicates that a WAF entity is in use.
func IsInUse(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrInUse)
	}
	return errors.Is(err, ErrInUse)
}

// IsLockTokenMismatch checks if the error indicates a lock token mismatch.
func IsLockTokenMismatch(err error) bool {
	var se *StoreError
	if errors.As(err, &se) {
		return errors.Is(se.Err, ErrLockTokenMismatch)
	}
	return errors.Is(err, ErrLockTokenMismatch)
}
