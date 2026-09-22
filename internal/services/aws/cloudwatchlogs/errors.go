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

// CloudWatch Logs is an awsJson1_1 service, so its client errors travel
// with the code in the response body and default to HTTP 400. The model's
// httpError census is not empty, though: TooManyTagsException carries
// httpError 400 and InternalServerException httpError 500 — any future
// InternalServerException mapping must carry 500, never the family
// default.
var (
	// ErrLogGroupNotFound is returned when a log group is not found.
	ErrLogGroupNotFound = NewLogsError("ResourceNotFoundException", "Log group not found", http.StatusBadRequest)
	// ErrLogGroupAlreadyExists is returned when a log group already exists.
	ErrLogGroupAlreadyExists = NewLogsError("ResourceAlreadyExistsException", "Log group already exists", http.StatusBadRequest)
	// ErrLogStreamNotFound is returned when a log stream is not found.
	ErrLogStreamNotFound = NewLogsError("ResourceNotFoundException", "Log stream not found", http.StatusBadRequest)
	// ErrLogStreamAlreadyExists is returned when a log stream already exists.
	ErrLogStreamAlreadyExists = NewLogsError("ResourceAlreadyExistsException", "Log stream already exists", http.StatusBadRequest)
	// ErrMetricFilterNotFound is returned when a metric filter is not found.
	ErrMetricFilterNotFound = NewLogsError("ResourceNotFoundException", "Metric filter not found", http.StatusBadRequest)
	// ErrInvalidParameter is returned when an invalid parameter is provided.
	ErrInvalidParameter = NewLogsError("InvalidParameterException", "Invalid parameter", http.StatusBadRequest)
	// ErrLimitExceeded is returned when a limit is exceeded.
	ErrLimitExceeded = NewLogsError("LimitExceededException", "Limit exceeded", http.StatusBadRequest)
	// ErrDestinationNotFound is returned when a destination is not found.
	// No destination-already-exists sentinel exists: PutDestination is the
	// documented create-or-update, so the conflict identity never fires.
	ErrDestinationNotFound = NewLogsError("ResourceNotFoundException", "Destination not found", http.StatusBadRequest)
)

// errRequiredMember is the identity a missing required member rejects
// with: InvalidParameterException, the parameter error every implemented
// operation's declared list carries. The legacy tag operations omit it
// from their lists; the least-contradictory-identity adjudication for
// those raw-HTTP-only rows lives in the plan register.
func errRequiredMember(member string) error {
	return NewLogsError("InvalidParameterException",
		"The "+member+" member is required", 400)
}

// errValidationMember is the identity the scheduled-query operation
// family rejects a missing member with: their declared lists carry
// ValidationException (with AccessDenied, Conflict, InternalServer,
// ResourceNotFound, Throttling) — not InvalidParameterException.
func errValidationMember(member string) error {
	return NewLogsError("ValidationException",
		"The "+member+" member is required", 400)
}

// storeErrorMappings maps store-level sentinel errors to CloudWatch Logs API errors.
//
// The two sentinel vocabularies are deliberate layering, not drift: the
// store layer speaks storage-identity sentinels (no wire knowledge),
// this layer owns the wire identities, and this table is the
// hand-maintained bridge between them — the established per-service
// pattern. The discipline the bridge imposes: every row's store sentinel
// must have a producer somewhere in the store (a row whose store side
// never fires is a dead pair — both sentinels and the row go), and every
// wire identity the table names must be one this file declares.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: logsstore.ErrLogGroupNotFound, AWS: ErrLogGroupNotFound},
	{Store: logsstore.ErrLogGroupAlreadyExists, AWS: ErrLogGroupAlreadyExists},
	{Store: logsstore.ErrLogStreamNotFound, AWS: ErrLogStreamNotFound},
	{Store: logsstore.ErrLogStreamAlreadyExists, AWS: ErrLogStreamAlreadyExists},
	{Store: logsstore.ErrMetricFilterNotFound, AWS: ErrMetricFilterNotFound},
	{Store: logsstore.ErrResourceNotFound, AWS: awserrors.NewAWSError("ResourceNotFoundException", "resource not found", http.StatusBadRequest)},
	{Store: logsstore.ErrSubscriptionFilterNotFound, AWS: awserrors.NewAWSError("ResourceNotFoundException", "subscription filter not found", http.StatusBadRequest)},
	{Store: logsstore.ErrDestinationNotFound, AWS: ErrDestinationNotFound},
	{Store: logsstore.ErrLimitExceeded, AWS: ErrLimitExceeded},
	{Store: logsstore.ErrInvalidPaginationToken, AWS: ErrInvalidParameter},
	// The tag store's cumulative ceiling needs no row: its construction
	// budget already carries the wire identity (TooManyTagsException,
	// declared on TagResource alone), and MapStoreError passes the
	// already-shaped error through unchanged. The legacy tag writers
	// convert it at their Core — their declared lists carry only
	// InvalidParameterException.
	// "If you attempt to delete a log group with deletion protection
	// enabled, you receive a ValidationException with the message:
	// 'Cannot delete log group with deletion protection enabled. Disable
	// deletion protection first.'" (user guide, protecting log groups
	// from deletion) — OperationAbortedException is the unrelated
	// concurrent-update conflict identity.
	{Store: logsstore.ErrLogGroupDeletionProtected, AWS: NewLogsError("ValidationException",
		"Cannot delete log group with deletion protection enabled. Disable deletion protection first.", http.StatusBadRequest)},
}

// mapStoreError converts a store error into an appropriate CloudWatch Logs API error.
func mapStoreError(err error) error {
	return awserrors.MapStoreError(err, storeErrorMappings)
}
