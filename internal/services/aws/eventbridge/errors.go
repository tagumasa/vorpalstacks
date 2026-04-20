package eventbridge

import (
	arnutil "vorpalstacks/internal/utils/aws/arn"

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

// NewResourceNotFoundException creates a resource not found error.
func NewResourceNotFoundException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceNotFoundException", message, 404)
}

// BuildEventBusARN constructs an ARN for an EventBridge event bus.
func BuildEventBusARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Events().EventBus(name)
}
