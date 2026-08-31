package s3

import (
	"vorpalstacks/internal/common/request"
)

// PutBucketEncryptionInput is the input for PutBucketEncryption.
type PutBucketEncryptionInput struct {
	Bucket                            string
	ServerSideEncryptionConfiguration *ServerSideEncryptionConfiguration
}

// ServerSideEncryptionConfiguration defines the encryption rules for a bucket.
type ServerSideEncryptionConfiguration struct {
	Rules []ServerSideEncryptionRule `xml:"Rule"`
}

// ServerSideEncryptionRule defines a server-side encryption rule.
type ServerSideEncryptionRule struct {
	ApplyServerSideEncryptionByDefault ApplyServerSideEncryptionByDefault `xml:"ApplyServerSideEncryptionByDefault"`
	BucketKeyEnabled                   *bool                              `xml:"BucketKeyEnabled"`
}

// ApplyServerSideEncryptionByDefault defines the default encryption settings.
type ApplyServerSideEncryptionByDefault struct {
	SSEAlgorithm   string `xml:"SSEAlgorithm"`
	KMSMasterKeyID string `xml:"KMSMasterKeyID,omitempty"`
}

// PutBucketEncryption sets the encryption configuration for an S3 bucket.
func (o *BucketOperations) PutBucketEncryption(ctx *request.RequestContext, input *PutBucketEncryptionInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketEncryptionCore(store.buckets, input)
}

// GetBucketEncryptionInput is the input for GetBucketEncryption.
type GetBucketEncryptionInput struct {
	Bucket string
}

// GetBucketEncryptionOutput is the output of GetBucketEncryption.
type GetBucketEncryptionOutput struct {
	ServerSideEncryptionConfiguration *ServerSideEncryptionConfiguration `xml:"ServerSideEncryptionConfiguration"`
}

// GetBucketEncryption retrieves the encryption configuration for an S3 bucket.
func (o *BucketOperations) GetBucketEncryption(ctx *request.RequestContext, input *GetBucketEncryptionInput) (*GetBucketEncryptionOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketEncryptionCore(store.buckets, input)
}

// DeleteBucketEncryptionInput is the input for DeleteBucketEncryption.
type DeleteBucketEncryptionInput struct {
	Bucket string
}

// DeleteBucketEncryption removes the encryption configuration from an S3 bucket.
func (o *BucketOperations) DeleteBucketEncryption(ctx *request.RequestContext, input *DeleteBucketEncryptionInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketEncryptionCore(store.buckets, input)
}
