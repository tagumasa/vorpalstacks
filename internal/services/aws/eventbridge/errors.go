package eventbridge

import (
	"fmt"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// EventsError represents an error from EventBridge.
type EventsError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *EventsError) Unwrap() error {
	return e.AWSError
}

// Error variables provide common EventBridge error instances.
var (
	ErrValidation            = &EventsError{awserrors.NewValidationException("Validation error")}
	ErrResourceNotFound      = &EventsError{awserrors.NewResourceNotFoundException("Resource", "")}
	ErrResourceAlreadyExists = &EventsError{awserrors.NewResourceAlreadyExistsException("Resource")}
	ErrInvalidParameter      = &EventsError{awserrors.NewInvalidParameterException("Invalid parameter")}
)

// NewValidationException creates a validation error.
func NewValidationException(message string) *EventsError {
	return &EventsError{awserrors.NewValidationException(message)}
}

// NewResourceNotFoundException creates a resource not found error.
func NewResourceNotFoundException(message string) *EventsError {
	return &EventsError{awserrors.NewAWSError("ResourceNotFoundException", message, 404)}
}

// NewResourceAlreadyExistsException creates a resource already exists error.
func NewResourceAlreadyExistsException(message string) *EventsError {
	return &EventsError{awserrors.NewResourceAlreadyExistsException(message)}
}

// NewInvalidParameterException creates an invalid parameter error.
func NewInvalidParameterException(message string) *EventsError {
	return &EventsError{awserrors.NewInvalidParameterException(message)}
}

// BuildEventBusARN constructs an ARN for an EventBridge event bus.
func BuildEventBusARN(accountID, region, name string) string {
	return fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", region, accountID, name)
}
