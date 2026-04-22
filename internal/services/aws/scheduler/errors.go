package scheduler

import (
	awserrors "vorpalstacks/internal/common/errors"
)

// Error variables provide common Scheduler error instances.
var (
	// ErrScheduleGroupNotFound is returned when the specified schedule group does not exist.
	ErrScheduleGroupNotFound = awserrors.NewResourceNotFoundException("Schedule group", "")
	// ErrScheduleGroupAlreadyExists is returned when a schedule group already exists.
	ErrScheduleGroupAlreadyExists = awserrors.NewConflictException("Schedule group already exists")
	// ErrScheduleNotFound is returned when the specified schedule does not exist.
	ErrScheduleNotFound = awserrors.NewResourceNotFoundException("Schedule", "")
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
	// ErrScheduleGroupNotEmpty is returned when the schedule group is not empty.
	ErrScheduleGroupNotEmpty = awserrors.NewConflictException("Schedule group is not empty")
	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = awserrors.NewInternalErrorException("Internal server error")
)
