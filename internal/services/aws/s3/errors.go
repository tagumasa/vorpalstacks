// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// Predefined S3 error variables.
var (
	// ErrNoSuchBucket is returned when the specified bucket does not exist.
	ErrNoSuchBucket = awserrors.NewAWSError("NoSuchBucket", "The specified bucket does not exist", http.StatusNotFound)
	// ErrBucketAlreadyExists is returned when the requested bucket name is not available.
	ErrBucketAlreadyExists = awserrors.NewAWSError("BucketAlreadyExists", "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.", http.StatusConflict)
	// ErrBucketNotEmpty is returned when attempting to delete a bucket that is not empty.
	ErrBucketNotEmpty = awserrors.NewAWSError("BucketNotEmpty", "The bucket you tried to delete is not empty", http.StatusConflict)
	// ErrNoSuchKey is returned when the specified key does not exist in the bucket.
	ErrNoSuchKey = awserrors.NewAWSError("NoSuchKey", "The specified key does not exist.", http.StatusNotFound)
	// ErrInvalidBucketName is returned when the specified bucket name is not valid.
	ErrInvalidBucketName = awserrors.NewAWSError("InvalidBucketName", "The specified bucket is not valid.", http.StatusBadRequest)
	// ErrInvalidRequest is returned when the request is malformed or invalid.
	ErrInvalidRequest = awserrors.NewAWSError("InvalidRequest", "Invalid Request", http.StatusBadRequest)
	// ErrMalformedXML is returned when the provided XML is not well-formed.
	ErrMalformedXML = awserrors.NewAWSError("MalformedXML", "The XML you provided was not well-formed.", http.StatusBadRequest)
	// ErrMissingContentLength is returned when the Content-Length header is missing.
	ErrMissingContentLength = awserrors.NewAWSError("MissingContentLength", "You must provide the Content-Length HTTP header.", http.StatusLengthRequired)
	// ErrAccessDenied is returned when the requester does not have access to the resource.
	ErrAccessDenied = awserrors.NewAWSError("AccessDenied", "Access Denied", http.StatusForbidden)
	// ErrInvalidCopySource is returned when the copy source is invalid.
	ErrInvalidCopySource = awserrors.NewAWSError("InvalidCopySource", "The copy source is invalid", http.StatusBadRequest)
	// ErrPreconditionFailed is returned when at least one of the pre-conditions did not hold.
	ErrPreconditionFailed = awserrors.NewAWSError("PreconditionFailed", "At least one of the pre-conditions you specified did not hold", http.StatusPreconditionFailed)
	// ErrNotModified is returned when the resource has not been modified.
	ErrNotModified = awserrors.NewAWSError("NotModified", "The resource was not modified", http.StatusNotModified)
	// ErrInvalidRange is returned when the requested range cannot be satisfied.
	ErrInvalidRange = awserrors.NewAWSError("InvalidRange", "The requested range cannot be satisfied", http.StatusRequestedRangeNotSatisfiable)
	// ErrInvalidSSECustomerKey is returned when the provided SSE customer key is not valid.
	ErrInvalidSSECustomerKey = awserrors.NewAWSError("InvalidSSECustomerKey", "The provided SSE customer key is not valid", http.StatusBadRequest)
	// ErrNoSuchUpload is returned when the specified multipart upload does not exist.
	ErrNoSuchUpload = awserrors.NewAWSError("NoSuchUpload", "The specified multipart upload does not exist.", http.StatusNotFound)
	// ErrInvalidPart is returned when one or more specified parts could not be found.
	ErrInvalidPart = awserrors.NewAWSError("InvalidPart", "One or more of the specified parts could not be found.", http.StatusBadRequest)
	// ErrEntityTooLarge is returned when the proposed upload exceeds the maximum allowed size.
	ErrEntityTooLarge = awserrors.NewAWSError("EntityTooLarge", "Your proposed upload exceeds the maximum allowed size", http.StatusBadRequest)
	// ErrInvalidArgument is returned when an invalid argument is provided.
	ErrInvalidArgument = awserrors.NewAWSError("InvalidArgument", "Invalid argument", http.StatusBadRequest)
	// ErrNoSuchBucketPolicy is returned when the bucket policy does not exist.
	ErrNoSuchBucketPolicy = awserrors.NewAWSError("NoSuchBucketPolicy", "The bucket policy does not exist", http.StatusNotFound)
	// ErrNoSuchCORS is returned when the CORS configuration does not exist.
	ErrNoSuchCORS = awserrors.NewAWSError("NoSuchCORSConfiguration", "The CORS configuration does not exist", http.StatusNotFound)
	// ErrNoSuchEncryption is returned when the server-side encryption configuration is not found.
	ErrNoSuchEncryption = awserrors.NewAWSError("ServerSideEncryptionConfigurationNotFoundError", "The server side encryption configuration was not found", http.StatusNotFound)
	// ErrNoSuchTagging is returned when the tag set does not exist.
	ErrNoSuchTagging = awserrors.NewAWSError("NoSuchTagSet", "The TagSet does not exist", http.StatusNotFound)
	// ErrNoSuchVersioning is returned when the bucket does not have a versioning configuration.
	ErrNoSuchVersioning = awserrors.NewAWSError("NoSuchBucketVersioningConfiguration", "The bucket does not have a versioning configuration", http.StatusNotFound)
	// ErrNoSuchWebsite is returned when the bucket does not have a website configuration.
	ErrNoSuchWebsite = awserrors.NewAWSError("NoSuchWebsiteConfiguration", "The specified bucket does not have a website configuration", http.StatusNotFound)
	// ErrNoSuchLifecycle is returned when the lifecycle configuration does not exist.
	ErrNoSuchLifecycle = awserrors.NewAWSError("NoSuchLifecycleConfiguration", "The lifecycle configuration does not exist", http.StatusNotFound)
	// ErrNoSuchOwnershipCtrls is returned when the bucket does not have ownership controls.
	ErrNoSuchOwnershipCtrls = awserrors.NewAWSError("OwnershipControlsNotFoundError", "The bucket does not have OwnershipControls", http.StatusNotFound)
	// ErrNoSuchPublicAccessBlk is returned when the public access block configuration is not found.
	ErrNoSuchPublicAccessBlk = awserrors.NewAWSError("NoSuchPublicAccessBlockConfiguration", "The public access block configuration was not found", http.StatusNotFound)
	// ErrNoSuchReplication is returned when the replication configuration is not found.
	ErrNoSuchReplication = awserrors.NewAWSError("ReplicationConfigurationNotFoundError", "The replication configuration was not found", http.StatusNotFound)
	// ErrNoSuchLogging is returned when the bucket logging configuration is not found.
	ErrNoSuchLogging = awserrors.NewAWSError("NoSuchBucketLogging", "There is no such bucket logging configuration", http.StatusNotFound)
	// ErrNoSuchRequestPayment is returned when the request payment configuration is not found.
	ErrNoSuchRequestPayment = awserrors.NewAWSError("NoSuchRequestPaymentConfiguration", "There is no such request payment configuration", http.StatusNotFound)
	// ErrNoSuchAccelerate is returned when the bucket does not have an accelerate configuration.
	ErrNoSuchAccelerate = awserrors.NewAWSError("NoSuchAccelerateConfiguration", "The specified bucket does not have an accelerate configuration", http.StatusNotFound)
	// ErrNoSuchObjectLock is returned when the object lock configuration does not exist for the bucket.
	ErrNoSuchObjectLock = awserrors.NewAWSError("ObjectLockConfigurationNotFoundError", "Object Lock configuration does not exist for this bucket", http.StatusNotFound)
	// ErrNoSuchNotification is returned when the notification configuration does not exist.
	ErrNoSuchNotification = awserrors.NewAWSError("NotificationConfigurationNotFoundError", "The notification configuration does not exist", http.StatusNotFound)
	// ErrInvalidObjectState is returned when the operation is not valid for the object's storage class.
	ErrInvalidObjectState = awserrors.NewAWSError("InvalidObjectState", "The operation is not valid for the object's storage class", http.StatusConflict)
	// ErrObjectAlreadyRestored is returned when the object is already in the active tier.
	ErrObjectAlreadyRestored = awserrors.NewAWSError("ObjectAlreadyInActiveTier", "The object is already in the active tier", http.StatusOK)
)

