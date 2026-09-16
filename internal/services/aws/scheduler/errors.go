package scheduler

import (
	awserrors "vorpalstacks/internal/common/errors"
)

// Error variables provide common Scheduler error instances.
var (
	// ErrScheduleGroupAlreadyExists is returned when a schedule group already exists.
	ErrScheduleGroupAlreadyExists = awserrors.NewConflictException("Schedule group already exists")
	// ErrScheduleGroupDeleting is returned when the target schedule group's
	// deletion is in progress — the group exists but admits no new or
	// updated schedules (the cascade would destroy the acknowledged write).
	ErrScheduleGroupDeleting = awserrors.NewConflictException("Schedule group is being deleted")
	// ErrScheduleAlreadyExists is returned when a schedule already exists.
	ErrScheduleAlreadyExists = awserrors.NewConflictException("Schedule already exists")
	// ErrValidation is returned when validation fails.
	ErrValidation = awserrors.NewValidationException("Validation error")
	// ErrInvalidScheduleExpression is returned when the schedule expression is invalid.
	ErrInvalidScheduleExpression = awserrors.NewValidationException("Invalid schedule expression")
	// ErrInvalidTarget is returned when the target is invalid.
	ErrInvalidTarget = awserrors.NewValidationException("Invalid target")
	// ErrInvalidFlexibleTimeWindow is returned when the flexible time window is invalid.
	ErrInvalidFlexibleTimeWindow = awserrors.NewValidationException("Invalid flexible time window")
	// ErrInvalidDate is returned when the date format is invalid.
	ErrInvalidDate = awserrors.NewValidationException("Invalid date format")
	// ErrInvalidState is returned when the state value is invalid.
	ErrInvalidState = awserrors.NewValidationException("Invalid state value")
	// ErrInvalidActionAfterCompletion is returned when the action after completion value is invalid.
	ErrInvalidActionAfterCompletion = awserrors.NewValidationException("Invalid action after completion value")
	// ErrInternalServer is returned when an internal server error occurs.
	// InternalServerException is the scheduler model's server-fault shape
	// (smithy.api#error "server", httpError 500).
	ErrInternalServer = awserrors.NewInternalServerException("Unexpected error encountered while processing the request.")
)

// scheduleNotFound builds the modelled ResourceNotFoundException naming the
// schedule the request addressed — every raise site knows the identifier,
// and a shared empty-identifier sentinel would format the message as
// "Schedule  not found".
func scheduleNotFound(name string) error {
	return awserrors.NewResourceNotFoundException("Schedule", name)
}

// scheduleGroupNotFound builds the modelled ResourceNotFoundException naming
// the schedule group the request addressed.
func scheduleGroupNotFound(name string) error {
	return awserrors.NewResourceNotFoundException("Schedule group", name)
}
