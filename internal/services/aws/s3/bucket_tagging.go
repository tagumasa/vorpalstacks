package s3

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

// PutBucketTaggingInput contains the request parameters for the PutBucketTagging operation.
type PutBucketTaggingInput struct {
	Bucket string
	Tags   []Tag
}

// Tag represents an S3 tag key-value pair.
type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

// ToCommon converts an S3 Tag to a common Tag.
func (t Tag) ToCommon() types.Tag {
	return types.Tag{Key: t.Key, Value: t.Value}
}

// TagsToCommon converts a slice of S3 Tags to common Tags.
func TagsToCommon(tags []Tag) []types.Tag {
	result := make([]types.Tag, len(tags))
	for i, t := range tags {
		result[i] = t.ToCommon()
	}
	return result
}

// CommonToTags converts a slice of common Tags to S3 Tags.
func CommonToTags(tags []types.Tag) []Tag {
	result := make([]Tag, len(tags))
	for i, t := range tags {
		result[i] = Tag{Key: t.Key, Value: t.Value}
	}
	return result
}

func validateBucketTags(tags []Tag) error {
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

// PutBucketTagging applies a set of tags to a bucket.
func (o *BucketOperations) PutBucketTagging(ctx *request.RequestContext, input *PutBucketTaggingInput) error {
	if err := validateBucketTags(input.Tags); err != nil {
		return err
	}
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetTags(input.Bucket, TagsToCommon(input.Tags))
}

// GetBucketTaggingInput contains the request parameters for the GetBucketTagging operation.
type GetBucketTaggingInput struct {
	Bucket string
}

// GetBucketTaggingOutput contains the result of the GetBucketTagging operation.
type GetBucketTaggingOutput struct {
	TagSet []Tag `xml:"TagSet>Tag"`
}

// GetBucketTagging retrieves the tags attached to a bucket.
func (o *BucketOperations) GetBucketTagging(ctx *request.RequestContext, input *GetBucketTaggingInput) (*GetBucketTaggingOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	return &GetBucketTaggingOutput{
		TagSet: CommonToTags(bucket.Tags),
	}, nil
}

// DeleteBucketTaggingInput contains the request parameters for the DeleteBucketTagging operation.
type DeleteBucketTaggingInput struct {
	Bucket string
}

// DeleteBucketTagging removes all tags from a bucket.
func (o *BucketOperations) DeleteBucketTagging(ctx *request.RequestContext, input *DeleteBucketTaggingInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetTags(input.Bucket, nil)
}
