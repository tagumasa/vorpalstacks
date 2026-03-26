// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	"errors"
	"fmt"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrCertificateNotFound is returned when a certificate is not found.
	ErrCertificateNotFound = errors.New("certificate not found")
	// ErrCertificateExists is returned when a certificate already exists.
	ErrCertificateExists = errors.New("certificate already exists")
	// ErrCertificateInUse is returned when a certificate is in use.
	ErrCertificateInUse = errors.New("certificate is in use")
	// ErrInvalidArn is returned when an ARN is invalid.
	ErrInvalidArn = errors.New("invalid ARN")
	// ErrInvalidDomain is returned when a domain name is invalid.
	ErrInvalidDomain = errors.New("invalid domain name")
	// ErrInvalidState is returned when a certificate state is invalid.
	ErrInvalidState = errors.New("invalid state")
)

// StoreError represents an error in the ACM store.
type StoreError struct {
	Service string
	Op      string
	Err     error
}

// Error returns a string representation of the store error.
func (e *StoreError) Error() string {
	return fmt.Sprintf("acm store %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e *StoreError) Unwrap() error {
	return e.Err
}

// NewStoreError creates a new store error for ACM operations.
func NewStoreError(op string, err error) *StoreError {
	return &StoreError{Service: "acm", Op: op, Err: err}
}

// IsNotFound checks whether the error indicates a certificate was not found.
func IsNotFound(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrCertificateNotFound) || errors.Is(storeErr.Err, common.ErrNotFound)
	}
	return errors.Is(err, ErrCertificateNotFound) || errors.Is(err, common.ErrNotFound)
}

// IsAlreadyExists checks whether the error indicates a certificate already exists.
func IsAlreadyExists(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrCertificateExists)
	}
	return errors.Is(err, ErrCertificateExists)
}

// IsInUse checks whether the error indicates a certificate is in use.
func IsInUse(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrCertificateInUse)
	}
	return errors.Is(err, ErrCertificateInUse)
}
