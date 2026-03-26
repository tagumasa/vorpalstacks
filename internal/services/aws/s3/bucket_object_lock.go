package s3

import (
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// PutObjectLockConfigurationInput contains the input parameters for the PutObjectLockConfiguration operation.
type PutObjectLockConfigurationInput struct {
	Bucket                  string
	ObjectLockConfiguration *ObjectLockConfigurationInput
}

// ObjectLockConfigurationInput contains the object lock configuration to set.
type ObjectLockConfigurationInput struct {
	ObjectLockEnabled string               `xml:"ObjectLockEnabled"`
	Rule              *ObjectLockRuleInput `xml:"Rule,omitempty"`
}

// ObjectLockRuleInput contains the default retention rule for object lock.
type ObjectLockRuleInput struct {
	DefaultRetention *DefaultRetentionInput `xml:"DefaultRetention,omitempty"`
}

// DefaultRetentionInput contains the default retention period for objects.
type DefaultRetentionInput struct {
	Mode  string `xml:"Mode"`
	Days  *int   `xml:"Days,omitempty"`
	Years *int   `xml:"Years,omitempty"`
}

// PutObjectLockConfiguration applies an object lock configuration to a bucket.
func (o *BucketOperations) PutObjectLockConfiguration(ctx *request.RequestContext, input *PutObjectLockConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return fmt.Errorf("object lock is not enabled on this bucket")
	}

	config := &s3store.ObjectLockConfiguration{
		ObjectLockEnabled: input.ObjectLockConfiguration.ObjectLockEnabled,
	}

	if input.ObjectLockConfiguration.Rule != nil && input.ObjectLockConfiguration.Rule.DefaultRetention != nil {
		config.Rule = &s3store.ObjectLockRule{
			DefaultRetention: &s3store.DefaultRetention{
				Mode:  s3store.ObjectLockRetentionMode(input.ObjectLockConfiguration.Rule.DefaultRetention.Mode),
				Days:  input.ObjectLockConfiguration.Rule.DefaultRetention.Days,
				Years: input.ObjectLockConfiguration.Rule.DefaultRetention.Years,
			},
		}
	}

	return store.buckets.SetObjectLockConfiguration(input.Bucket, config)
}

// GetObjectLockConfigurationInput contains the input parameters for the GetObjectLockConfiguration operation.
type GetObjectLockConfigurationInput struct {
	Bucket string
}

// GetObjectLockConfigurationOutput contains the output result of the GetObjectLockConfiguration operation.
type GetObjectLockConfigurationOutput struct {
	ObjectLockConfiguration *ObjectLockConfigurationOutput `xml:"ObjectLockConfiguration"`
}

// ObjectLockConfigurationOutput contains the object lock configuration for a bucket.
type ObjectLockConfigurationOutput struct {
	ObjectLockEnabled string                `xml:"ObjectLockEnabled"`
	Rule              *ObjectLockRuleOutput `xml:"Rule,omitempty"`
}

// ObjectLockRuleOutput contains the default retention rule output.
type ObjectLockRuleOutput struct {
	DefaultRetention *DefaultRetentionOutput `xml:"DefaultRetention,omitempty"`
}

// DefaultRetentionOutput contains the default retention period output.
type DefaultRetentionOutput struct {
	Mode  string `xml:"Mode"`
	Days  *int   `xml:"Days,omitempty"`
	Years *int   `xml:"Years,omitempty"`
}

// GetObjectLockConfiguration retrieves the object lock configuration for a bucket.
func (o *BucketOperations) GetObjectLockConfiguration(ctx *request.RequestContext, input *GetObjectLockConfigurationInput) (*GetObjectLockConfigurationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if !bucket.ObjectLockEnabled {
		return nil, ErrNoSuchObjectLock
	}

	if bucket.ObjectLockConfig == nil {
		return &GetObjectLockConfigurationOutput{
			ObjectLockConfiguration: &ObjectLockConfigurationOutput{
				ObjectLockEnabled: "Enabled",
			},
		}, nil
	}

	output := &GetObjectLockConfigurationOutput{
		ObjectLockConfiguration: &ObjectLockConfigurationOutput{
			ObjectLockEnabled: bucket.ObjectLockConfig.ObjectLockEnabled,
		},
	}

	if bucket.ObjectLockConfig.Rule != nil && bucket.ObjectLockConfig.Rule.DefaultRetention != nil {
		output.ObjectLockConfiguration.Rule = &ObjectLockRuleOutput{
			DefaultRetention: &DefaultRetentionOutput{
				Mode:  string(bucket.ObjectLockConfig.Rule.DefaultRetention.Mode),
				Days:  bucket.ObjectLockConfig.Rule.DefaultRetention.Days,
				Years: bucket.ObjectLockConfig.Rule.DefaultRetention.Years,
			},
		}
	}

	return output, nil
}
