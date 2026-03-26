package scheduler

import (
	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// SchedulerError represents an error from the EventBridge Scheduler service.
type SchedulerError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *SchedulerError) Unwrap() error {
	return e.AWSError
}

// Error variables provide common Scheduler error instances.
var (
	ErrScheduleGroupNotFound        = &SchedulerError{awserrors.NewResourceNotFoundException("Schedule group", "")}
	ErrScheduleGroupAlreadyExists   = &SchedulerError{awserrors.NewConflictException("Schedule group already exists")}
	ErrScheduleNotFound             = &SchedulerError{awserrors.NewResourceNotFoundException("Schedule", "")}
	ErrScheduleAlreadyExists        = &SchedulerError{awserrors.NewConflictException("Schedule already exists")}
	ErrValidation                   = &SchedulerError{awserrors.NewValidationException("Validation error")}
	ErrInvalidScheduleExpression    = &SchedulerError{awserrors.NewValidationException("Invalid schedule expression")}
	ErrInvalidTarget                = &SchedulerError{awserrors.NewValidationException("Invalid target")}
	ErrInvalidFlexibleTimeWindow    = &SchedulerError{awserrors.NewValidationException("Invalid flexible time window")}
	ErrInvalidDate                  = &SchedulerError{awserrors.NewValidationException("Invalid date format")}
	ErrInvalidState                 = &SchedulerError{awserrors.NewValidationException("Invalid state value")}
	ErrInvalidActionAfterCompletion = &SchedulerError{awserrors.NewValidationException("Invalid action after completion value")}
	ErrScheduleGroupNotEmpty        = &SchedulerError{awserrors.NewConflictException("Schedule group is not empty")}
	ErrInternalServer               = &SchedulerError{awserrors.NewInternalErrorException("Internal server error")}
)

// NewValidationException creates a validation error.
func NewValidationException(message string) *SchedulerError {
	return &SchedulerError{awserrors.NewValidationException(message)}
}

// NewResourceNotFoundException creates a resource not found error.
func NewResourceNotFoundException(resource string) *SchedulerError {
	return &SchedulerError{awserrors.NewNotFoundException(resource)}
}

// NewConflictException creates a conflict error.
func NewConflictException(message string) *SchedulerError {
	return &SchedulerError{awserrors.NewConflictException(message)}
}
