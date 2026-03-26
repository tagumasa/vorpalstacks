package s3

import (
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// PutBucketVersioningInput contains the request parameters for the PutBucketVersioning operation.
type PutBucketVersioningInput struct {
	Bucket    string
	Status    string
	MFADelete string
	MFA       string
}

// PutBucketVersioning sets the versioning state of a bucket.
func (o *BucketOperations) PutBucketVersioning(ctx *request.RequestContext, input *PutBucketVersioningInput) error {
	status := s3store.BucketVersioningStatus(input.Status)
	if status != s3store.BucketVersioningEnabled && status != s3store.BucketVersioningSuspended {
		return fmt.Errorf("invalid versioning status")
	}
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetVersioning(input.Bucket, status)
}

// GetBucketVersioningInput contains the request parameters for the GetBucketVersioning operation.
type GetBucketVersioningInput struct {
	Bucket string
}

// GetBucketVersioningOutput contains the result of the GetBucketVersioning operation.
type GetBucketVersioningOutput struct {
	Status    string `xml:"Status,omitempty"`
	MFADelete string `xml:"MFADelete,omitempty"`
}

// GetBucketVersioning retrieves the versioning configuration of a bucket.
func (o *BucketOperations) GetBucketVersioning(ctx *request.RequestContext, input *GetBucketVersioningInput) (*GetBucketVersioningOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.VersioningStatus == "" {
		return &GetBucketVersioningOutput{}, nil
	}

	return &GetBucketVersioningOutput{
		Status: string(bucket.VersioningStatus),
	}, nil
}
