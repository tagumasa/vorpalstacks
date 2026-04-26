package dynamodb

import (
	"encoding/json"
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
	ErrResourceNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Requested resource not found", http.StatusNotFound)
	// ErrResourceAlreadyExists is returned when a resource already exists with the same identifier.
	ErrResourceAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Resource already exists", http.StatusConflict)
	// ErrInvalidParameter is returned when a request parameter fails validation.
	ErrInvalidParameter = NewAPIError("com.amazon.coral.validate#ValidationException", "Invalid parameter", http.StatusBadRequest)
	// ErrTableNotFound is returned when the specified table does not exist.
	ErrTableNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Requested resource not found: Table not found", http.StatusNotFound)
	// ErrTableAlreadyExists is returned when a table already exists with the same name.
	ErrTableAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table already exists", http.StatusConflict)
	// ErrTableNotActive is returned when the table is not in ACTIVE state.
	ErrTableNotActive = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is not in ACTIVE state", http.StatusBadRequest)
	// ErrTableDeletionProtected is returned when deletion protection is enabled on the table.
	ErrTableDeletionProtected = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Table is protected from deletion", http.StatusBadRequest)
	// ErrConditionalCheckFailed is returned when a conditional write operation fails.
	ErrConditionalCheckFailed = NewAPIError("com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException", "The conditional request failed", http.StatusBadRequest)
	// ErrItemNotFound is returned when the specified item does not exist in the table.
	ErrItemNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Item not found", http.StatusNotFound)
	// ErrInvalidKey is returned when the provided key is invalid.
	ErrInvalidKey = NewAPIError("com.amazon.coral.validate#ValidationException", "Invalid key", http.StatusBadRequest)
	// ErrMissingKey is returned when a required key is missing from the request.
	ErrMissingKey = NewAPIError("com.amazon.coral.validate#ValidationException", "Missing required key in request", http.StatusBadRequest)
	// ErrBackupNotFound is returned when the specified backup does not exist.
	ErrBackupNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#BackupNotFoundException", "Backup not found", http.StatusNotFound)
	// ErrBackupAlreadyExists is returned when a backup already exists with the same name.
	ErrBackupAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#BackupInUseException", "Backup already exists", http.StatusConflict)
	// ErrGlobalTableNotFound is returned when the specified global table does not exist.
	ErrGlobalTableNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#GlobalTableNotFoundException", "Global table not found", http.StatusNotFound)
	// ErrGlobalTableAlreadyExists is returned when a global table already exists with the same name.
	ErrGlobalTableAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#GlobalTableAlreadyExistsException", "Global table already exists", http.StatusConflict)
	// ErrReplicaAlreadyExists is returned when the specified replica already exists on the global table.
	ErrReplicaAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicaAlreadyExistsException", "Replica already exists", http.StatusConflict)
	// ErrReplicaNotFound is returned when the specified replica does not exist on the global table.
	ErrReplicaNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ReplicaNotFoundException", "Replica not found", http.StatusNotFound)
	// ErrTransactionConflict is returned when an operation conflicts with an ongoing transaction for the item.
	ErrTransactionConflict = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionConflictException", "TransactionConflict", http.StatusBadRequest)
	// ErrTransactionCanceled is returned when a transaction request is cancelled.
	ErrTransactionCanceled = NewAPIError("com.amazonaws.dynamodb.v20120810#TransactionCanceledException", "Transaction canceled", http.StatusBadRequest)
	// ErrIdempotentParameterMismatch is returned when a retried request has a different payload but the same idempotency token.
	ErrIdempotentParameterMismatch = NewAPIError("com.amazonaws.dynamodb.v20120810#IdempotentParameterMismatchException", "Idempotent parameter mismatch", http.StatusBadRequest)
	// ErrExportNotFound is returned when the specified export does not exist.
	ErrExportNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ExportNotFoundException", "Export not found", http.StatusNotFound)
	// ErrImportNotFound is returned when the specified import does not exist.
	ErrImportNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ImportNotFoundException", "Import not found", http.StatusNotFound)
	// ErrPolicyNotFound is returned when a resource-based policy is not found for the specified resource.
	ErrPolicyNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#PolicyNotFoundException", "Policy not found", http.StatusNotFound)
	// ErrIndexNotFound is returned when the specified secondary index does not exist on the table.
	ErrIndexNotFound = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceNotFoundException", "Index not found", http.StatusNotFound)
	// ErrIndexAlreadyExists is returned when a secondary index already exists with the same name on the table.
	ErrIndexAlreadyExists = NewAPIError("com.amazonaws.dynamodb.v20120810#ResourceInUseException", "Index already exists", http.StatusConflict)
	// ErrPITRNotEnabled is returned when point-in-time recovery is not enabled for the table.
	ErrPITRNotEnabled = NewAPIError("com.amazonaws.dynamodb.v20120810#PointInTimeRecoveryUnavailableException", "Point in time recovery is not enabled for this table", http.StatusBadRequest)
)
