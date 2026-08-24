package sfn

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// Error variables for common Step Functions error conditions. The set
// follows the Smithy AWSStepFunctions model: every error code raised by a
// handler must be one the model attaches to that operation, so call sites
// use these named errors (or the New* constructors below) instead of
// ad-hoc NewAWSError calls.
var (
	// ErrStateMachineDoesNotExist is returned when the specified state machine does not exist.
	ErrStateMachineDoesNotExist = awserrors.NewAWSError("StateMachineDoesNotExist", "State Machine Does not exist", http.StatusBadRequest)
	// ErrExecutionDoesNotExist is returned when the specified execution does not exist.
	ErrExecutionDoesNotExist = awserrors.NewAWSError("ExecutionDoesNotExist", "Execution Does not exist", http.StatusBadRequest)
	// ErrActivityDoesNotExist is returned when the specified activity does not exist.
	ErrActivityDoesNotExist = awserrors.NewAWSError("ActivityDoesNotExist", "Activity Does not exist", http.StatusBadRequest)
	// ErrInvalidArn is returned when an ARN is invalid.
	ErrInvalidArn = awserrors.NewAWSError("InvalidArn", "Invalid Arn", http.StatusBadRequest)
	// ErrInvalidDefinition is returned when the state machine definition is invalid.
	ErrInvalidDefinition = awserrors.NewAWSError("InvalidDefinition", "Invalid Definition", http.StatusBadRequest)
	// ErrInvalidName is returned when a name is invalid.
	ErrInvalidName = awserrors.NewAWSError("InvalidName", "Invalid Name", http.StatusBadRequest)
	// ErrTaskTimedOut is returned when a task timed out.
	ErrTaskTimedOut = awserrors.NewAWSError("TaskTimedOut", "Task timed out", http.StatusBadRequest)
	// ErrTaskDoesNotExist is returned when a task token does not exist or
	// the task it refers to is no longer in a state that accepts reports.
	ErrTaskDoesNotExist = awserrors.NewAWSError("TaskDoesNotExist", "Task does not exist", http.StatusBadRequest)
)

// NewStateMachineDoesNotExist creates a new StateMachineDoesNotExist error.
func NewStateMachineDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("StateMachineDoesNotExist", message, http.StatusBadRequest)
}

// NewExecutionDoesNotExist creates a new ExecutionDoesNotExist error.
func NewExecutionDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ExecutionDoesNotExist", message, http.StatusBadRequest)
}

// NewActivityDoesNotExist creates a new ActivityDoesNotExist error.
func NewActivityDoesNotExist(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ActivityDoesNotExist", message, http.StatusBadRequest)
}

// NewActivityAlreadyExists creates a new ActivityAlreadyExists error
// (CreateActivity).
func NewActivityAlreadyExists(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ActivityAlreadyExists", message, http.StatusBadRequest)
}

// NewInvalidArnException creates a new InvalidArn error.
func NewInvalidArnException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidArn", message, http.StatusBadRequest)
}

// NewInvalidDefinitionException creates a new InvalidDefinition error.
func NewInvalidDefinitionException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidDefinition", message, http.StatusBadRequest)
}

// NewInvalidName creates a new InvalidName error.
func NewInvalidName(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidName", message, http.StatusBadRequest)
}

// NewInvalidToken creates a new InvalidToken error (List operations,
// SendTask* task tokens).
func NewInvalidToken(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidToken", message, http.StatusBadRequest)
}

// NewStateMachineAlreadyExists creates a new StateMachineAlreadyExists
// error (CreateStateMachine).
func NewStateMachineAlreadyExists(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("StateMachineAlreadyExists", message, http.StatusBadRequest)
}

// NewExecutionAlreadyExists creates a new ExecutionAlreadyExists error
// (StartExecution).
func NewExecutionAlreadyExists(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ExecutionAlreadyExists", message, http.StatusBadRequest)
}

// NewResourceNotFound creates a new ResourceNotFound error. The Smithy
// model attaches @httpError(404) to this shape, and it is raised by the
// alias, version and map-run operations where a referenced resource does
// not exist.
func NewResourceNotFound(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ResourceNotFound", message, http.StatusNotFound)
}

// NewConflictException creates a new ConflictException error. Used for
// already-exists races on CreateStateMachineAlias and concurrent
// publish/update requests.
func NewConflictException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ConflictException", message, http.StatusConflict)
}

// NewValidationException creates a new ValidationException error. The
// Smithy model attaches it to most operations as the generic
// input-constraint failure.
func NewValidationException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ValidationException", message, http.StatusBadRequest)
}

