package s3

import (
	"fmt"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// ObjectOperations handles S3 object operations.
type ObjectOperations struct {
	svc *S3Service
}

// NewObjectOperations creates a new ObjectOperations instance.
func NewObjectOperations(svc *S3Service) *ObjectOperations {
	return &ObjectOperations{svc: svc}
}

func (o *ObjectOperations) validateBucketExists(stores *s3Stores, bucket string) error {
	if !stores.buckets.Exists(bucket) {
		return ErrNoSuchBucket
	}
	return nil
}

func formatETag(etag string) string {
	return fmt.Sprintf("\"%s\"", etag)
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
