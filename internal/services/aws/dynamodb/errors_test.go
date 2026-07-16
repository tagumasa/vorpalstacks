package dynamodb

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	awserrors "vorpalstacks/internal/common/errors"
)

func TestDynamoDBErrors(t *testing.T) {
	t.Run("NewAPIError", func(t *testing.T) {
		err := NewAPIError("TestCode", "Test message", 400)
		assert.Equal(t, "TestCode: Test message", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ResourceNotFoundException: Requested resource not found", ErrResourceNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrResourceNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ResourceInUseException: Resource already exists", ErrResourceAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrResourceAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "com.amazon.coral.validate#ValidationException: Invalid parameter", ErrInvalidParameter.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidParameter.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ResourceNotFoundException: Requested resource not found: Table not found", ErrTableNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrTableNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ResourceInUseException: Table already exists", ErrTableAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrTableAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ResourceInUseException: Table is not in ACTIVE state", ErrTableNotActive.Error())
		assert.Equal(t, http.StatusBadRequest, ErrTableNotActive.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException: The conditional request failed", ErrConditionalCheckFailed.Error())
		assert.Equal(t, http.StatusBadRequest, ErrConditionalCheckFailed.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ResourceNotFoundException: Item not found", ErrItemNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrItemNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazon.coral.validate#ValidationException: Invalid key", ErrInvalidKey.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidKey.GetHTTPStatusCode())

		assert.Equal(t, "com.amazon.coral.validate#ValidationException: Missing required key in request", ErrMissingKey.Error())
		assert.Equal(t, http.StatusBadRequest, ErrMissingKey.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#BackupNotFoundException: Backup not found", ErrBackupNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrBackupNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#BackupInUseException: Backup already exists", ErrBackupAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrBackupAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#GlobalTableNotFoundException: Global table not found", ErrGlobalTableNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrGlobalTableNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#GlobalTableAlreadyExistsException: Global table already exists", ErrGlobalTableAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrGlobalTableAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ReplicaAlreadyExistsException: Replica already exists", ErrReplicaAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrReplicaAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ReplicaNotFoundException: Replica not found", ErrReplicaNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrReplicaNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#TransactionConflictException: TransactionConflict", ErrTransactionConflict.Error())
		assert.Equal(t, http.StatusBadRequest, ErrTransactionConflict.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#IdempotentParameterMismatchException: Idempotent parameter mismatch", ErrIdempotentParameterMismatch.Error())
		assert.Equal(t, http.StatusBadRequest, ErrIdempotentParameterMismatch.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ExportNotFoundException: Export not found", ErrExportNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrExportNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#ImportNotFoundException: Import not found", ErrImportNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrImportNotFound.GetHTTPStatusCode())

		assert.Equal(t, "com.amazonaws.dynamodb.v20120810#PolicyNotFoundException: Policy not found", ErrPolicyNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrPolicyNotFound.GetHTTPStatusCode())
	})
}

func TestTransactionCanceledError_ToJSON(t *testing.T) {
	t.Run("normal serialisation", func(t *testing.T) {
		err := &TransactionCanceledError{
			APIError: &APIError{awserrors.NewAWSError("TransactionCanceledException", "Transaction cancelled", http.StatusBadRequest)},
			CancellationReasons: []CancellationReason{
				{Code: "ConditionalCheckFailed"},
				{Code: "None"},
			},
		}
		jsonStr := err.ToJSON()

		var parsed map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(jsonStr), &parsed))
		assert.Equal(t, "TransactionCanceledException", parsed["__type"])
		assert.Equal(t, "Transaction cancelled", parsed["message"])

		reasons, ok := parsed["CancellationReasons"].([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 2, len(reasons))
	})

	t.Run("with item data", func(t *testing.T) {
		err := &TransactionCanceledError{
			APIError: &APIError{awserrors.NewAWSError("TransactionCanceledException", "msg", http.StatusBadRequest)},
			CancellationReasons: []CancellationReason{
				{Code: "ConditionalCheckFailed", Item: map[string]interface{}{"id": "test"}},
			},
		}
		jsonStr := err.ToJSON()

		var parsed map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(jsonStr), &parsed))
		reasons := parsed["CancellationReasons"].([]interface{})
		reason := reasons[0].(map[string]interface{})
		assert.Equal(t, "ConditionalCheckFailed", reason["Code"])
		assert.NotNil(t, reason["Item"])
	})

	t.Run("fallback on marshal failure (Bug D-7)", func(t *testing.T) {
		err := &TransactionCanceledError{
			APIError: &APIError{awserrors.NewAWSError("TestCode", "TestMsg", 400)},
		}
		jsonStr := err.ToJSON()
		assert.Contains(t, jsonStr, "TestCode")
		assert.Contains(t, jsonStr, "TestMsg")
	})
}