// NewNoSuchBucketError creates a NoSuchBucket error for the given bucket name.
func NewNoSuchBucketError(bucketName string) *awserrors.AWSError {
	return awserrors.NewAWSError("NoSuchBucket", "The specified bucket "+bucketName+" does not exist", http.StatusNotFound)
}

// NewBucketAlreadyExistsError creates a BucketAlreadyExists error for the given bucket name.
func NewBucketAlreadyExistsError(bucketName string) *awserrors.AWSError {
	return awserrors.NewAWSError("BucketAlreadyExists", "The requested bucket name "+bucketName+" is not available.", http.StatusConflict)
}

// NewNoSuchKeyError creates a NoSuchKey error for the given key.
func NewNoSuchKeyError(key string) *awserrors.AWSError {
	return awserrors.NewAWSError("NoSuchKey", "The specified key "+key+" does not exist.", http.StatusNotFound)
}

// NewInvalidBucketNameError creates an InvalidBucketName error for the given bucket name.
func NewInvalidBucketNameError(bucketName string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidBucketName", "The specified bucket "+bucketName+" is not valid.", http.StatusBadRequest)
}

// NewInvalidArgumentError creates an InvalidArgument error.
func NewInvalidArgumentError(message string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidArgument", message, http.StatusBadRequest)
}
