// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// ObjectOperations handles S3 object operations.
type ObjectOperations struct {
	svc *S3Service
}

const (
	maxObjectKeyLength  = 1024
	maxObjectSize       = 5 * 1024 * 1024 * 1024 * 1024
	maxSingleUploadSize = 5 * 1024 * 1024 * 1024
	maxTags             = 50
	maxTagKeyLength     = 128
	maxTagValueLength   = 256
	minPartNumber       = 1
	maxPartNumber       = 10000
	minPartSize         = 5 * 1024 * 1024
)

// NewObjectOperations creates a new ObjectOperations instance.
func NewObjectOperations(svc *S3Service) *ObjectOperations {
	return &ObjectOperations{svc: svc}
}

func validateObjectKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("object key cannot be empty")
	}
	if len(key) > maxObjectKeyLength {
		return fmt.Errorf("object key cannot exceed 1024 bytes")
	}
	if strings.Contains(key, "..") {
		return fmt.Errorf("invalid object key: path traversal detected")
	}
	if strings.Contains(key, "\x00") {
		return fmt.Errorf("invalid object key: null byte detected")
	}
	for i, r := range key {
		if r < 0x20 && r != 0x09 && r != 0x0A && r != 0x0D {
			return fmt.Errorf("object key contains invalid control character at position %d", i)
		}
	}
	return nil
}

func validateTags(tags []Tag) error {
	if len(tags) > maxTags {
		return fmt.Errorf("too many tags (maximum %d)", maxTags)
	}
	for _, tag := range tags {
		if len(tag.Key) == 0 {
			return fmt.Errorf("tag key cannot be empty")
		}
		if len(tag.Key) > maxTagKeyLength {
			return fmt.Errorf("tag key cannot exceed 128 characters")
		}
		if len(tag.Value) > maxTagValueLength {
			return fmt.Errorf("tag value cannot exceed 256 characters")
		}
		if strings.HasPrefix(strings.ToLower(tag.Key), "aws:") {
			return fmt.Errorf("tag key cannot start with 'aws:' (reserved prefix)")
		}
	}
	return nil
}

func validatePartNumber(partNumber int) error {
	if partNumber < minPartNumber || partNumber > maxPartNumber {
		return fmt.Errorf("part number must be between %d and %d", minPartNumber, maxPartNumber)
	}
	return nil
}

func (o *ObjectOperations) validateBucketExists(reqCtx *request.RequestContext, bucket string) error {
	store, err := o.svc.store(reqCtx)
	if err != nil {
		return err
	}
	if !store.buckets.Exists(bucket) {
		return ErrNoSuchBucket
	}
	return nil
}

func formatETag(etag string) string {
	return fmt.Sprintf("\"%s\"", etag)
}

// ObjectStoreObject defines the interface for object storage operations.
type ObjectStoreObject interface {
	GetKey() string
	GetLastModified() time.Time
	GetETag() string
	GetSize() int64
	GetStorageClass() string
	IsDeleteMarker() bool
}

func buildObjectContents(objects []*s3store.Object) []*ObjectContent {
	var contents []*ObjectContent
	for _, obj := range objects {
		if !obj.IsDeleteMarker {
			contents = append(contents, &ObjectContent{
				Key:          obj.Key,
				LastModified: obj.LastModified,
				ETag:         formatETag(obj.ETag),
				Size:         obj.Size,
				StorageClass: string(obj.StorageClass),
			})
		}
	}
	return contents
}
