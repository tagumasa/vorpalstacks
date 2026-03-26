package s3

import (
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

func boolPtr(v bool) *bool {
	return &v
}

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
	if len(input.ServerSideEncryptionConfiguration.Rules) == 0 {
		return fmt.Errorf("at least one rule is required")
	}

	rule := input.ServerSideEncryptionConfiguration.Rules[0]
	sseAlgorithm := rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm
	if sseAlgorithm != "AES256" && sseAlgorithm != "aws:kms" {
		return fmt.Errorf("invalid SSE algorithm: %s (must be AES256 or aws:kms)", sseAlgorithm)
	}

	config := &s3store.EncryptionConfig{
		SSEAlgorithm:   rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm,
		KMSMasterKeyID: rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID,
	}
	if rule.BucketKeyEnabled != nil {
		config.BucketKeyEnabled = rule.BucketKeyEnabled
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetEncryption(input.Bucket, config)
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
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.EncryptionConfig == nil {
		return nil, ErrNoSuchEncryption
	}

	return &GetBucketEncryptionOutput{
		ServerSideEncryptionConfiguration: &ServerSideEncryptionConfiguration{
			Rules: []ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: ApplyServerSideEncryptionByDefault{
						SSEAlgorithm:   bucket.EncryptionConfig.SSEAlgorithm,
						KMSMasterKeyID: bucket.EncryptionConfig.KMSMasterKeyID,
					},
					BucketKeyEnabled: boolPtr(bucket.EncryptionConfig.BucketKeyEnabled != nil && *bucket.EncryptionConfig.BucketKeyEnabled),
				},
			},
		},
	}, nil
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
	return store.buckets.SetEncryption(input.Bucket, nil)
}
