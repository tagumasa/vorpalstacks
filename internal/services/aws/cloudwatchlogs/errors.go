package cloudwatchlogs

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
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
	// ErrDestinationNotFound is returned when a destination is not found.
	ErrDestinationNotFound = NewLogsError("ResourceNotFoundException", "Destination not found", http.StatusNotFound)
	// ErrDestinationAlreadyExists is returned when a destination already exists.
	ErrDestinationAlreadyExists = NewLogsError("ResourceAlreadyExistsException", "Destination already exists", http.StatusConflict)
)

// storeErrorMappings maps store-level sentinel errors to CloudWatch Logs API errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: logsstore.ErrLogGroupNotFound, AWS: ErrLogGroupNotFound},
	{Store: logsstore.ErrLogGroupAlreadyExists, AWS: ErrLogGroupAlreadyExists},
	{Store: logsstore.ErrLogStreamNotFound, AWS: ErrLogStreamNotFound},
	{Store: logsstore.ErrLogStreamAlreadyExists, AWS: ErrLogStreamAlreadyExists},
	{Store: logsstore.ErrMetricFilterNotFound, AWS: ErrMetricFilterNotFound},
	{Store: logsstore.ErrMetricFilterAlreadyExists, AWS: awserrors.NewAWSError("ResourceAlreadyExistsException", "metric filter already exists", http.StatusConflict)},
	{Store: logsstore.ErrResourceNotFound, AWS: awserrors.NewAWSError("ResourceNotFoundException", "resource not found", http.StatusNotFound)},
	{Store: logsstore.ErrResourceAlreadyExists, AWS: awserrors.NewAWSError("ResourceAlreadyExistsException", "resource already exists", http.StatusConflict)},
	{Store: logsstore.ErrDataAlreadyAccepted, AWS: awserrors.NewAWSError("DataAlreadyAcceptedException", "data already accepted", http.StatusBadRequest)},
	{Store: logsstore.ErrInvalidSequenceToken, AWS: awserrors.NewAWSError("InvalidSequenceTokenException", "invalid sequence token", http.StatusBadRequest)},
	{Store: logsstore.ErrSubscriptionFilterNotFound, AWS: awserrors.NewAWSError("ResourceNotFoundException", "subscription filter not found", http.StatusNotFound)},
	{Store: logsstore.ErrDestinationNotFound, AWS: ErrDestinationNotFound},
	{Store: logsstore.ErrDestinationAlreadyExists, AWS: ErrDestinationAlreadyExists},
	{Store: logsstore.ErrLimitExceeded, AWS: ErrLimitExceeded},
	{Store: logsstore.ErrInvalidPaginationToken, AWS: ErrInvalidParameter},
}

// mapStoreError converts a store error into an appropriate CloudWatch Logs API error.
func mapStoreError(err error) error {
	return awserrors.MapStoreError(err, storeErrorMappings)
}
