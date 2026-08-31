package s3

import (
	"vorpalstacks/internal/common/request"
)

// PublicAccessBlockConfiguration defines the public access block settings for a bucket.
type PublicAccessBlockConfiguration struct {
	BlockPublicAcls       bool `xml:"BlockPublicAcls"`
	BlockPublicPolicy     bool `xml:"BlockPublicPolicy"`
	IgnorePublicAcls      bool `xml:"IgnorePublicAcls"`
	RestrictPublicBuckets bool `xml:"RestrictPublicBuckets"`
}

// PutPublicAccessBlockInput represents the input for setting public access block configuration.
type PutPublicAccessBlockInput struct {
	Bucket                         string
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration
}

// PutPublicAccessBlock applies public access block configuration to an S3 bucket.
func (o *BucketOperations) PutPublicAccessBlock(ctx *request.RequestContext, input *PutPublicAccessBlockInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putPublicAccessBlockCore(store.buckets, input)
}

// GetPublicAccessBlockInput represents the input for retrieving public access block configuration.
type GetPublicAccessBlockInput struct {
	Bucket string
}

// GetPublicAccessBlockOutput represents the output from retrieving public access block configuration.
type GetPublicAccessBlockOutput struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `xml:"PublicAccessBlockConfiguration"`
}

// GetPublicAccessBlock retrieves public access block configuration for an S3 bucket.
func (o *BucketOperations) GetPublicAccessBlock(ctx *request.RequestContext, input *GetPublicAccessBlockInput) (*GetPublicAccessBlockOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getPublicAccessBlockCore(store.buckets, input)
}

// DeletePublicAccessBlockInput represents the input for deleting public access block configuration.
type DeletePublicAccessBlockInput struct {
	Bucket string
}

// DeletePublicAccessBlock removes public access block configuration from an S3 bucket.
func (o *BucketOperations) DeletePublicAccessBlock(ctx *request.RequestContext, input *DeletePublicAccessBlockInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deletePublicAccessBlockCore(store.buckets, input)
}
