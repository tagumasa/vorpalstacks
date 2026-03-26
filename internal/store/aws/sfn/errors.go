package sfn

import "errors"

var (
	// ErrStateMachineNotFound is returned when the specified Step Functions state
	// machine does not exist.
	ErrStateMachineNotFound = errors.New("state machine not found")

	// ErrStateMachineAlreadyExists is returned when attempting to create a state
	// machine that already exists.
	ErrStateMachineAlreadyExists = errors.New("state machine already exists")

	// ErrExecutionNotFound is returned when the specified execution does not exist.
	ErrExecutionNotFound = errors.New("execution not found")

	// ErrExecutionAlreadyExists is returned when attempting to create an execution
	// that already exists.
	ErrExecutionAlreadyExists = errors.New("execution already exists")

	// ErrActivityNotFound is returned when the specified activity does not exist.
	ErrActivityNotFound = errors.New("activity not found")

	// ErrActivityAlreadyExists is returned when attempting to create an activity
	// that already exists.
	ErrActivityAlreadyExists = errors.New("activity already exists")

	// ErrInvalidARN is returned when the Amazon Resource Name (ARN) is not valid.
	ErrInvalidARN = errors.New("invalid ARN")

	// ErrInvalidDefinition is returned when the state machine definition
	// is not valid.
	ErrInvalidDefinition = errors.New("invalid state machine definition")

	// ErrTaskNotFound is returned when the specified task does not exist.
	ErrTaskNotFound = errors.New("task not found")

	// ErrTaskTimeout is returned when the task has timed out.
	ErrTaskTimeout = errors.New("task timed out")

	// ErrActivityQueueFull is returned when the activity task queue is full.
	ErrActivityQueueFull = errors.New("activity task queue is full")

	// ErrTaskNotRunning is returned when the task is not running.
	ErrTaskNotRunning = errors.New("task is not running")
)
