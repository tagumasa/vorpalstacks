package sqs

import (
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

// TestSQSErrorWireContract pins every service error's JSON error code, HTTP
// status, and query-protocol error code to the values the SQS service model
// defines (the shape name, the httpError/awsQueryError httpResponseCode, and
// the awsQueryError code respectively). Errors without an awsQueryError trait
// in the model have no query code and fall back to the shape name on the
// query wire. A new error that misses this table fails the test, as does an
// edit that diverges from the model.
func TestSQSErrorWireContract(t *testing.T) {
	cases := []struct {
		name       string
		err        *awserrors.AWSError
		code       string
		httpStatus int
		queryCode  string
	}{
		{"QueueDoesNotExist", ErrQueueDoesNotExist, "QueueDoesNotExist", 400, "AWS.SimpleQueueService.NonExistentQueue"},
		{"QueueDeletedRecently", ErrQueueDeletedRecently, "QueueDeletedRecently", 400, "AWS.SimpleQueueService.QueueDeletedRecently"},
		{"QueueNameExists", ErrQueueNameExists, "QueueNameExists", 400, "QueueAlreadyExists"},
		{"InvalidQueueName", ErrInvalidQueueName, "InvalidParameterValue", 400, ""},
		{"InvalidParameterValue", ErrInvalidParameterValue, "InvalidParameterValue", 400, ""},
		{"MissingParameter", ErrMissingParameter, "MissingParameter", 400, ""},
		{"InvalidBatchEntryId", ErrInvalidBatchEntryId, "InvalidBatchEntryId", 400, "AWS.SimpleQueueService.InvalidBatchEntryId"},
		{"BatchEntryIdsNotDistinct", ErrBatchEntryIdsNotDistinct, "BatchEntryIdsNotDistinct", 400, "AWS.SimpleQueueService.BatchEntryIdsNotDistinct"},
		{"TooManyEntriesInBatch", ErrTooManyEntriesInBatch, "TooManyEntriesInBatchRequest", 400, "AWS.SimpleQueueService.TooManyEntriesInBatchRequest"},
		{"EmptyBatchRequest", ErrEmptyBatchRequest, "EmptyBatchRequest", 400, "AWS.SimpleQueueService.EmptyBatchRequest"},
		{"ReceiptHandleIsInvalid", ErrReceiptHandleIsInvalid, "ReceiptHandleIsInvalid", 404, "ReceiptHandleIsInvalid"},
		{"InvalidAttributeName", ErrInvalidAttributeName, "InvalidAttributeName", 400, ""},
		{"MessageNotInflight", ErrMessageNotInflight, "MessageNotInflight", 400, "AWS.SimpleQueueService.MessageNotInflight"},
		{"PurgeQueueInProgress", ErrPurgeQueueInProgress, "PurgeQueueInProgress", 403, "AWS.SimpleQueueService.PurgeQueueInProgress"},
		{"MessageTooLarge", ErrMessageTooLarge, "InvalidMessageContents", 400, ""},
		{"InvalidMessageContents", ErrInvalidMessageContents, "InvalidMessageContents", 400, ""},
		{"TooManyTags", ErrTooManyTags, "InvalidParameterValue", 400, ""},
		{"InvalidTagKey", ErrInvalidTagKey, "InvalidParameterValue", 400, ""},
		{"InvalidTagValue", ErrInvalidTagValue, "InvalidParameterValue", 400, ""},
		{"MissingMessageGroupId", ErrMissingMessageGroupId, "InvalidParameterValue", 400, ""},
		{"MissingDeduplicationId", ErrMissingDeduplicationId, "InvalidParameterValue", 400, ""},
		{"InvalidParameterCombination", ErrInvalidParameterCombination, "InvalidParameterCombination", 400, ""},
		{"ResourceNotFound", ErrResourceNotFound, "ResourceNotFoundException", 404, "ResourceNotFoundException"},
		{"OverLimit", ErrOverLimit, "OverLimit", 403, "OverLimit"},
		{"InvalidAttributeValue", ErrInvalidAttributeValue, "InvalidAttributeValue", 400, ""},
		{"BatchRequestTooLong", ErrBatchRequestTooLong, "BatchRequestTooLong", 400, "AWS.SimpleQueueService.BatchRequestTooLong"},
		{"KmsAccessDenied", ErrKmsAccessDenied, "KmsAccessDenied", 400, "KMS.AccessDeniedException"},
		{"KmsDisabled", ErrKmsDisabled, "KmsDisabled", 400, "KMS.DisabledException"},
		{"KmsInvalidKeyUsage", ErrKmsInvalidKeyUsage, "KmsInvalidKeyUsage", 400, "KMS.InvalidKeyUsageException"},
		{"KmsInvalidState", ErrKmsInvalidState, "KmsInvalidState", 400, "KMS.InvalidStateException"},
		{"KmsNotFound", ErrKmsNotFound, "KmsNotFound", 400, "KMS.NotFoundException"},
		{"KmsOptInRequired", ErrKmsOptInRequired, "KmsOptInRequired", 403, "KMS.OptInRequired"},
		{"KmsThrottled", ErrKmsThrottled, "KmsThrottled", 400, "KMS.ThrottlingException"},
		{"InvalidAddress", ErrInvalidAddress, "InvalidAddress", 404, "InvalidAddress"},
		{"InvalidSecurity", ErrInvalidSecurity, "InvalidSecurity", 403, "InvalidSecurity"},
		{"RequestThrottled", ErrRequestThrottled, "RequestThrottled", 403, "RequestThrottled"},
		{"UnsupportedOperation", ErrUnsupportedOperation, "UnsupportedOperation", 400, "AWS.SimpleQueueService.UnsupportedOperation"},
		{"InvalidIdFormat", ErrInvalidIdFormat, "InvalidIdFormat", 400, ""},
	}
	for _, tc := range cases {
		if tc.err.GetCode() != tc.code {
			t.Errorf("%s: code = %q, want %q", tc.name, tc.err.GetCode(), tc.code)
		}
		if tc.err.GetHTTPStatusCode() != tc.httpStatus {
			t.Errorf("%s: HTTP status = %d, want %d", tc.name, tc.err.GetHTTPStatusCode(), tc.httpStatus)
		}
		if tc.err.QueryErrorCode != tc.queryCode {
			t.Errorf("%s: query code = %q, want %q", tc.name, tc.err.QueryErrorCode, tc.queryCode)
		}
	}
}
