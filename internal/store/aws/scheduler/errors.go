package scheduler

import "errors"

var (
	// ErrScheduleGroupNotFound is returned when the specified EventBridge Scheduler
	// schedule group does not exist.
	ErrScheduleGroupNotFound = errors.New("schedule group not found")

	// ErrScheduleGroupAlreadyExists is returned when attempting to create a schedule
	// group that already exists.
	ErrScheduleGroupAlreadyExists = errors.New("schedule group already exists")

	// ErrScheduleNotFound is returned when the specified schedule does not exist.
	ErrScheduleNotFound = errors.New("schedule not found")

	// ErrScheduleAlreadyExists is returned when attempting to create a schedule
	// that already exists.
	ErrScheduleAlreadyExists = errors.New("schedule already exists")

	// ErrInvalidARN is returned when the Amazon Resource Name (ARN) is not valid.
	ErrInvalidARN = errors.New("invalid ARN")

	// ErrInvalidScheduleExpression is returned when the schedule expression
	// is not valid.
	ErrScheduleGroupNotEmpty = errors.New("schedule group not empty")
)
