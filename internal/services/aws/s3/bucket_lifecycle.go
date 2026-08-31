package s3

import (
	"time"

	"vorpalstacks/internal/common/request"
)

// PutBucketLifecycleConfigurationInput is the input for PutBucketLifecycleConfiguration.
type PutBucketLifecycleConfigurationInput struct {
	Bucket                 string
	LifecycleConfiguration *LifecycleConfigurationInput
}

// LifecycleConfigurationInput defines the lifecycle configuration for a bucket.
type LifecycleConfigurationInput struct {
	Rules []LifecycleRuleInput `xml:"Rule"`
}

// LifecycleRuleInput defines a lifecycle rule for bucket objects.
type LifecycleRuleInput struct {
	ID                             string                             `xml:"ID"`
	Status                         string                             `xml:"Status"`
	Filter                         *LifecycleRuleFilterInput          `xml:"Filter,omitempty"`
	Expiration                     *LifecycleExpirationInput          `xml:"Expiration,omitempty"`
	Transitions                    []LifecycleTransitionInput         `xml:"Transition,omitempty"`
	NoncurrentVersionExpiration    *NoncurrentVersionExpirationInput  `xml:"NoncurrentVersionExpiration,omitempty"`
	NoncurrentVersionTransitions   []NoncurrentVersionTransitionInput `xml:"NoncurrentVersionTransition,omitempty"`
	AbortIncompleteMultipartUpload *AbortIncompleteUploadInput        `xml:"AbortIncompleteMultipartUpload,omitempty"`
}

// LifecycleRuleFilterInput defines the filter for a lifecycle rule.
type LifecycleRuleFilterInput struct {
	Prefix                string                         `xml:"Prefix,omitempty"`
	ObjectSizeGreaterThan *int64                         `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64                         `xml:"ObjectSizeLessThan,omitempty"`
	And                   *LifecycleRuleAndOperatorInput `xml:"And,omitempty"`
	Tag                   *Tag                           `xml:"Tag,omitempty"`
}

// LifecycleRuleAndOperatorInput defines multiple filters for a lifecycle rule.
type LifecycleRuleAndOperatorInput struct {
	Prefix                string `xml:"Prefix,omitempty"`
	Tags                  []Tag  `xml:"Tags>Tag,omitempty"`
	ObjectSizeGreaterThan *int64 `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64 `xml:"ObjectSizeLessThan,omitempty"`
}

// LifecycleExpirationInput defines when objects expire.
type LifecycleExpirationInput struct {
	Date                      *time.Time `xml:"Date,omitempty"`
	Days                      *int32     `xml:"Days,omitempty"`
	ExpiredObjectDeleteMarker *bool      `xml:"ExpiredObjectDeleteMarker,omitempty"`
}

// LifecycleTransitionInput defines when objects transition to another storage class.
type LifecycleTransitionInput struct {
	Date         *time.Time `xml:"Date,omitempty"`
	Days         *int32     `xml:"Days,omitempty"`
	StorageClass string     `xml:"StorageClass"`
}

// NoncurrentVersionExpirationInput defines when noncurrent versions expire.
type NoncurrentVersionExpirationInput struct {
	NoncurrentDays          *int32 `xml:"NoncurrentDays,omitempty"`
	NewerNoncurrentVersions *int32 `xml:"NewerNoncurrentVersions,omitempty"`
}

// NoncurrentVersionTransitionInput defines when noncurrent versions transition.
type NoncurrentVersionTransitionInput struct {
	NoncurrentDays          *int32 `xml:"NoncurrentDays,omitempty"`
	NewerNoncurrentVersions *int32 `xml:"NewerNoncurrentVersions,omitempty"`
	StorageClass            string `xml:"StorageClass"`
}

// AbortIncompleteUploadInput defines when incomplete multipart uploads are aborted.
type AbortIncompleteUploadInput struct {
	DaysAfterInitiation *int32 `xml:"DaysAfterInitiation,omitempty"`
}

// PutBucketLifecycleConfiguration sets the lifecycle configuration for an S3 bucket.
func (o *BucketOperations) PutBucketLifecycleConfiguration(ctx *request.RequestContext, input *PutBucketLifecycleConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketLifecycleConfigurationCore(store.buckets, input)
}

// GetBucketLifecycleConfigurationInput is the input for GetBucketLifecycleConfiguration.
type GetBucketLifecycleConfigurationInput struct {
	Bucket string
}