// NewStateMachineTypeNotSupported creates a new
// StateMachineTypeNotSupported error (CreateStateMachine,
// StartSyncExecution, ListExecutions).
func NewStateMachineTypeNotSupported(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("StateMachineTypeNotSupported", message, http.StatusBadRequest)
}

// NewExecutionNotRedrivable creates a new ExecutionNotRedrivable error
// (RedriveExecution).
func NewExecutionNotRedrivable(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ExecutionNotRedrivable", message, http.StatusBadRequest)
}

// NewInvalidLoggingConfiguration creates a new InvalidLoggingConfiguration
// error (CreateStateMachine, UpdateStateMachine, CreateActivity).
func NewInvalidLoggingConfiguration(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidLoggingConfiguration", message, http.StatusBadRequest)
}

// NewInvalidEncryptionConfiguration creates a new
// InvalidEncryptionConfiguration error (CreateStateMachine,
// UpdateStateMachine, CreateActivity).
func NewInvalidEncryptionConfiguration(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidEncryptionConfiguration", message, http.StatusBadRequest)
}

// NewInvalidTracingConfiguration creates a new InvalidTracingConfiguration
// error (CreateStateMachine, UpdateStateMachine).
func NewInvalidTracingConfiguration(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidTracingConfiguration", message, http.StatusBadRequest)
}

// NewTooManyTags creates a new TooManyTags error (CreateStateMachine,
// CreateActivity, TagResource). The tagging quota is fifty tags per
// resource.
func NewTooManyTags(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("TooManyTags", message, http.StatusBadRequest)
}

// NewMissingRequiredParameter creates a new MissingRequiredParameter error
// (UpdateStateMachine).
func NewMissingRequiredParameter(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("MissingRequiredParameter", message, http.StatusBadRequest)
}

// NewInvalidExecutionInput creates a new InvalidExecutionInput error
// (StartExecution, StartSyncExecution): the provided JSON input data is
// not valid.
func NewInvalidExecutionInput(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidExecutionInput", message, http.StatusBadRequest)
}

// NewInvalidOutput creates a new InvalidOutput error (SendTaskSuccess):
// the provided JSON output data is not valid.
func NewInvalidOutput(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidOutput", message, http.StatusBadRequest)
}

// NewActivityLimitExceeded creates a new ActivityLimitExceeded error
// (CreateActivity): the maximum number of activities has been reached.
func NewActivityLimitExceeded(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ActivityLimitExceeded", message, http.StatusBadRequest)
}

// NewActivityWorkerLimitExceeded creates a new ActivityWorkerLimitExceeded
// error (GetActivityTask): the maximum number of activity workers has been
// reached.
func NewActivityWorkerLimitExceeded(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ActivityWorkerLimitExceeded", message, http.StatusBadRequest)
}

// NewExecutionLimitExceeded creates a new ExecutionLimitExceeded error
// (StartExecution, RedriveExecution): the maximum number of running
// executions has been reached.
func NewExecutionLimitExceeded(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ExecutionLimitExceeded", message, http.StatusBadRequest)
}

// NewStateMachineLimitExceeded creates a new StateMachineLimitExceeded
// error (CreateStateMachine): the maximum number of state machines has
// been reached.
func NewStateMachineLimitExceeded(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("StateMachineLimitExceeded", message, http.StatusBadRequest)
}

// NewStateMachineDeleting creates a new StateMachineDeleting error: the
// specified state machine is being deleted.
func NewStateMachineDeleting(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("StateMachineDeleting", message, http.StatusBadRequest)
}

// NewKmsAccessDeniedException creates a new KmsAccessDeniedException: the
// KMS key policy or API caller lacks the required permissions.
func NewKmsAccessDeniedException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("KmsAccessDeniedException", message, http.StatusBadRequest)
}

// NewKmsInvalidStateException creates a new KmsInvalidStateException: the
// KMS key is not in a valid state, for example Disabled or Deleted.
func NewKmsInvalidStateException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("KmsInvalidStateException", message, http.StatusBadRequest)
}

// NewKmsThrottlingException creates a new KmsThrottlingException: KMS
// returned ThrottlingException for a call made on behalf of the caller.
func NewKmsThrottlingException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("KmsThrottlingException", message, http.StatusBadRequest)
}

// NewServiceQuotaExceededException creates a new
// ServiceQuotaExceededException error. The Smithy model attaches
// @httpError(402) to this shape, unlike the default 400 of the other
// client errors.
func NewServiceQuotaExceededException(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("ServiceQuotaExceededException", message, http.StatusPaymentRequired)
}
