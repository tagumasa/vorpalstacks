package cloudtrail

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
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
	// InternalFailure is the only 500 code the CloudTrail contract
	// documents (Common Error Types).
	ErrInternalError = awserrors.NewAWSError("InternalFailure", "The request can't be processed right now because of an internal server issue.", http.StatusInternalServerError)
	// ErrConflictException is returned when an operation conflicts with the current state.
	ErrConflictException = awserrors.NewAWSError("ConflictException", "Operation conflicts with current state.", http.StatusConflict)
	// ErrEventDataStoreNotFoundException is returned when the specified event data store does not exist.
	ErrEventDataStoreNotFoundException = awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found.", http.StatusNotFound)
	// ErrOperationNotPermitted is returned when the operation is not permitted in the current state.
	ErrOperationNotPermitted = awserrors.NewAWSError("OperationNotPermittedException", "Operation not permitted.", http.StatusBadRequest)
)

// storeErrorMappings maps store-level sentinel errors to CloudTrail API errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: cloudtrailstore.ErrTrailNotFound, AWS: ErrTrailNotFound},
	{Store: cloudtrailstore.ErrTrailAlreadyExists, AWS: ErrTrailAlreadyExists},
}
