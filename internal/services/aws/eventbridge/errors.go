package eventbridge

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
)

// Error variables provide common EventBridge error instances.
var (
	// ErrValidation is returned when validation fails.
	ErrValidation = awserrors.NewValidationException("Validation error")
	// ErrResourceNotFound is returned when a resource is not found.
	ErrResourceNotFound = awserrors.NewResourceNotFoundException("Resource", "")
	// ErrResourceAlreadyExists is returned when a resource already exists.
	ErrResourceAlreadyExists = awserrors.NewResourceAlreadyExistsException("Resource")
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = awserrors.NewInvalidParameterException("Invalid parameter")
)

// NewValidationException creates a validation error.
func NewValidationException(message string) *awserrors.AWSError {
	return awserrors.NewValidationException(message)
}

// NewResourceNotFoundException creates a resource not found error.
func NewResourceNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceNotFoundException", message, 404)
}

// NewResourceAlreadyExistsException creates a resource already exists error.
func NewResourceAlreadyExistsException(message string) *awserrors.AWSError {
	return awserrors.NewResourceAlreadyExistsException(message)
}

// NewInvalidParameterException creates an invalid parameter error.
func NewInvalidParameterException(message string) *awserrors.AWSError {
	return awserrors.NewInvalidParameterException(message)
}

// BuildEventBusARN constructs an ARN for an EventBridge event bus.
func BuildEventBusARN(accountID, region, name string) string {
	return fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", region, accountID, name)
}
