package cloudtrail

import (
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

var (
	// ErrTrailNotFound is returned when the specified trail does not exist.
	ErrTrailNotFound = awserrors.NewAWSError("TrailNotFoundException", "Trail not found.", http.StatusNotFound)
	// ErrTrailAlreadyExists is returned when attempting to create a trail that already exists.
	ErrTrailAlreadyExists = awserrors.NewAWSError("TrailAlreadyExistsException", "Trail already exists.", http.StatusConflict)
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = awserrors.NewAWSError("InvalidParameterException", "Invalid parameter.", http.StatusBadRequest)
	// ErrS3BucketNotFound is returned when the specified S3 bucket does not exist.
	ErrS3BucketNotFound = awserrors.NewAWSError("S3BucketDoesNotExistException", "S3 bucket does not exist.", http.StatusNotFound)
	// ErrInsufficientSnsPolicy is returned when the SNS topic policy is insufficient.
	ErrInsufficientSnsPolicy = awserrors.NewAWSError("InsufficientSnsTopicPolicyException", "Insufficient SNS topic policy.", http.StatusBadRequest)
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = awserrors.NewAWSError("AccessDeniedException", "Access denied.", http.StatusForbidden)
	// ErrInternalError is returned when an internal server error occurs.
	ErrInternalError = awserrors.NewAWSError("InternalError", "Internal server error.", http.StatusInternalServerError)
)
