package s3

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// PutObjectTaggingInput contains the parameters for setting tags on an object.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally targets a specific object version.
// Tags is the list of tag key-value pairs to associate with the object.
type PutObjectTaggingInput struct {
	Bucket    string
	Key       string
	VersionId string
	Tags      []Tag
}

// PutObjectTagging replaces all tags on an object with the specified tag set.
// When input.VersionId is set, only that version's tag set is replaced.
func (o *ObjectOperations) PutObjectTagging(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectTaggingInput) error {
	return o.svc.putObjectTaggingCore(ctx, reqCtx, stores, input)
}

// GetObjectTaggingInput contains the parameters for retrieving an object's tags.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally targets a specific object version.
type GetObjectTaggingInput struct {
	Bucket    string
	Key       string
	VersionId string
}

// GetObjectTaggingOutput contains the result of retrieving an object's tags.
// TagSet contains all tags associated with the object.
type GetObjectTaggingOutput struct {
	TagSet []Tag `xml:"TagSet>Tag"`
}

// GetObjectTagging retrieves all tags associated with an object. When
// input.VersionId is set, the tags of that version are returned.
func (o *ObjectOperations) GetObjectTagging(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *GetObjectTaggingInput) (*GetObjectTaggingOutput, error) {
	return o.svc.getObjectTaggingCore(ctx, stores, input)
}

// DeleteObjectTaggingInput contains the parameters for deleting all tags from
// an object. Bucket is the name of the S3 bucket. Key is the object key
// within the bucket. VersionId optionally targets a specific object version.
type DeleteObjectTaggingInput struct {
	Bucket    string
	Key       string
	VersionId string
}

// DeleteObjectTagging removes all tags from an object. When input.VersionId
// is set, only that version's tag set is removed.
func (o *ObjectOperations) DeleteObjectTagging(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectTaggingInput) error {
	return o.svc.deleteObjectTaggingCore(ctx, reqCtx, stores, input)
}
