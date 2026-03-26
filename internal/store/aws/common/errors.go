// Package common provides common AWS store utilities for vorpalstacks.
package common

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when a resource already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrInvalidState is returned when a resource is in an invalid state.
	ErrInvalidState = errors.New("invalid state")
	// ErrInvalidInput is returned when input is invalid.
	ErrInvalidInput = errors.New("invalid input")
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = errors.New("access denied")
	// ErrConflict is returned when a conflict occurs.
	ErrConflict = errors.New("conflict")
	// ErrInternal is returned when an internal error occurs.
	ErrInternal = errors.New("internal error")
)

// StoreError represents an error in a store operation.
type StoreError struct {
	Service string
	Op      string
	Key     string
	Err     error
}

// Error returns the error message.
func (e *StoreError) Error() string {
	if e.Key != "" {
		return fmt.Sprintf("%s store %s %s: %v", e.Service, e.Op, e.Key, e.Err)
	}
	return fmt.Sprintf("%s store %s: %v", e.Service, e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e *StoreError) Unwrap() error {
	return e.Err
}

// NewStoreError creates a new store error.
func NewStoreError(service, op string, err error) *StoreError {
	return &StoreError{Service: service, Op: op, Err: err}
}

// NewStoreErrorWithKey creates a new store error with a key.
func NewStoreErrorWithKey(service, op, key string, err error) *StoreError {
	return &StoreError{Service: service, Op: op, Key: key, Err: err}
}

// IsNotFound checks whether the error indicates a resource was not found.
func IsNotFound(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrNotFound)
	}
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists checks whether the error indicates a resource already exists.
func IsAlreadyExists(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrAlreadyExists)
	}
	return errors.Is(err, ErrAlreadyExists)
}
