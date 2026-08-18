package s3

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
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

// resolveTaggingTarget returns the object record a tagging operation
// targets: the specific version when versionId is set, the current object
// otherwise. Both tagging reads and writes go through it so that a
// nonexistent object and a nonexistent version are reported uniformly.
func (o *ObjectOperations) resolveTaggingTarget(ctx context.Context, stores *s3Stores, bucket, key, versionId string) (*s3store.Object, error) {
	if versionId != "" {
		return stores.objects.HeadWithVersion(ctx, bucket, key, versionId)
	}
	return stores.objects.GetMetadata(bucket, key)
}

// PutObjectTagging replaces all tags on an object with the specified tag set.
// When input.VersionId is set, only that version's tag set is replaced.
func (o *ObjectOperations) PutObjectTagging(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectTaggingInput) error {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	obj, err := o.resolveTaggingTarget(ctx, stores, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return versionLookupError(input.Key, input.VersionId)
	}

	if err := validateTags(input.Tags); err != nil {
		return err
	}

	if err := stores.objects.SetTags(input.Bucket, input.Key, input.VersionId, TagsToCommon(input.Tags)); err != nil {
		return err
	}

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectTaggingPut)
	return nil
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
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	obj, err := o.resolveTaggingTarget(ctx, stores, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, versionLookupError(input.Key, input.VersionId)
	}

	return &GetObjectTaggingOutput{TagSet: CommonToTags(obj.Tags)}, nil
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
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	obj, err := o.resolveTaggingTarget(ctx, stores, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return versionLookupError(input.Key, input.VersionId)
	}

	if err := stores.objects.SetTags(input.Bucket, input.Key, input.VersionId, nil); err != nil {
		return err
	}

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectTaggingDelete)
	return nil
}
