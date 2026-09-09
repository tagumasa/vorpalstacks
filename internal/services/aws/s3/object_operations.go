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
	return o.svc.validateBucketExists(stores, bucket)
}

func formatETag(etag string) string {
	return fmt.Sprintf("\"%s\"", etag)
}

// bucketOwner is the owner identity list operations render on their
// entries. AWS identifies the bucket owner by account, and the platform's
// single-tenant account doubles as the display name.
func (s *S3Service) bucketOwner() *Owner {
	return &Owner{ID: s.accountID, DisplayName: s.accountID}
}

func buildObjectContents(objects []*s3store.Object, owner *Owner) []*ObjectContent {
	var contents []*ObjectContent
	for _, obj := range objects {
		if !obj.IsDeleteMarker {
			contents = append(contents, &ObjectContent{
				Key:          obj.Key,
				LastModified: obj.LastModified,
				ETag:         formatETag(obj.ETag),
				Size:         obj.Size,
				StorageClass: string(obj.StorageClass),
				Owner:        owner,
			})
		}
	}
	return contents
}
