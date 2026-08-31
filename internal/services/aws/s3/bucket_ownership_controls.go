package s3

import (
	"vorpalstacks/internal/common/request"
)

// OwnershipControlsInput represents the input for bucket ownership controls configuration.
type OwnershipControlsInput struct {
	Rules []OwnershipControlsRuleInput `xml:"Rule"`
}

// OwnershipControlsRuleInput defines an ownership rule for bucket ownership controls.
type OwnershipControlsRuleInput struct {
	ObjectOwnership string `xml:"ObjectOwnership"`
}

// PutBucketOwnershipControlsInput represents the input for setting bucket ownership controls.
type PutBucketOwnershipControlsInput struct {
	Bucket            string
	OwnershipControls *OwnershipControlsInput
}

// PutBucketOwnershipControls sets ownership controls for an S3 bucket.
func (o *BucketOperations) PutBucketOwnershipControls(ctx *request.RequestContext, input *PutBucketOwnershipControlsInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketOwnershipControlsCore(store.buckets, input)
}

// GetBucketOwnershipControlsInput represents the input for retrieving bucket ownership controls.
type GetBucketOwnershipControlsInput struct {
	Bucket string
}

// GetBucketOwnershipControlsOutput represents the output from retrieving bucket ownership controls.
type GetBucketOwnershipControlsOutput struct {
	OwnershipControls *OwnershipControlsInput `xml:"OwnershipControls"`
}

// GetBucketOwnershipControls retrieves ownership controls for an S3 bucket.
func (o *BucketOperations) GetBucketOwnershipControls(ctx *request.RequestContext, input *GetBucketOwnershipControlsInput) (*GetBucketOwnershipControlsOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketOwnershipControlsCore(store.buckets, input)
}

// DeleteBucketOwnershipControlsInput represents the input for deleting bucket ownership controls.
type DeleteBucketOwnershipControlsInput struct {
	Bucket string
}

// DeleteBucketOwnershipControls deletes ownership controls from an S3 bucket.
func (o *BucketOperations) DeleteBucketOwnershipControls(ctx *request.RequestContext, input *DeleteBucketOwnershipControlsInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketOwnershipControlsCore(store.buckets, input)
}
