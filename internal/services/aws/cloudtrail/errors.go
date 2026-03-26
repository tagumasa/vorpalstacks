package cloudtrail

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// CloudTrailError represents an error returned by the CloudTrail service.
type CloudTrailError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWS error.
func (e *CloudTrailError) Unwrap() error {
	return e.AWSError
}

var (
	// ErrTrailNotFound is returned when the specified trail does not exist.
	ErrTrailNotFound = &CloudTrailError{awserrors.NewAWSError("TrailNotFoundException", "Trail not found.", http.StatusNotFound)}
	// ErrTrailAlreadyExists is returned when attempting to create a trail that already exists.
	ErrTrailAlreadyExists = &CloudTrailError{awserrors.NewAWSError("TrailAlreadyExistsException", "Trail already exists.", http.StatusConflict)}
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = &CloudTrailError{awserrors.NewAWSError("InvalidParameterException", "Invalid parameter.", http.StatusBadRequest)}
	// ErrS3BucketNotFound is returned when the specified S3 bucket does not exist.
	ErrS3BucketNotFound = &CloudTrailError{awserrors.NewAWSError("S3BucketDoesNotExistException", "S3 bucket does not exist.", http.StatusNotFound)}
	// ErrInsufficientSnsPolicy is returned when the SNS topic policy is insufficient.
	ErrInsufficientSnsPolicy = &CloudTrailError{awserrors.NewAWSError("InsufficientSnsTopicPolicyException", "Insufficient SNS topic policy.", http.StatusBadRequest)}
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = &CloudTrailError{awserrors.NewAWSError("AccessDeniedException", "Access denied.", http.StatusForbidden)}
	// ErrInternalError is returned when an internal server error occurs.
	ErrInternalError = &CloudTrailError{awserrors.NewAWSError("InternalError", "Internal server error.", http.StatusInternalServerError)}
)
