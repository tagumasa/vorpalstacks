package scheduler

import "errors"

var (
	// ErrScheduleGroupNotFound is returned when the specified EventBridge Scheduler
	// schedule group does not exist.
	ErrScheduleGroupNotFound = errors.New("schedule group not found")

	// ErrScheduleGroupAlreadyExists is returned when attempting to create a schedule
	// group that already exists.
	ErrScheduleGroupAlreadyExists = errors.New("schedule group already exists")

	// ErrScheduleGroupDeleting is returned when an operation would add or
	// modify schedule records in a group whose deletion is in progress —
	// the engine cascade would destroy the acknowledged write.
	ErrScheduleGroupDeleting = errors.New("schedule group is being deleted")

	// ErrScheduleNotFound is returned when the specified schedule does not exist.
	ErrScheduleNotFound = errors.New("schedule not found")

	// ErrScheduleAlreadyExists is returned when attempting to create a schedule
	// that already exists.
	ErrScheduleAlreadyExists = errors.New("schedule already exists")

	// ErrInvalidName is returned when a schedule or group record arrives
	// with an empty name — records are keyed by name, so the write has no
	// identity to persist under.
	ErrInvalidName = errors.New("invalid name")

	// ErrScheduleGroupNotEmpty is returned when attempting to purge a
	// schedule group that still contains member schedules.
	ErrScheduleGroupNotEmpty = errors.New("schedule group not empty")
)