// GetBucketLifecycleConfigurationOutput is the output of GetBucketLifecycleConfiguration.
type GetBucketLifecycleConfigurationOutput struct {
	Rules []LifecycleRuleOutput `xml:"Rule"`
}

// LifecycleRuleOutput represents a lifecycle rule in the output.
type LifecycleRuleOutput struct {
	ID                             string                              `xml:"ID"`
	Status                         string                              `xml:"Status"`
	Filter                         *LifecycleRuleFilterOutput          `xml:"Filter,omitempty"`
	Expiration                     *LifecycleExpirationOutput          `xml:"Expiration,omitempty"`
	Transitions                    []LifecycleTransitionOutput         `xml:"Transition,omitempty"`
	NoncurrentVersionExpiration    *NoncurrentVersionExpirationOutput  `xml:"NoncurrentVersionExpiration,omitempty"`
	NoncurrentVersionTransitions   []NoncurrentVersionTransitionOutput `xml:"NoncurrentVersionTransition,omitempty"`
	AbortIncompleteMultipartUpload *AbortIncompleteUploadOutput        `xml:"AbortIncompleteMultipartUpload,omitempty"`
}

// LifecycleRuleFilterOutput represents the filter in the output.
type LifecycleRuleFilterOutput struct {
	Prefix                string                          `xml:"Prefix,omitempty"`
	ObjectSizeGreaterThan *int64                          `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64                          `xml:"ObjectSizeLessThan,omitempty"`
	And                   *LifecycleRuleAndOperatorOutput `xml:"And,omitempty"`
	Tag                   *Tag                            `xml:"Tag,omitempty"`
}

// LifecycleRuleAndOperatorOutput represents multiple filters in the output.
type LifecycleRuleAndOperatorOutput struct {
	Prefix                string `xml:"Prefix,omitempty"`
	Tags                  []Tag  `xml:"Tags>Tag,omitempty"`
	ObjectSizeGreaterThan *int64 `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64 `xml:"ObjectSizeLessThan,omitempty"`
}

// LifecycleExpirationOutput represents expiration in the output.
type LifecycleExpirationOutput struct {
	Date                      *time.Time `xml:"Date,omitempty"`
	Days                      *int32     `xml:"Days,omitempty"`
	ExpiredObjectDeleteMarker *bool      `xml:"ExpiredObjectDeleteMarker,omitempty"`
}

// LifecycleTransitionOutput represents a transition in the output.
type LifecycleTransitionOutput struct {
	Date         *time.Time `xml:"Date,omitempty"`
	Days         *int32     `xml:"Days,omitempty"`
	StorageClass string     `xml:"StorageClass"`
}

// NoncurrentVersionExpirationOutput represents noncurrent version expiration.
type NoncurrentVersionExpirationOutput struct {
	NoncurrentDays          *int32 `xml:"NoncurrentDays,omitempty"`
	NewerNoncurrentVersions *int32 `xml:"NewerNoncurrentVersions,omitempty"`
}

// NoncurrentVersionTransitionOutput represents noncurrent version transition.
type NoncurrentVersionTransitionOutput struct {
	NoncurrentDays          *int32 `xml:"NoncurrentDays,omitempty"`
	NewerNoncurrentVersions *int32 `xml:"NewerNoncurrentVersions,omitempty"`
	StorageClass            string `xml:"StorageClass"`
}

// AbortIncompleteUploadOutput represents abort incomplete upload settings.
type AbortIncompleteUploadOutput struct {
	DaysAfterInitiation *int32 `xml:"DaysAfterInitiation,omitempty"`
}

// GetBucketLifecycleConfiguration retrieves the lifecycle configuration for an S3 bucket.
func (o *BucketOperations) GetBucketLifecycleConfiguration(ctx *request.RequestContext, input *GetBucketLifecycleConfigurationInput) (*GetBucketLifecycleConfigurationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketLifecycleConfigurationCore(store.buckets, input)
}

// DeleteBucketLifecycleConfigurationInput is the input for DeleteBucketLifecycleConfiguration.
type DeleteBucketLifecycleConfigurationInput struct {
	Bucket string
}

// DeleteBucketLifecycleConfiguration removes the lifecycle configuration from an S3 bucket.
func (o *BucketOperations) DeleteBucketLifecycleConfiguration(ctx *request.RequestContext, input *DeleteBucketLifecycleConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketLifecycleConfigurationCore(store.buckets, input)
}
