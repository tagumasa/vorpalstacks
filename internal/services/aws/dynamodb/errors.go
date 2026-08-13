package dynamodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// APIError represents a DynamoDB API error.
type APIError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *APIError) Unwrap() error {
	return e.AWSError
}

// NewAPIError creates a new DynamoDB API error.
func NewAPIError(code, message string, httpStatus int) *APIError {
	return &APIError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

// TransactionCanceledError represents a transaction cancellation error.
type TransactionCanceledError struct {
	*APIError
	CancellationReasons []CancellationReason
}

// CancellationReason represents a reason for transaction cancellation.
type CancellationReason struct {
	Code string
	Item map[string]interface{}
}

// NewTransactionCanceledError creates a new TransactionCanceledError.
func NewTransactionCanceledError(message string, reasons []CancellationReason) *TransactionCanceledError {
	return &TransactionCanceledError{
		APIError:            NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionCanceledException", message, http.StatusBadRequest),
		CancellationReasons: reasons,
	}
}

// ToJSON serialises the error to JSON format.
func (e *TransactionCanceledError) ToJSON() string {
	type cancellationReasonJSON struct {
		Item map[string]interface{} `json:"Item,omitempty"`
		Code string                 `json:"Code"`
	}

	type errorJSON struct {
		Type                string                   `json:"__type"`
		Message             string                   `json:"message"`
		CancellationReasons []cancellationReasonJSON `json:"CancellationReasons"`
	}

	reasons := make([]cancellationReasonJSON, len(e.CancellationReasons))
	for i, r := range e.CancellationReasons {
		reasons[i] = cancellationReasonJSON{
			Item: r.Item,
			Code: r.Code,
		}
	}

	resp := errorJSON{
		Type:                e.APIError.AWSError.Code,
		Message:             e.APIError.AWSError.Message,
		CancellationReasons: reasons,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return fmt.Sprintf(`{"__type":"%s","message":"%s"}`, e.APIError.AWSError.Code, e.APIError.AWSError.Message)
	}
	return string(b)
}

// Predefined DynamoDB error variables.
var (
	// ErrResourceNotFound is returned when a requested DynamoDB resource does not exist.
	// DynamoDB uses the awsJson1_0 protocol where ALL client errors are HTTP 400.
	ErrResourceNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Requested resource not found", http.StatusBadRequest)
	// ErrResourceAlreadyExists is returned when a resource already exists with the same identifier.
	ErrResourceAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Resource already exists", http.StatusConflict)
	// ErrInvalidParameter is returned when a request parameter fails validation.
	ErrInvalidParameter = NewAPIError("com.amazon.coral.validate#ValidationException", "Invalid parameter", http.StatusBadRequest)

	// ErrTypeMismatch is a typed sentinel returned by the expression
	// evaluator when an ADD or SET operation encounters incompatible
	// attribute types. Callers should use errors.Is to detect it.
	ErrTypeMismatch = errors.New("TYPE_MISMATCH: Type mismatch for attribute to update")
	// ErrInternal is returned when an internal error occurs during stream operations.
	ErrInternal = NewAPIError("com.amazonaws.dynamodb.v20120810#InternalServerError", "Internal error", http.StatusInternalServerError)
	// ErrTableNotFound is returned when the specified table does not exist.
	// Uses the general ResourceNotFoundException code so that callers checked
	// against ResourceNotFoundException (the vast majority of DynamoDB ops)
	// match. Operations that the Smithy model declares as throwing the
	// distinct TableNotFoundException (CreateBackup, CreateGlobalTable,
	// DescribeContinuousBackups, ExportTableToPointInTime,
	// RestoreTableToPointInTime, UpdateContinuousBackups, UpdateGlobalTable)
	// must use ErrTableNotFoundException below instead.
	ErrTableNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Requested resource not found: Table not found", http.StatusBadRequest)
	// ErrTableNotFoundException is the Smithy-specific TableNotFoundException
	// used by the seven operations whose Smithy model declares it. Most
	// operations throw ResourceNotFoundException via ErrTableNotFound above.
	ErrTableNotFoundException = NewAPIError("com.amazonaws.dynamodb.v20120810#TableNotFoundException", "Requested resource not found: Table not found", http.StatusBadRequest)
	// ErrTableAlreadyExists is returned when CreateTable or ImportTable
	// targets a name that is already in use. Uses the general
	// ResourceInUseException code per Smithy (CreateTable/ImportTable declare
	// ResourceInUseException, not TableAlreadyExistsException).
	ErrTableAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table already exists", http.StatusConflict)
	// ErrTableAlreadyExistsException is the Smithy-specific
	// TableAlreadyExistsException used by RestoreTableFromBackup and
	// RestoreTableToPointInTime, whose Smithy models declare it instead of
	// the general ResourceInUseException.
	ErrTableAlreadyExistsException = NewAPIError("com.amazonaws.dynamodb.v20120810#TableAlreadyExistsException", "Table already exists", http.StatusConflict)
	// ErrTableInUseException is the Smithy-specific TableInUseException
	// declared by CreateBackup, RestoreTableFromBackup, and
	// RestoreTableToPointInTime. It signals that the target table is in a
	// transitional (CREATING/DELETING/UPDATING) state.
	ErrTableInUseException = NewAPIError("com.amazonaws.dynamodb.v20120810#TableInUseException", "Table is in use; it is being created or deleted", http.StatusConflict)
	// ErrInvalidEndpoint is the Smithy InvalidEndpointException (declared by
	// 44 DynamoDB operations). Smithy sets httpError 421 (Misdirected
	// Request). Currently defined for completeness; the vorpalstacks endpoint
	// resolver does not yet surface this condition.
	ErrInvalidEndpoint = NewAPIError("com.amazonaws.dynamodb.v20120810#InvalidEndpointException", "Invalid endpoint", 421)
	// ErrRequestLimitExceeded is the Smithy RequestLimitExceeded (declared by
	// 13 data-plane operations). Currently defined for completeness;
	// throughput-based throttling is handled by ErrThrottling /
	// ErrProvisionedThroughputExceeded for now.
	ErrRequestLimitExceeded = NewAPIError("com.amazonaws.dynamodb.v20120810#RequestLimitExceeded", "Request limit exceeded; throughput quota for the account has been exceeded", http.StatusBadRequest)
	// ErrTableNotActive is returned when the table is not in ACTIVE state.
	ErrTableNotActive = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is not in ACTIVE state", http.StatusBadRequest)
	// ErrTableDeletionProtected is returned when deletion protection is enabled on the table.
	ErrTableDeletionProtected = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is protected from deletion", http.StatusBadRequest)
	// ErrConditionalCheckFailed is returned when a conditional write operation fails.
	ErrConditionalCheckFailed = NewAPIError("com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException", "The conditional request failed", http.StatusBadRequest)
	// ErrItemNotFound is returned when the specified item does not exist in the table.
	ErrItemNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Item not found", http.StatusBadRequest)
	// ErrInvalidKey is returned when the provided key is invalid.
	ErrInvalidKey = NewAPIError("com.amazon.coral.validate#ValidationException", "Invalid key", http.StatusBadRequest)
	// ErrMissingKey is returned when a required key is missing from the request.
	ErrMissingKey = NewAPIError("com.amazon.coral.validate#ValidationException", "Missing required key in request", http.StatusBadRequest)
	// ErrBackupNotFound is returned when the specified backup does not exist.
	ErrBackupNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#BackupNotFoundException", "Backup not found", http.StatusBadRequest)
	// ErrBackupAlreadyExists is returned when a backup already exists with the same name.
	ErrBackupAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#BackupInUseException", "Backup already exists", http.StatusConflict)
	// ErrGlobalTableNotFound is returned when the specified global table does not exist.
	ErrGlobalTableNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#GlobalTableNotFoundException", "Global table not found", http.StatusBadRequest)
	// ErrGlobalTableAlreadyExists is returned when a global table already exists with the same name.
	ErrGlobalTableAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#GlobalTableAlreadyExistsException", "Global table already exists", http.StatusConflict)
	// ErrReplicaAlreadyExists is returned when the specified replica already exists on the global table.
	ErrReplicaAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicaAlreadyExistsException", "Replica already exists", http.StatusConflict)
	// ErrReplicaNotFound is returned when the specified replica does not exist on the global table.
	ErrReplicaNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicaNotFoundException", "Replica not found", http.StatusBadRequest)
	// ErrTransactionConflict is returned when an operation conflicts with an ongoing transaction for the item.
	ErrTransactionConflict = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionConflictException", "TransactionConflict", http.StatusBadRequest)
	// ErrTransactionCanceled is returned when a transaction request is cancelled.
	ErrTransactionCanceled = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionCanceledException", "Transaction canceled", http.StatusBadRequest)
	// ErrIdempotentParameterMismatch is returned when a retried request has a different payload but the same idempotency token.
	ErrIdempotentParameterMismatch = NewAPIError("com.amazonaws.dynamodb.v20120810#IdempotentParameterMismatchException", "Idempotent parameter mismatch", http.StatusBadRequest)
	// ErrExportNotFound is returned when the specified export does not exist.
	ErrExportNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ExportNotFoundException", "Export not found", http.StatusBadRequest)
	// ErrImportNotFound is returned when the specified import does not exist.
	ErrImportNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ImportNotFoundException", "Import not found", http.StatusBadRequest)
	// ErrPolicyNotFound is returned when a resource-based policy is not found for the specified resource.
	ErrPolicyNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#PolicyNotFoundException", "Policy not found", http.StatusBadRequest)
	// ErrIndexNotFound is returned when the specified secondary index does not exist on the table.
	ErrIndexNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#IndexNotFoundException", "Index not found", http.StatusBadRequest)
	// ErrIndexAlreadyExists is returned when a secondary index already exists with the same name on the table.
	ErrIndexAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Index already exists", http.StatusConflict)
	// ErrPITRNotEnabled is returned when point-in-time recovery is not enabled for the table.
	ErrPITRNotEnabled = NewAPIError("com.amazonaws.dynamodb.v20120810#PointInTimeRecoveryUnavailableException", "Point in time recovery is not enabled for this table", http.StatusBadRequest)
	// ErrContinuousBackupsUnavailable is returned when continuous backups are unavailable.
	ErrContinuousBackupsUnavailable = NewAPIError("com.amazonaws.dynamodb.v20120810#ContinuousBackupsUnavailableException", "Backups are not available for this table", http.StatusBadRequest)
	// ErrDuplicateItem is returned when a batch write contains duplicate items.
	ErrDuplicateItem = NewAPIError("com.amazonaws.dynamodb.v20120810#DuplicateItemException", "Duplicate item in request", http.StatusBadRequest)
	// ErrExportConflict is returned when an export operation conflicts with an existing export.
	ErrExportConflict = NewAPIError("com.amazonaws.dynamodb.v20120810#ExportConflictException", "Export conflict", http.StatusBadRequest)
	// ErrFailure is returned when a generic failure occurs during import or export.
	ErrFailure = NewAPIError("com.amazonaws.dynamodb.v20120810#FailureException", "Failure", http.StatusInternalServerError)
	// ErrImportConflict is returned when an import operation conflicts with an existing import.
	ErrImportConflict = NewAPIError("com.amazonaws.dynamodb.v20120810#ImportConflictException", "Import conflict", http.StatusBadRequest)
	// ErrItemCollectionSizeLimitExceeded is returned when an item collection exceeds the size limit.
	ErrItemCollectionSizeLimitExceeded = NewAPIError("com.amazonaws.dynamodb.v20120810#ItemCollectionSizeLimitExceededException", "Item collection size limit exceeded", http.StatusBadRequest)
	// ErrLimitExceeded is returned when a service limit is exceeded.
	ErrLimitExceeded = NewAPIError("com.amazonaws.dynamodb.v20120810#LimitExceededException", "Limit exceeded", http.StatusBadRequest)
	// ErrProvisionedThroughputExceeded is returned when provisioned throughput is exceeded.
	ErrProvisionedThroughputExceeded = NewAPIError("com.amazonaws.dynamodb.v20120810#ProvisionedThroughputExceededException", "Provisioned throughput exceeded", http.StatusBadRequest)
	// ErrReplicatedWriteConflict is returned when a replicated write conflict occurs.
	ErrReplicatedWriteConflict = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicatedWriteConflictException", "Replicated write conflict", http.StatusBadRequest)
	// ErrThrottling is returned when the request is throttled.
	ErrThrottling = NewAPIError("com.amazonaws.dynamodb.v20120810#ThrottlingException", "Rate of requests exceeds throughput limit", http.StatusBadRequest)
	// ErrTransactionInProgress is returned when a transaction is already in progress.
	ErrTransactionInProgress = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionInProgressException", "Transaction in progress", http.StatusBadRequest)
	// ErrInvalidExportTime is returned when the requested export time is invalid.
	ErrInvalidExportTime = NewAPIError("com.amazonaws.dynamodb.v20120810#InvalidExportTimeException", "Invalid export time", http.StatusBadRequest)
	// ErrInvalidRestoreTime is returned when the requested restore time is invalid.
	ErrInvalidRestoreTime = NewAPIError("com.amazonaws.dynamodb.v20120810#InvalidRestoreTimeException", "Invalid restore time", http.StatusBadRequest)
)
