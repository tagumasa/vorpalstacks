// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	"errors"

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
type StoreError = common.StoreError

// NewStoreError creates a new store error for ACM operations.
func NewStoreError(op string, err error) *StoreError {
	return common.NewStoreError("acm", op, err)
}

// IsNotFound checks whether the error indicates a certificate was not found.
func IsNotFound(err error) bool {
	return common.IsNotFound(err) ||
		errors.Is(err, ErrCertificateNotFound)
}

// IsAlreadyExists checks whether the error indicates a certificate already exists.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrCertificateExists)
}

// IsInUse checks whether the error indicates a certificate is in use.
func IsInUse(err error) bool {
	return errors.Is(err, ErrCertificateInUse)
}
