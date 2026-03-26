package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import (
	"errors"
	"fmt"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrHostedZoneNotFound is returned when the specified Route 53 hosted zone
	// does not exist.
	ErrHostedZoneNotFound = errors.New("hosted zone not found")

	// ErrHostedZoneExists is returned when attempting to create a hosted zone
	// that already exists.
	ErrHostedZoneExists = errors.New("hosted zone already exists")

	// ErrHealthCheckNotFound is returned when the specified health check
	// does not exist.
	ErrHealthCheckNotFound = errors.New("health check not found")

	// ErrHealthCheckExists is returned when attempting to create a health check
	// that already exists.
	ErrHealthCheckExists = errors.New("health check already exists")

	// ErrRecordSetNotFound is returned when the specified resource record set
	// does not exist.
	ErrRecordSetNotFound = errors.New("resource record set not found")

	// ErrRecordSetExists is returned when attempting to create a resource record set
	// that already exists.
	ErrRecordSetExists = errors.New("resource record set already exists")

	// ErrChangeNotFound is returned when the specified change batch
	// does not exist.
	ErrChangeNotFound = errors.New("change not found")

	// ErrInvalidInput is returned when the input is not valid.
	ErrInvalidInput = errors.New("invalid input")
)

// StoreError represents an error that occurs during Route 53 store operations.
// It encapsulates the service name, operation, and underlying error.
type StoreError struct {
	Service string
	Op      string
	Err     error
}

// Error returns a string representation of the StoreError.
func (e *StoreError) Error() string {
	return fmt.Sprintf("route53 store %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error.
func (e *StoreError) Unwrap() error {
	return e.Err
}

// NewStoreError creates a new Route 53 store error with the given operation and error.
func NewStoreError(op string, err error) *StoreError {
	return &StoreError{Service: "route53", Op: op, Err: err}
}

// IsNotFound checks if the error indicates that a Route 53 resource
// was not found. This includes hosted zones, health checks, record sets, and changes.
func IsNotFound(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrHostedZoneNotFound) ||
			errors.Is(storeErr.Err, ErrHealthCheckNotFound) ||
			errors.Is(storeErr.Err, ErrRecordSetNotFound) ||
			errors.Is(storeErr.Err, ErrChangeNotFound) ||
			errors.Is(storeErr.Err, common.ErrNotFound)
	}
	return errors.Is(err, ErrHostedZoneNotFound) ||
		errors.Is(err, ErrHealthCheckNotFound) ||
		errors.Is(err, ErrRecordSetNotFound) ||
		errors.Is(err, ErrChangeNotFound) ||
		errors.Is(err, common.ErrNotFound)
}

// IsAlreadyExists checks if the error indicates that a Route 53 resource
// already exists. This includes hosted zones, health checks, and record sets.
func IsAlreadyExists(err error) bool {
	var storeErr *StoreError
	if errors.As(err, &storeErr) {
		return errors.Is(storeErr.Err, ErrHostedZoneExists) ||
			errors.Is(storeErr.Err, ErrHealthCheckExists) ||
			errors.Is(storeErr.Err, ErrRecordSetExists) ||
			errors.Is(storeErr.Err, common.ErrAlreadyExists)
	}
	return errors.Is(err, ErrHostedZoneExists) ||
		errors.Is(err, ErrHealthCheckExists) ||
		errors.Is(err, ErrRecordSetExists) ||
		errors.Is(err, common.ErrAlreadyExists)
}
