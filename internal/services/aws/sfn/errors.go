package sfn

import (
	awserrors "vorpalstacks/internal/common/errors"
)

// Error variables for common Step Functions error conditions.
var (
	// ErrStateMachineDoesNotExist is returned when the specified state machine does not exist.
	ErrStateMachineDoesNotExist = awserrors.NewAWSError("StateMachineDoesNotExist", "State Machine Does not exist", 400)
	// ErrExecutionDoesNotExist is returned when the specified execution does not exist.
	ErrExecutionDoesNotExist = awserrors.NewAWSError("ExecutionDoesNotExist", "Execution Does not exist", 400)
	// ErrActivityDoesNotExist is returned when the specified activity does not exist.
	ErrActivityDoesNotExist = awserrors.NewAWSError("ActivityDoesNotExist", "Activity Does not exist", 400)
	// ErrInvalidArn is returned when an ARN is invalid.
	ErrInvalidArn = awserrors.NewAWSError("InvalidArn", "Invalid Arn", 400)
	// ErrInvalidDefinition is returned when the state machine definition is invalid.
	ErrInvalidDefinition = awserrors.NewAWSError("InvalidDefinition", "Invalid Definition", 400)
	// ErrInvalidName is returned when a name is invalid.
	ErrInvalidName = awserrors.NewAWSError("InvalidName", "Invalid Name", 400)
	// ErrTaskTimedOut is returned when a task timed out.
	ErrTaskTimedOut = awserrors.NewAWSError("TaskTimedOut", "Task timed out", 400)
	// ErrTaskNotRunning is returned when a task is not running.
	ErrTaskNotRunning = awserrors.NewAWSError("TaskNotRunning", "Task is not running", 400)
	// ErrQueryEvaluationError is returned when a query evaluation fails.
	ErrQueryEvaluationError = awserrors.NewAWSError("States.QueryEvaluationError", "The query evaluation failed", 400)
)

// NewStateMachineDoesNotExist creates a new StateMachineDoesNotExist error.
func NewStateMachineDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("StateMachineDoesNotExist", message, 400)
}

// NewExecutionDoesNotExist creates a new ExecutionDoesNotExist error.
func NewExecutionDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ExecutionDoesNotExist", message, 400)
}

// NewActivityDoesNotExist creates a new ActivityDoesNotExist error.
func NewActivityDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ActivityDoesNotExist", message, 400)
}

// NewInvalidArnException creates a new InvalidArn error.
func NewInvalidArnException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidArn", message, 400)
}

// NewInvalidDefinitionException creates a new InvalidDefinition error.
func NewInvalidDefinitionException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidDefinition", message, 400)
}

// NewInvalidName creates a new InvalidName error.
func NewInvalidName(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidName", message, 400)
}

// NewInvalidExecutionType creates a new InvalidExecutionType error.
func NewInvalidExecutionType(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidExecutionType", message, 400)
}

// NewMapRunDoesNotExist creates a new MapRunDoesNotExist error.
func NewMapRunDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("MapRunDoesNotExist", message, 400)
}
