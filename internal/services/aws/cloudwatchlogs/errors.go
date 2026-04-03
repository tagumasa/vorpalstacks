package cloudwatchlogs

import (
	"errors"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// NewLogsError creates a new CloudWatch Logs error.
func NewLogsError(code, message string, statusCode int) *awserrors.AWSError {
	return awserrors.NewAWSError(code, message, statusCode)
}

var (
	// ErrLogGroupNotFound is returned when a log group is not found.
	ErrLogGroupNotFound = NewLogsError("ResourceNotFoundException", "Log group not found", http.StatusNotFound)
	// ErrLogGroupAlreadyExists is returned when a log group already exists.
	ErrLogGroupAlreadyExists = NewLogsError("ResourceAlreadyExistsException", "Log group already exists", http.StatusConflict)
	// ErrLogStreamNotFound is returned when a log stream is not found.
	ErrLogStreamNotFound = NewLogsError("ResourceNotFoundException", "Log stream not found", http.StatusNotFound)
	// ErrLogStreamAlreadyExists is returned when a log stream already exists.
	ErrLogStreamAlreadyExists = NewLogsError("ResourceAlreadyExistsException", "Log stream already exists", http.StatusConflict)
	// ErrMetricFilterNotFound is returned when a metric filter is not found.
	ErrMetricFilterNotFound = NewLogsError("ResourceNotFoundException", "Metric filter not found", http.StatusNotFound)
	// ErrInvalidParameter is returned when an invalid parameter is provided.
	ErrInvalidParameter = NewLogsError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = NewLogsError("MissingParameterException", "Missing required parameter", http.StatusBadRequest)
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = NewLogsError("AccessDeniedException", "Access denied", http.StatusForbidden)
	// ErrLimitExceeded is returned when a limit is exceeded.
	ErrLimitExceeded = NewLogsError("LimitExceededException", "Limit exceeded", http.StatusTooManyRequests)
	// ErrOperationAborted is returned when an operation is aborted.
	ErrOperationAborted = NewLogsError("OperationAbortedException", "Operation aborted", http.StatusBadRequest)
	// ErrServiceUnavailable is returned when the service is unavailable.
	ErrServiceUnavailable = NewLogsError("ServiceUnavailableException", "Service unavailable", http.StatusServiceUnavailable)
	// ErrDestinationNotFound is returned when a destination is not found.
	ErrDestinationNotFound = NewLogsError("ResourceNotFoundException", "Destination not found", http.StatusNotFound)
	// ErrDestinationAlreadyExists is returned when a destination already exists.
	ErrDestinationAlreadyExists = NewLogsError("ResourceAlreadyExistsException", "Destination already exists", http.StatusConflict)
)

// mapStoreError converts a store error into an appropriate CloudWatch Logs API error.
func mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, logsstore.ErrLogGroupNotFound):
		return ErrLogGroupNotFound
	case errors.Is(err, logsstore.ErrLogGroupAlreadyExists):
		return ErrLogGroupAlreadyExists
	case errors.Is(err, logsstore.ErrLogStreamNotFound):
		return ErrLogStreamNotFound
	case errors.Is(err, logsstore.ErrLogStreamAlreadyExists):
		return ErrLogStreamAlreadyExists
	case errors.Is(err, logsstore.ErrMetricFilterNotFound):
		return ErrMetricFilterNotFound
	case errors.Is(err, logsstore.ErrMetricFilterAlreadyExists):
		return NewLogsError("ResourceAlreadyExistsException", "metric filter already exists", http.StatusConflict)
	case errors.Is(err, logsstore.ErrResourceNotFound):
		return NewLogsError("ResourceNotFoundException", "resource not found", http.StatusNotFound)
	case errors.Is(err, logsstore.ErrResourceAlreadyExists):
		return NewLogsError("ResourceAlreadyExistsException", "resource already exists", http.StatusConflict)
	case errors.Is(err, logsstore.ErrDataAlreadyAccepted):
		return NewLogsError("DataAlreadyAcceptedException", "data already accepted", http.StatusBadRequest)
	case errors.Is(err, logsstore.ErrInvalidSequenceToken):
		return NewLogsError("InvalidSequenceTokenException", "invalid sequence token", http.StatusBadRequest)
	case errors.Is(err, logsstore.ErrSubscriptionFilterNotFound):
		return NewLogsError("ResourceNotFoundException", "subscription filter not found", http.StatusNotFound)
	case errors.Is(err, logsstore.ErrDestinationNotFound):
		return ErrDestinationNotFound
	case errors.Is(err, logsstore.ErrDestinationAlreadyExists):
		return ErrDestinationAlreadyExists
	default:
		return err
	}
}
