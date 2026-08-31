package s3

import (
	"vorpalstacks/internal/common/request"
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
	return o.svc.getBucketLocationCore(store.buckets, input)
}
