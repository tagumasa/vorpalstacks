package sqs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQSErrors(t *testing.T) {
	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "QueueDoesNotExist: The specified queue does not exist.", ErrQueueDoesNotExist.Error())
		assert.Equal(t, 400, ErrQueueDoesNotExist.GetHTTPStatusCode())

		assert.Equal(t, "QueueNameExists: A queue already exists with the same name and a different configuration.", ErrQueueNameExists.Error())
		assert.Equal(t, 400, ErrQueueNameExists.GetHTTPStatusCode())

		assert.Equal(t, "InvalidParameterValue: Invalid queue name.", ErrInvalidQueueName.Error())
		assert.Equal(t, 400, ErrInvalidQueueName.GetHTTPStatusCode())

		assert.Equal(t, "InvalidParameterValue: Invalid parameter value.", ErrInvalidParameterValue.Error())
		assert.Equal(t, 400, ErrInvalidParameterValue.GetHTTPStatusCode())

		assert.Equal(t, "MissingParameter: The request must contain a required parameter.", ErrMissingParameter.Error())
		assert.Equal(t, 400, ErrMissingParameter.GetHTTPStatusCode())

		assert.Equal(t, "InvalidBatchEntryId: The batch entry ID is invalid.", ErrInvalidBatchEntryId.Error())
		assert.Equal(t, 400, ErrInvalidBatchEntryId.GetHTTPStatusCode())

		assert.Equal(t, "BatchEntryIdsNotDistinct: Two or more batch entries have the same ID.", ErrBatchEntryIdsNotDistinct.Error())
		assert.Equal(t, 400, ErrBatchEntryIdsNotDistinct.GetHTTPStatusCode())

		assert.Equal(t, "TooManyEntriesInBatchRequest: Maximum number of entries per request are 10.", ErrTooManyEntriesInBatch.Error())
		assert.Equal(t, 400, ErrTooManyEntriesInBatch.GetHTTPStatusCode())

		assert.Equal(t, "EmptyBatchRequest: The batch request doesn't contain any entries.", ErrEmptyBatchRequest.Error())
		assert.Equal(t, 400, ErrEmptyBatchRequest.GetHTTPStatusCode())

		assert.Equal(t, "ReceiptHandleIsInvalid: The receipt handle provided is not valid.", ErrReceiptHandleIsInvalid.Error())
		assert.Equal(t, 400, ErrReceiptHandleIsInvalid.GetHTTPStatusCode())

		assert.Equal(t, "InvalidAttributeName: The attribute name is invalid.", ErrInvalidAttributeName.Error())
		assert.Equal(t, 400, ErrInvalidAttributeName.GetHTTPStatusCode())

		assert.Equal(t, "MessageNotInflight: The message referred to is not in flight.", ErrMessageNotInflight.Error())
		assert.Equal(t, 400, ErrMessageNotInflight.GetHTTPStatusCode())

		assert.Equal(t, "PurgeQueueInProgress: There is already a purge queue operation running.", ErrPurgeQueueInProgress.Error())
		assert.Equal(t, 403, ErrPurgeQueueInProgress.GetHTTPStatusCode())
	})
}
