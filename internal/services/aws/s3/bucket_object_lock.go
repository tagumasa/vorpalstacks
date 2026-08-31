package s3

import (
	"vorpalstacks/internal/common/request"
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
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putObjectLockConfigurationCore(store.buckets, input)
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
	return o.svc.getObjectLockConfigurationCore(store.buckets, input)
}
