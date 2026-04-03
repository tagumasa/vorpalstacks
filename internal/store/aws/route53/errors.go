package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import (
	"errors"

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
type StoreError = common.StoreError

// NewStoreError creates a new Route 53 store error with the given operation and error.
func NewStoreError(op string, err error) *StoreError {
	return common.NewStoreError("route53", op, err)
}

// IsNotFound checks if the error indicates that a Route 53 resource
// was not found. This includes hosted zones, health checks, record sets, and changes.
func IsNotFound(err error) bool {
	return common.IsNotFound(err) ||
		errors.Is(err, ErrHostedZoneNotFound) ||
		errors.Is(err, ErrHealthCheckNotFound) ||
		errors.Is(err, ErrRecordSetNotFound) ||
		errors.Is(err, ErrChangeNotFound)
}

// IsAlreadyExists checks if the error indicates that a Route 53 resource
// already exists. This includes hosted zones, health checks, and record sets.
func IsAlreadyExists(err error) bool {
	return common.IsAlreadyExists(err) ||
		errors.Is(err, ErrHostedZoneExists) ||
		errors.Is(err, ErrHealthCheckExists) ||
		errors.Is(err, ErrRecordSetExists)
}
