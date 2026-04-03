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
	ErrInvalidParameter = common.ErrInvalidParameter

	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = errors.New("internal error")
)

// StoreError represents an error that occurs during CloudFront store operations.
type StoreError = common.StoreError

// NewStoreError creates a new CloudFront store error with the given operation and error.
func NewStoreError(operation string, err error) *StoreError {
	return common.NewStoreError("cloudfront", operation, err)
}

// IsNotFound checks if the error indicates that a CloudFront entity was not found.
func IsNotFound(err error) bool {
	return common.IsNotFound(err) ||
		errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrDistributionNotFound)
}

// IsAlreadyExists checks if the error indicates that a CloudFront entity
// already exists.
func IsAlreadyExists(err error) bool {
	return common.IsAlreadyExists(err) ||
		errors.Is(err, ErrAlreadyExists)
}

// IsInUse checks if the error indicates that a CloudFront entity is in use.
func IsInUse(err error) bool {
	return errors.Is(err, ErrInUse)
}
