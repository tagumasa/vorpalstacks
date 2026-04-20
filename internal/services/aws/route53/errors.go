// Package route53 provides Route 53 service operations for vorpalstacks.
package route53

import (
	stderrors "errors"
	"net/http"

	"vorpalstacks/internal/common/errors"
	route53store "vorpalstacks/internal/store/aws/route53"
)

var (
	// ErrNoSuchHostedZone is returned when the specified hosted zone does not exist.
	ErrNoSuchHostedZone = errors.NewAWSError("NoSuchHostedZone", "The hosted zone does not exist.", http.StatusNotFound)
	// ErrHostedZoneAlreadyExists is returned when attempting to create a hosted zone that already exists.
	ErrHostedZoneAlreadyExists = errors.NewAWSError("HostedZoneAlreadyExists", "The hosted zone already exists.", http.StatusConflict)
	// ErrHostedZoneNotEmpty is returned when attempting to delete a non-empty hosted zone.
	ErrHostedZoneNotEmpty = errors.NewAWSError("HostedZoneNotEmpty", "The hosted zone must be empty before it can be deleted.", http.StatusBadRequest)
	// ErrNoSuchHealthCheck is returned when the specified health check does not exist.
	ErrNoSuchHealthCheck = errors.NewAWSError("NoSuchHealthCheck", "The health check does not exist.", http.StatusNotFound)
	// ErrHealthCheckAlreadyExists is returned when attempting to create a health check that already exists.
	ErrHealthCheckAlreadyExists = errors.NewAWSError("HealthCheckAlreadyExists", "The health check already exists.", http.StatusConflict)
	// ErrNoSuchChange is returned when the specified change does not exist.
	ErrNoSuchChange = errors.NewAWSError("NoSuchChange", "The change does not exist.", http.StatusNotFound)
	// ErrInvalidInput is returned when the input is invalid.
	ErrInvalidInput = errors.NewAWSError("InvalidInput", "The input is invalid.", http.StatusBadRequest)
)

// NewNoSuchHostedZoneError creates a new error for a hosted zone that does not exist.
func NewNoSuchHostedZoneError(id string) *errors.AWSError {
	return errors.NewAWSError("NoSuchHostedZone", "No hosted zone found with id: "+id, http.StatusNotFound)
}

// NewHostedZoneAlreadyExistsError creates a new error for a hosted zone that already exists.
func NewHostedZoneAlreadyExistsError(name string) *errors.AWSError {
	return errors.NewAWSError("HostedZoneAlreadyExists", "Hosted zone already exists: "+name, http.StatusConflict)
}

// NewHostedZoneNotEmptyError creates a new error for a non-empty hosted zone.
func NewHostedZoneNotEmptyError() *errors.AWSError {
	return errors.NewAWSError("HostedZoneNotEmpty", "The hosted zone must be empty before it can be deleted.", http.StatusBadRequest)
}

// NewNoSuchHealthCheckError creates a new error for a health check that does not exist.
func NewNoSuchHealthCheckError(id string) *errors.AWSError {
	return errors.NewAWSError("NoSuchHealthCheck", "No health check found with id: "+id, http.StatusNotFound)
}

// NewHealthCheckAlreadyExistsError creates a new error for a health check that already exists.
func NewHealthCheckAlreadyExistsError() *errors.AWSError {
	return errors.NewAWSError("HealthCheckAlreadyExists", "The health check already exists.", http.StatusConflict)
}

// NewNoSuchChangeError creates a new error for a change that does not exist.
func NewNoSuchChangeError(id string) *errors.AWSError {
	return errors.NewAWSError("NoSuchChange", "No change found with id: "+id, http.StatusNotFound)
}

// mapStoreError converts a route53 store error into an appropriate AWS API error.
// Infrastructure errors (s.store() failures) should NOT be passed through this function.
func mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	if stderrors.Is(err, route53store.ErrHostedZoneExists) {
		return NewHostedZoneAlreadyExistsError("")
	}
	if stderrors.Is(err, route53store.ErrHealthCheckExists) {
		return NewHealthCheckAlreadyExistsError()
	}
	if stderrors.Is(err, route53store.ErrRecordSetNotFound) {
		return errors.NewAWSError("InvalidChangeBatch", "The resource record set does not exist.", http.StatusBadRequest)
	}
	if stderrors.Is(err, route53store.ErrRecordSetExists) {
		return errors.NewAWSError("InvalidChangeBatch", "The resource record set already exists.", http.StatusBadRequest)
	}
	return err
}
