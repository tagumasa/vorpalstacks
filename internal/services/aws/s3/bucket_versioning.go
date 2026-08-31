package s3

import (
	"vorpalstacks/internal/common/request"
)

// PutBucketVersioningInput contains the request parameters for the PutBucketVersioning operation.
type PutBucketVersioningInput struct {
	Bucket    string
	Status    string
	MFADelete string
}

// PutBucketVersioning sets the versioning state of a bucket.
func (o *BucketOperations) PutBucketVersioning(ctx *request.RequestContext, input *PutBucketVersioningInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketVersioningCore(store.buckets, input)
}

// GetBucketVersioningInput contains the request parameters for the GetBucketVersioning operation.
type GetBucketVersioningInput struct {
	Bucket string
}

// GetBucketVersioningOutput contains the result of the GetBucketVersioning operation.
type GetBucketVersioningOutput struct {
	Status    string `xml:"Status,omitempty"`
	MFADelete string `xml:"MfaDelete,omitempty"`
}

// GetBucketVersioning retrieves the versioning configuration of a bucket.
func (o *BucketOperations) GetBucketVersioning(ctx *request.RequestContext, input *GetBucketVersioningInput) (*GetBucketVersioningOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketVersioningCore(store.buckets, input)
}
