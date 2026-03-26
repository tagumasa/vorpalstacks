// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"net/http"

	"vorpalstacks/internal/services/aws/common/errors"
)

// Predefined S3 error variables.
var (
	ErrNoSuchBucket          = errors.NewAWSError("NoSuchBucket", "The specified bucket does not exist", http.StatusNotFound)
	ErrBucketAlreadyExists   = errors.NewAWSError("BucketAlreadyExists", "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.", http.StatusConflict)
	ErrBucketNotEmpty        = errors.NewAWSError("BucketNotEmpty", "The bucket you tried to delete is not empty", http.StatusConflict)
	ErrNoSuchKey             = errors.NewAWSError("NoSuchKey", "The specified key does not exist.", http.StatusNotFound)
	ErrInvalidBucketName     = errors.NewAWSError("InvalidBucketName", "The specified bucket is not valid.", http.StatusBadRequest)
	ErrInvalidRequest        = errors.NewAWSError("InvalidRequest", "Invalid Request", http.StatusBadRequest)
	ErrMalformedXML          = errors.NewAWSError("MalformedXML", "The XML you provided was not well-formed.", http.StatusBadRequest)
	ErrMissingContentLength  = errors.NewAWSError("MissingContentLength", "You must provide the Content-Length HTTP header.", http.StatusLengthRequired)
	ErrAccessDenied          = errors.NewAWSError("AccessDenied", "Access Denied", http.StatusForbidden)
	ErrInvalidCopySource     = errors.NewAWSError("InvalidCopySource", "The copy source is invalid", http.StatusBadRequest)
	ErrPreconditionFailed    = errors.NewAWSError("PreconditionFailed", "At least one of the pre-conditions you specified did not hold", http.StatusPreconditionFailed)
	ErrNotModified           = errors.NewAWSError("NotModified", "The resource was not modified", http.StatusNotModified)
	ErrInvalidRange          = errors.NewAWSError("InvalidRange", "The requested range cannot be satisfied", http.StatusRequestedRangeNotSatisfiable)
	ErrInvalidSSECustomerKey = errors.NewAWSError("InvalidSSECustomerKey", "The provided SSE customer key is not valid", http.StatusBadRequest)
	ErrNoSuchUpload          = errors.NewAWSError("NoSuchUpload", "The specified multipart upload does not exist.", http.StatusNotFound)
	ErrInvalidPart           = errors.NewAWSError("InvalidPart", "One or more of the specified parts could not be found.", http.StatusBadRequest)
	ErrEntityTooLarge        = errors.NewAWSError("EntityTooLarge", "Your proposed upload exceeds the maximum allowed size", http.StatusBadRequest)
	ErrInvalidArgument       = errors.NewAWSError("InvalidArgument", "Invalid argument", http.StatusBadRequest)
	ErrNoSuchBucketPolicy    = errors.NewAWSError("NoSuchBucketPolicy", "The bucket policy does not exist", http.StatusNotFound)
	ErrNoSuchCORS            = errors.NewAWSError("NoSuchCORSConfiguration", "The CORS configuration does not exist", http.StatusNotFound)
	ErrNoSuchEncryption      = errors.NewAWSError("ServerSideEncryptionConfigurationNotFoundError", "The server side encryption configuration was not found", http.StatusNotFound)
	ErrNoSuchTagging         = errors.NewAWSError("NoSuchTagSet", "The TagSet does not exist", http.StatusNotFound)
	ErrNoSuchVersioning      = errors.NewAWSError("NoSuchBucketVersioningConfiguration", "The bucket does not have a versioning configuration", http.StatusNotFound)
	ErrNoSuchWebsite         = errors.NewAWSError("NoSuchWebsiteConfiguration", "The specified bucket does not have a website configuration", http.StatusNotFound)
	ErrNoSuchLifecycle       = errors.NewAWSError("NoSuchLifecycleConfiguration", "The lifecycle configuration does not exist", http.StatusNotFound)
	ErrNoSuchOwnershipCtrls  = errors.NewAWSError("OwnershipControlsNotFoundError", "The bucket does not have OwnershipControls", http.StatusNotFound)
	ErrNoSuchPublicAccessBlk = errors.NewAWSError("NoSuchPublicAccessBlockConfiguration", "The public access block configuration was not found", http.StatusNotFound)
	ErrNoSuchReplication     = errors.NewAWSError("ReplicationConfigurationNotFoundError", "The replication configuration was not found", http.StatusNotFound)
	ErrNoSuchLogging         = errors.NewAWSError("NoSuchBucketLogging", "There is no such bucket logging configuration", http.StatusNotFound)
	ErrNoSuchRequestPayment  = errors.NewAWSError("NoSuchRequestPaymentConfiguration", "There is no such request payment configuration", http.StatusNotFound)
	ErrNoSuchAccelerate      = errors.NewAWSError("NoSuchAccelerateConfiguration", "The specified bucket does not have an accelerate configuration", http.StatusNotFound)
	ErrNoSuchObjectLock      = errors.NewAWSError("ObjectLockConfigurationNotFoundError", "Object Lock configuration does not exist for this bucket", http.StatusNotFound)
	ErrNoSuchNotification    = errors.NewAWSError("NotificationConfigurationNotFoundError", "The notification configuration does not exist", http.StatusNotFound)
	ErrInvalidObjectState    = errors.NewAWSError("InvalidObjectState", "The operation is not valid for the object's storage class", http.StatusConflict)
	ErrObjectAlreadyRestored = errors.NewAWSError("ObjectAlreadyInActiveTier", "The object is already in the active tier", http.StatusOK)
)

// NewNoSuchBucketError creates a NoSuchBucket error for the given bucket name.
func NewNoSuchBucketError(bucketName string) *errors.AWSError {
	return errors.NewAWSError("NoSuchBucket", "The specified bucket "+bucketName+" does not exist", http.StatusNotFound)
}

// NewBucketAlreadyExistsError creates a BucketAlreadyExists error for the given bucket name.
func NewBucketAlreadyExistsError(bucketName string) *errors.AWSError {
	return errors.NewAWSError("BucketAlreadyExists", "The requested bucket name "+bucketName+" is not available.", http.StatusConflict)
}

// NewNoSuchKeyError creates a NoSuchKey error for the given key.
func NewNoSuchKeyError(key string) *errors.AWSError {
	return errors.NewAWSError("NoSuchKey", "The specified key "+key+" does not exist.", http.StatusNotFound)
}

// NewInvalidBucketNameError creates an InvalidBucketName error for the given bucket name.
func NewInvalidBucketNameError(bucketName string) *errors.AWSError {
	return errors.NewAWSError("InvalidBucketName", "The specified bucket "+bucketName+" is not valid.", http.StatusBadRequest)
}

// NewInvalidArgumentError creates an InvalidArgument error.
func NewInvalidArgumentError(message string) *errors.AWSError {
	return errors.NewAWSError("InvalidArgument", message, http.StatusBadRequest)
}
