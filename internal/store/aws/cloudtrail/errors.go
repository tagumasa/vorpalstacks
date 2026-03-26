package cloudtrail

import "errors"

var (
	// ErrTrailNotFound is returned when the specified CloudTrail trail
	// does not exist.
	ErrTrailNotFound = errors.New("trail not found")

	// ErrTrailAlreadyExists is returned when attempting to create a trail
	// that already exists.
	ErrTrailAlreadyExists = errors.New("trail already exists")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = errors.New("invalid parameter")

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
)
