package s3

import (
	"vorpalstacks/internal/services/aws/common/request"
)

// GetBucketLocationInput contains the input parameters for the GetBucketLocation operation.
type GetBucketLocationInput struct {
	Bucket string
}

// GetBucketLocationOutput contains the output result of the GetBucketLocation operation.
type GetBucketLocationOutput struct {
	LocationConstraint string `xml:",innerxml"`
}

// GetBucketLocation retrieves the location of a bucket.
func (o *BucketOperations) GetBucketLocation(ctx *request.RequestContext, input *GetBucketLocationInput) (*GetBucketLocationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.Region == "us-east-1" {
		return &GetBucketLocationOutput{LocationConstraint: ""}, nil
	}

	return &GetBucketLocationOutput{LocationConstraint: bucket.Region}, nil
}
