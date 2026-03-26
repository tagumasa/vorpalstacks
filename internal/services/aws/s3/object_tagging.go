package s3

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
)

// PutObjectTaggingInput contains the parameters for setting tags on an object.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// Tags is the list of tag key-value pairs to associate with the object.
type PutObjectTaggingInput struct {
	Bucket string
	Key    string
	Tags   []Tag
}

// PutObjectTagging replaces all tags on an object with the specified tag set.
func (o *ObjectOperations) PutObjectTagging(ctx context.Context, reqCtx *request.RequestContext, input *PutObjectTaggingInput) error {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return err
	}

	if err := validateTags(input.Tags); err != nil {
		return err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return err
	}

	return store.objects.SetTags(input.Bucket, input.Key, TagsToCommon(input.Tags))
}

// GetObjectTaggingInput contains the parameters for retrieving an object's tags.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
type GetObjectTaggingInput struct {
	Bucket string
	Key    string
}

// GetObjectTaggingOutput contains the result of retrieving an object's tags.
// TagSet contains all tags associated with the object.
type GetObjectTaggingOutput struct {
	TagSet []Tag `xml:"TagSet>Tag"`
}

// GetObjectTagging retrieves all tags associated with an object.
func (o *ObjectOperations) GetObjectTagging(ctx context.Context, reqCtx *request.RequestContext, input *GetObjectTaggingInput) (*GetObjectTaggingOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	obj, err := store.objects.GetMetadata(input.Bucket, input.Key)
	if err != nil {
		return nil, err
	}

	return &GetObjectTaggingOutput{TagSet: CommonToTags(obj.Tags)}, nil
}

// DeleteObjectTaggingInput contains the parameters for deleting all tags from an object.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
type DeleteObjectTaggingInput struct {
	Bucket string
	Key    string
}

// DeleteObjectTagging removes all tags from an object.
func (o *ObjectOperations) DeleteObjectTagging(ctx context.Context, reqCtx *request.RequestContext, input *DeleteObjectTaggingInput) error {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return err
	}

	return store.objects.SetTags(input.Bucket, input.Key, nil)
}
