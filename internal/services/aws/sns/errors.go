package sns

import (
	"errors"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

var (
	// ErrAuthorizationError is returned when the request is not authorized.
	ErrAuthorizationError = awserrors.NewAWSError("AuthorizationError", "Authorization error", 403)
	// ErrFilterLimitExceeded is returned when the filter limit is exceeded.
	ErrFilterLimitExceeded = awserrors.NewAWSError("FilterPolicyLimitExceeded", "Filter policy limit exceeded", 403)
	// ErrTopicNotFound is returned when the topic does not exist.
	ErrTopicNotFound = awserrors.NewAWSError("NotFound", "Topic does not exist", 404)
	// ErrSubscriptionNotFound is returned when the subscription does not exist.
	ErrSubscriptionNotFound = awserrors.NewAWSError("NotFound", "Subscription does not exist", 404)
	// ErrPlatformAppNotFound is returned when the platform application does not exist.
	ErrPlatformAppNotFound = awserrors.NewAWSError("NotFound", "Platform application does not exist", 404)
	// ErrEndpointNotFound is returned when the endpoint does not exist.
	ErrEndpointNotFound = awserrors.NewAWSError("NotFound", "Endpoint does not exist", 404)
	// ErrResourceNotFound is returned by the tag family when the resource
	// ARN addresses no existing topic. The tag operations declare the
	// model's ResourceNotFoundException shape, whose awsQueryError wire
	// code is "ResourceNotFound" — a different code from the
	// NotFoundException "NotFound" the rest of the API reports.
	ErrResourceNotFound = awserrors.NewAWSError("ResourceNotFound", "Topic does not exist", 404)
	// ErrTagLimitExceeded is returned when the tag limit is exceeded. The
	// message is the model's own TagLimitExceededException documentation.
	ErrTagLimitExceeded = awserrors.NewAWSError("TagLimitExceeded", "Can't add more than 50 tags to a topic.", 400)
	// ErrTopicLimitExceeded is returned when the account already owns the
	// maximum allowed number of topics for the type being created. The
	// code and message are the model's TopicLimitExceededException
	// (awsQueryError "TopicLimitExceeded", HTTP 403).
	ErrTopicLimitExceeded = awserrors.NewAWSError("TopicLimitExceeded", "Indicates that the customer already owns the maximum allowed number of topics.", 403)
	// ErrSubscriptionLimitExceeded is returned when the topic already
	// carries the maximum allowed number of subscriptions. The code and
	// message are the model's SubscriptionLimitExceededException
	// (awsQueryError "SubscriptionLimitExceeded", HTTP 403).
	ErrSubscriptionLimitExceeded = awserrors.NewAWSError("SubscriptionLimitExceeded", "Indicates that the customer already owns the maximum allowed number of subscriptions.", 403)
	// ErrBatchEntryIdsNotDistinct is returned when two or more batch entries have the same ID.
	ErrBatchEntryIdsNotDistinct = awserrors.NewAWSError("BatchEntryIdsNotDistinct", "Two or more batch entries have the same ID", 400)
	// ErrEmptyBatchRequest is returned when a batch request carries no entries.
	ErrEmptyBatchRequest = awserrors.NewAWSError("EmptyBatchRequest", "Batch request does not contain any entries", 400)
	// ErrTooManyEntriesInBatch is returned when the batch request exceeds the
	// entry limit.
	ErrTooManyEntriesInBatch = awserrors.NewAWSError("TooManyEntriesInBatchRequest", fmt.Sprintf("Maximum number of entries per request are %d", snsstore.MaxBatchEntries), 400)
	// ErrBatchRequestTooLong is returned when the batch request's total size
	// exceeds the platform's batch size ceiling.
	ErrBatchRequestTooLong = awserrors.NewAWSError("BatchRequestTooLong", fmt.Sprintf("Total batch request size exceeds the maximum of %d", snsstore.MaxBatchTotalSize), 400)
)

// mapStoreError is the single store→service error-translation seam: every
// Core funnels its store errors through this table, so the wire identity
// each store sentinel maps to is decided once, here, against the model's
// awsQueryError traits — not per call site. The not-found family reports
// the NotFoundException code "NotFound"; the platform already-exists
// refusal and the token/label constraint violations report
// "InvalidParameter" (the InvalidParameterException trait code); errors
// the table does not know (storage failures) pass through unchanged.
func mapStoreError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, snsstore.ErrTopicNotFound):
		return ErrTopicNotFound
	case errors.Is(err, snsstore.ErrSubscriptionNotFound):
		return ErrSubscriptionNotFound
	case errors.Is(err, snsstore.ErrPlatformApplicationNotFound):
		return ErrPlatformAppNotFound
	case errors.Is(err, snsstore.ErrEndpointNotFound):
		return ErrEndpointNotFound
	case errors.Is(err, snsstore.ErrPlatformApplicationAlreadyExists):
		return NewInvalidParameter("Platform application already exists with the same name")
	case errors.Is(err, snsstore.ErrTopicLimitExceeded):
		return ErrTopicLimitExceeded
	case errors.Is(err, snsstore.ErrSubscriptionLimitExceeded):
		return ErrSubscriptionLimitExceeded
	case errors.Is(err, snsstore.ErrFilterPolicyLimitExceeded):
		return ErrFilterLimitExceeded
	case errors.Is(err, snsstore.ErrPermissionLabelExists):
		return NewInvalidParameter("Invalid parameter: Label already exists on the topic's policy")
	case errors.Is(err, snsstore.ErrInvalidToken):
		return NewInvalidParameter("Invalid parameter: Token")
	default:
		return err
	}
}

// NewInvalidParameter creates an InvalidParameterException error carrying
// the wire code the SNS model assigns to that shape: its awsQueryError
// trait maps the shape to "InvalidParameter", and AWS query responses must
// report the trait code rather than the shape name.
func NewInvalidParameter(message string) *awserrors.AWSError {
	return awserrors.NewInvalidParameterException(message).SetQueryErrorCode("InvalidParameter")
}
