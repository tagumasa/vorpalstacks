package s3

import (
	"fmt"

	"vorpalstacks/internal/common/request"
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
	Days  *int32 `xml:"Days,omitempty"`
	Years *int32 `xml:"Years,omitempty"`
}

// PutObjectLockConfiguration applies an object lock configuration to a bucket.
func (o *BucketOperations) PutObjectLockConfiguration(ctx *request.RequestContext, input *PutObjectLockConfigurationInput) error {
	if input.ObjectLockConfiguration == nil {
		return NewInvalidArgumentError("ObjectLockConfiguration is required")
	}
	if err := validateObjectLockEnabled(input.ObjectLockConfiguration.ObjectLockEnabled); err != nil {
		return err
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return ErrObjectLockNotEnabled
	}

	config := &s3store.ObjectLockConfiguration{
		ObjectLockEnabled: input.ObjectLockConfiguration.ObjectLockEnabled,
	}

	if input.ObjectLockConfiguration.Rule != nil && input.ObjectLockConfiguration.Rule.DefaultRetention != nil {
		dr := input.ObjectLockConfiguration.Rule.DefaultRetention

		// Validate Mode: must be GOVERNANCE or COMPLIANCE
		if dr.Mode != string(s3store.ObjectLockRetentionModeGovernance) && dr.Mode != string(s3store.ObjectLockRetentionModeCompliance) {
			return NewInvalidArgumentError(fmt.Sprintf("invalid retention mode: %s (must be GOVERNANCE or COMPLIANCE)", dr.Mode))
		}

		// Days/Years are mutually exclusive: exactly one must be specified
		hasDays := dr.Days != nil
		hasYears := dr.Years != nil
		if hasDays && hasYears {
			return NewInvalidArgumentError("Days and Years are mutually exclusive; specify exactly one")
		}
		if !hasDays && !hasYears {
			return NewInvalidArgumentError("either Days or Years must be specified in DefaultRetention")
		}
		if hasDays {
			if *dr.Days < 1 || *dr.Days > 3650 {
				return NewInvalidArgumentError(fmt.Sprintf("Days must be between 1 and 3650, got %d", *dr.Days))
			}
		}
		if hasYears {
			if *dr.Years < 1 || *dr.Years > 100 {
				return NewInvalidArgumentError(fmt.Sprintf("Years must be between 1 and 100, got %d", *dr.Years))
			}
		}

		config.Rule = &s3store.ObjectLockRule{
			DefaultRetention: &s3store.DefaultRetention{
				Mode:  s3store.ObjectLockRetentionMode(dr.Mode),
				Days:  dr.Days,
				Years: dr.Years,
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
	Days  *int32 `xml:"Days,omitempty"`
	Years *int32 `xml:"Years,omitempty"`
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
