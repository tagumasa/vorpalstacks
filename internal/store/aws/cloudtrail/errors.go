package cloudtrail

import "errors"

var (
	// ErrTrailNotFound is returned when the specified CloudTrail trail
	// does not exist.
	ErrTrailNotFound = errors.New("trail not found")

	// ErrTrailAlreadyExists is returned when attempting to create a trail
	// that already exists.
	ErrTrailAlreadyExists = errors.New("trail already exists")

	// ErrS3BucketNotFound is returned when the specified S3 bucket
	// does not exist.
	ErrS3BucketNotFound = errors.New("s3 bucket not found")

	// ErrInvalidTrailName is returned when the trail name is not valid.
	ErrInvalidTrailName = errors.New("invalid trail name")

	// ErrInsufficientSnsTopicPolicy is returned when the SNS topic policy
	// does not have the required permissions.
	ErrInsufficientSnsTopicPolicy = errors.New("insufficient sns topic policy")

	// ErrEventNotFound is returned when the specified CloudTrail event
	// does not exist.
	ErrEventNotFound = errors.New("event not found")

	// ErrEventDataStoreNotFound is returned when the specified event data
	// store does not exist.
	ErrEventDataStoreNotFound = errors.New("event data store not found")

	// ErrEventDataStoreAlreadyExists is returned when attempting to create
	// an event data store with a name that is already in use.
	ErrEventDataStoreAlreadyExists = errors.New("event data store already exists")

	// ErrInvalidEventDataStoreState is returned when an operation is
	// attempted on an event data store in an invalid state.
	ErrInvalidEventDataStoreState = errors.New("invalid event data store state")

	// ErrEventDataStoreNotPendingDeletion is returned when attempting to
	// restore an event data store that is not in PENDING_DELETION state.
	ErrEventDataStoreNotPendingDeletion = errors.New("event data store is not pending deletion")

	// ErrQueryNotFound is returned when the specified query does not exist.
	ErrQueryNotFound = errors.New("query not found")

	// ErrChannelNotFound is returned when the specified channel does not exist.
	ErrChannelNotFound = errors.New("channel not found")

	// ErrChannelAlreadyExists is returned when a channel name already exists.
	ErrChannelAlreadyExists = errors.New("channel already exists")
)
