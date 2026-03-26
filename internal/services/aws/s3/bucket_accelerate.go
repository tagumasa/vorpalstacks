package s3

import (
	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// AccelerateConfigurationInput represents the configuration for S3 Accelerate.
type AccelerateConfigurationInput struct {
	Status string `xml:"Status"`
}

// PutBucketAccelerateConfigurationInput is the input for PutBucketAccelerateConfiguration.
type PutBucketAccelerateConfigurationInput struct {
	Bucket                  string
	AccelerateConfiguration *AccelerateConfigurationInput
}

// PutBucketAccelerateConfiguration sets the accelerate configuration for an S3 bucket.
// This enables or disables S3 Transfer Acceleration on the bucket.
func (o *BucketOperations) PutBucketAccelerateConfiguration(ctx *request.RequestContext, input *PutBucketAccelerateConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	config := &s3store.AccelerateConfiguration{
		Status: input.AccelerateConfiguration.Status,
	}

	return store.buckets.SetAccelerateConfiguration(input.Bucket, config)
}

// GetBucketAccelerateConfigurationInput is the input for GetBucketAccelerateConfiguration.
type GetBucketAccelerateConfigurationInput struct {
	Bucket string
}

// GetBucketAccelerateConfigurationOutput is the output of GetBucketAccelerateConfiguration.
type GetBucketAccelerateConfigurationOutput struct {
	AccelerateConfiguration *AccelerateConfigurationInput `xml:"AccelerateConfiguration"`
}

// GetBucketAccelerateConfiguration retrieves the accelerate configuration for an S3 bucket.
func (o *BucketOperations) GetBucketAccelerateConfiguration(ctx *request.RequestContext, input *GetBucketAccelerateConfigurationInput) (*GetBucketAccelerateConfigurationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	config, err := store.buckets.GetAccelerateConfiguration(input.Bucket)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return &GetBucketAccelerateConfigurationOutput{
			AccelerateConfiguration: &AccelerateConfigurationInput{},
		}, nil
	}

	return &GetBucketAccelerateConfigurationOutput{
		AccelerateConfiguration: &AccelerateConfigurationInput{
			Status: config.Status,
		},
	}, nil
}
