package sfn

import (
	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// SFNError represents an error returned by the Step Functions service.
type SFNError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error for error chaining.
func (e *SFNError) Unwrap() error {
	return e.AWSError
}

// Error variables for common Step Functions error conditions.
var (
	ErrStateMachineDoesNotExist = &SFNError{awserrors.NewAWSError("StateMachineDoesNotExist", "State Machine Does not exist", 400)}
	ErrExecutionDoesNotExist    = &SFNError{awserrors.NewAWSError("ExecutionDoesNotExist", "Execution Does not exist", 400)}
	ErrActivityDoesNotExist     = &SFNError{awserrors.NewAWSError("ActivityDoesNotExist", "Activity Does not exist", 400)}
	ErrInvalidArn               = &SFNError{awserrors.NewAWSError("InvalidArn", "Invalid Arn", 400)}
	ErrInvalidDefinition        = &SFNError{awserrors.NewAWSError("InvalidDefinition", "Invalid Definition", 400)}
	ErrInvalidName              = &SFNError{awserrors.NewAWSError("InvalidName", "Invalid Name", 400)}
	ErrTaskTimedOut             = &SFNError{awserrors.NewAWSError("TaskTimedOut", "Task timed out", 400)}
	ErrTaskNotRunning           = &SFNError{awserrors.NewAWSError("TaskNotRunning", "Task is not running", 400)}
	ErrQueryEvaluationError     = &SFNError{awserrors.NewAWSError("States.QueryEvaluationError", "The query evaluation failed", 400)}
)

// NewStateMachineDoesNotExist creates a new StateMachineDoesNotExist error.
func NewStateMachineDoesNotExist(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("StateMachineDoesNotExist", message, 400)}
}

// NewExecutionDoesNotExist creates a new ExecutionDoesNotExist error.
func NewExecutionDoesNotExist(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("ExecutionDoesNotExist", message, 400)}
}

// NewActivityDoesNotExist creates a new ActivityDoesNotExist error.
func NewActivityDoesNotExist(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("ActivityDoesNotExist", message, 400)}
}

// NewInvalidArnException creates a new InvalidArn error.
func NewInvalidArnException(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("InvalidArn", message, 400)}
}

// NewInvalidDefinitionException creates a new InvalidDefinition error.
func NewInvalidDefinitionException(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("InvalidDefinition", message, 400)}
}

// NewInvalidName creates a new InvalidName error.
func NewInvalidName(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("InvalidName", message, 400)}
}

// NewInvalidExecutionType creates a new InvalidExecutionType error.
func NewInvalidExecutionType(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("InvalidExecutionType", message, 400)}
}

// NewMapRunDoesNotExist creates a new MapRunDoesNotExist error.
func NewMapRunDoesNotExist(message string) *SFNError {
	return &SFNError{awserrors.NewAWSError("MapRunDoesNotExist", message, 400)}
}
