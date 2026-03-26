package cloudwatch

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

var (
	// ErrInvalidParameter is returned when an invalid parameter is provided.
	ErrInvalidParameter = awserrors.NewInvalidParameterException("Invalid parameter")
	// ErrResourceNotFound is returned when a resource is not found.
	ErrResourceNotFound = awserrors.NewResourceNotFoundException("Resource", "")
	// ErrResourceAlreadyExists is returned when a resource already exists.
	ErrResourceAlreadyExists = awserrors.NewResourceAlreadyExistsException("Resource")
	// ErrInvalidSequenceToken is returned when an invalid sequence token is provided.
	ErrInvalidSequenceToken = awserrors.NewAWSError("InvalidSequenceTokenException", "Invalid sequence token", http.StatusBadRequest)
	// ErrDataAlreadyAccepted is returned when data has already been accepted.
	ErrDataAlreadyAccepted = awserrors.NewAWSError("DataAlreadyAcceptedException", "Data already accepted", http.StatusBadRequest)
	// ErrOperationAborted is returned when an operation is aborted.
	ErrOperationAborted = awserrors.NewAWSError("OperationAbortedException", "Operation aborted", http.StatusBadRequest)
	// ErrLimitExceeded is returned when a limit is exceeded.
	ErrLimitExceeded = awserrors.NewLimitExceededException("Limit exceeded")
	// ErrServiceUnavailable is returned when the service is unavailable.
	ErrServiceUnavailable = awserrors.NewServiceUnavailableException("Service unavailable")
	// ErrAlarmNotFound is returned when an alarm is not found.
	ErrAlarmNotFound = awserrors.NewResourceNotFoundException("Alarm", "")
	// ErrAlarmLimitExceeded is returned when the alarm limit is exceeded.
	ErrAlarmLimitExceeded = awserrors.NewLimitExceededException("Alarm limit exceeded")
)

// NewInvalidSequenceTokenError creates an error for an invalid sequence token.
func NewInvalidSequenceTokenError(expectedToken string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidSequenceTokenException", "The given sequenceToken is invalid. The next expected sequenceToken is: "+expectedToken, http.StatusBadRequest)
}
