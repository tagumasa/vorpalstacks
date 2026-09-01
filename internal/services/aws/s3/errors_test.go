package s3

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3Errors(t *testing.T) {
	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "NoSuchBucket: The specified bucket does not exist", ErrNoSuchBucket.Error())
		assert.Equal(t, http.StatusNotFound, ErrNoSuchBucket.GetHTTPStatusCode())

		assert.Equal(t, "BucketAlreadyExists: The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.", ErrBucketAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrBucketAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "BucketNotEmpty: The bucket you tried to delete is not empty", ErrBucketNotEmpty.Error())
		assert.Equal(t, http.StatusConflict, ErrBucketNotEmpty.GetHTTPStatusCode())

		assert.Equal(t, "NoSuchKey: The specified key does not exist.", ErrNoSuchKey.Error())
		assert.Equal(t, http.StatusNotFound, ErrNoSuchKey.GetHTTPStatusCode())

		assert.Equal(t, "InvalidBucketName: The specified bucket is not valid.", ErrInvalidBucketName.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidBucketName.GetHTTPStatusCode())

		assert.Equal(t, "InvalidRequest: Invalid Request", ErrInvalidRequest.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidRequest.GetHTTPStatusCode())

		assert.Equal(t, "MalformedXML: The XML you provided was not well-formed.", ErrMalformedXML.Error())
		assert.Equal(t, http.StatusBadRequest, ErrMalformedXML.GetHTTPStatusCode())

		assert.Equal(t, "MissingContentLength: You must provide the Content-Length HTTP header.", ErrMissingContentLength.Error())
		assert.Equal(t, http.StatusLengthRequired, ErrMissingContentLength.GetHTTPStatusCode())

		assert.Equal(t, "AccessDenied: Access Denied", ErrAccessDenied.Error())
		assert.Equal(t, http.StatusForbidden, ErrAccessDenied.GetHTTPStatusCode())

		assert.Equal(t, "InvalidCopySource: The copy source is invalid", ErrInvalidCopySource.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidCopySource.GetHTTPStatusCode())

		assert.Equal(t, "PreconditionFailed: At least one of the pre-conditions you specified did not hold", ErrPreconditionFailed.Error())
		assert.Equal(t, http.StatusPreconditionFailed, ErrPreconditionFailed.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchBucketError", func(t *testing.T) {
		err := NewNoSuchBucketError("my-bucket")
		assert.Equal(t, "NoSuchBucket: The specified bucket my-bucket does not exist", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewBucketAlreadyExistsError", func(t *testing.T) {
		err := NewBucketAlreadyExistsError("my-bucket")
		assert.Equal(t, "BucketAlreadyExists: The requested bucket name my-bucket is not available.", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchKeyError", func(t *testing.T) {
		err := NewNoSuchKeyError("my-key")
		assert.Equal(t, "NoSuchKey: The specified key my-key does not exist.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewInvalidBucketNameError", func(t *testing.T) {
		err := NewInvalidBucketNameError("invalid..bucket")
		assert.Equal(t, "InvalidBucketName: The specified bucket invalid..bucket is not valid.", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatusCode())
	})
}

// TestAdminObjectWriteCoresRejectEmptyMembers pins the empty-member
// rejections in the admin object-write cores: bucket and key are required
// on every write, DeleteObjects also requires a non-empty object list,
// and CopyObject requires the copy source. All fire before any store
// access.
func TestAdminObjectWriteCoresRejectEmptyMembers(t *testing.T) {
	svc := &S3Service{}

	assertMsg := func(t *testing.T, err error, want string) {
		t.Helper()
		if err == nil {
			t.Fatalf("expected %q, got nil", want)
		}
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not contain %q", err.Error(), want)
		}
	}

	t.Run("deleteObject", func(t *testing.T) {
		_, err := svc.deleteObjectCore(nil, nil, AdminDeleteObjectInput{})
		assertMsg(t, err, "bucket is required")
		_, err = svc.deleteObjectCore(nil, nil, AdminDeleteObjectInput{Bucket: "b"})
		assertMsg(t, err, "key is required")
	})

	t.Run("putObject", func(t *testing.T) {
		_, err := svc.putObjectCore(nil, nil, nil, AdminPutObjectInput{})
		assertMsg(t, err, "bucket is required")
		_, err = svc.putObjectCore(nil, nil, nil, AdminPutObjectInput{Bucket: "b"})
		assertMsg(t, err, "key is required")
	})

	t.Run("deleteObjects", func(t *testing.T) {
		_, err := svc.deleteObjectsCore(nil, nil, AdminDeleteObjectsInput{})
		assertMsg(t, err, "bucket is required")
		_, err = svc.deleteObjectsCore(nil, nil, AdminDeleteObjectsInput{Bucket: "b"})
		assertMsg(t, err, "no objects specified for deletion")
	})

	t.Run("copyObject", func(t *testing.T) {
		_, err := svc.copyObjectCore(nil, nil, nil, AdminCopyObjectInput{})
		assertMsg(t, err, "bucket is required")
		_, err = svc.copyObjectCore(nil, nil, nil, AdminCopyObjectInput{Bucket: "b"})
		assertMsg(t, err, "key is required")
		_, err = svc.copyObjectCore(nil, nil, nil, AdminCopyObjectInput{Bucket: "b", Key: "k"})
		assertMsg(t, err, "copy source is required")
	})
}
