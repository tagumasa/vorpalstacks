package timestreamwrite

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

var (
	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFoundException", "The requested resource could not be found.", http.StatusNotFound)
	// ErrConflictException is returned when attempting to create a resource that already exists.
	// Smithy declares ConflictException (HTTP 409) for CreateDatabase, CreateTable, and
	// CreateBatchLoadTask. ResourceAlreadyExistsException does not exist in the SDK.
	ErrConflictException = awserrors.NewAWSError("ConflictException", "The resource already exists.", http.StatusConflict)
	// ErrRejectedRecordsException is returned when WriteRecords rejects all records.
	// Smithy declares RejectedRecordsException (HTTP 419) for WriteRecords.
	ErrRejectedRecordsException = awserrors.NewAWSError("RejectedRecordsException", "One or more records were rejected.", 419)
	// ErrValidationException is returned when validation fails.
	ErrValidationException = awserrors.NewAWSError("ValidationException", "The input fails to satisfy the constraints specified by an AWS service.", http.StatusBadRequest)
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = awserrors.NewAWSError("AccessDeniedException", "You are not authorised to perform this action.", http.StatusForbidden)
	// ErrThrottling is returned when the request is throttled.
	ErrThrottling = awserrors.NewAWSError("ThrottlingException", "The request was denied due to request throttling.", http.StatusTooManyRequests)
	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = awserrors.NewAWSError("InternalServerException", "Internal server error.", http.StatusInternalServerError)
	// ErrServiceQuotaExceeded is returned when a service quota is exceeded.
	ErrServiceQuotaExceeded = awserrors.NewAWSError("ServiceQuotaExceededException", "The request exceeded the service quota.", http.StatusTooManyRequests)
)

// storeErrorMappings maps store-level sentinel errors to Timestream API errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: tsstore.ErrDatabaseNotFound, AWS: ErrResourceNotFound},
	{Store: tsstore.ErrTableNotFound, AWS: ErrResourceNotFound},
	{Store: tsstore.ErrBatchLoadTaskNotFound, AWS: ErrResourceNotFound},
	{Store: tsstore.ErrDatabaseAlreadyExists, AWS: ErrConflictException},
	{Store: tsstore.ErrTableAlreadyExists, AWS: ErrConflictException},
	{Store: tsstore.ErrBatchLoadTaskAlreadyExists, AWS: ErrConflictException},
	{Store: tsstore.ErrDatabaseNotEmpty, AWS: ErrValidationException},
}
