package dynamodb

import (
	"encoding/json"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
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

	b, _ := json.Marshal(resp)
	return string(b)
}

// Predefined DynamoDB error variables.
var (
	ErrResourceNotFound            = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Requested resource not found", http.StatusNotFound)
	ErrResourceAlreadyExists       = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Resource already exists", http.StatusConflict)
	ErrInvalidParameter            = NewAPIError("com.amazon.coral.validate#ValidationException", "Invalid parameter", http.StatusBadRequest)
	ErrTableNotFound               = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Requested resource not found: Table not found", http.StatusNotFound)
	ErrTableAlreadyExists          = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table already exists", http.StatusConflict)
	ErrTableNotActive              = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is not in ACTIVE state", http.StatusBadRequest)
	ErrTableDeletionProtected      = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is protected from deletion", http.StatusBadRequest)
	ErrConditionalCheckFailed      = NewAPIError("com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException", "The conditional request failed", http.StatusBadRequest)
	ErrItemNotFound                = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Item not found", http.StatusNotFound)
	ErrInvalidKey                  = NewAPIError("com.amazon.coral.validate#ValidationException", "Invalid key", http.StatusBadRequest)
	ErrMissingKey                  = NewAPIError("com.amazon.coral.validate#ValidationException", "Missing required key in request", http.StatusBadRequest)
	ErrBackupNotFound              = NewAPIError("com.amazonaws.dynamodb.v20120810#BackupNotFoundException", "Backup not found", http.StatusNotFound)
	ErrBackupAlreadyExists         = NewAPIError("com.amazonaws.dynamodb.v20120810#BackupInUseException", "Backup already exists", http.StatusConflict)
	ErrGlobalTableNotFound         = NewAPIError("com.amazonaws.dynamodb.v20120810#GlobalTableNotFoundException", "Global table not found", http.StatusNotFound)
	ErrGlobalTableAlreadyExists    = NewAPIError("com.amazonaws.dynamodb.v20120810#GlobalTableAlreadyExistsException", "Global table already exists", http.StatusConflict)
	ErrReplicaAlreadyExists        = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicaAlreadyExistsException", "Replica already exists", http.StatusConflict)
	ErrReplicaNotFound             = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicaNotFoundException", "Replica not found", http.StatusNotFound)
	ErrTransactionConflict         = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionConflictException", "TransactionConflict", http.StatusBadRequest)
	ErrTransactionCanceled         = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionCanceledException", "Transaction canceled", http.StatusBadRequest)
	ErrIdempotentParameterMismatch = NewAPIError("com.amazonaws.dynamodb.v20120810#IdempotentParameterMismatchException", "Idempotent parameter mismatch", http.StatusBadRequest)
	ErrExportNotFound              = NewAPIError("com.amazonaws.dynamodb.v20120810#ExportNotFoundException", "Export not found", http.StatusNotFound)
	ErrImportNotFound              = NewAPIError("com.amazonaws.dynamodb.v20120810#ImportNotFoundException", "Import not found", http.StatusNotFound)
	ErrPolicyNotFound              = NewAPIError("com.amazonaws.dynamodb.v20120810#PolicyNotFoundException", "Policy not found", http.StatusNotFound)
	ErrIndexNotFound               = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Index not found", http.StatusNotFound)
	ErrIndexAlreadyExists          = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Index already exists", http.StatusConflict)
	ErrPITRNotEnabled              = NewAPIError("com.amazonaws.dynamodb.v20120810#PointInTimeRecoveryUnavailableException", "Point in time recovery is not enabled for this table", http.StatusBadRequest)
)
