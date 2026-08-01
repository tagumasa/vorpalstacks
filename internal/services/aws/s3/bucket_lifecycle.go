package s3

import (
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/aws/types"
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
	if input.LifecycleConfiguration == nil || len(input.LifecycleConfiguration.Rules) == 0 {
		return NewInvalidArgumentError("lifecycle configuration must contain at least one rule")
	}

	validStorageClasses := map[string]bool{
		"STANDARD": true, "STANDARD_IA": true, "ONEZONE_IA": true,
		"INTELLIGENT_TIERING": true, "GLACIER": true, "DEEP_ARCHIVE": true,
		"GLACIER_IR": true, "OUTPOSTS": true, "EXPRESS_ONEZONE": true,
		"SNOW": true,
	}

	for _, rule := range input.LifecycleConfiguration.Rules {
		if rule.Status != "Enabled" && rule.Status != "Disabled" {
			return NewInvalidArgumentError("rule status must be Enabled or Disabled")
		}

		if rule.Expiration != nil {
			hasDays := rule.Expiration.Days != nil
			hasDate := rule.Expiration.Date != nil
			if hasDays && hasDate {
				return NewInvalidArgumentError("Expiration Days and Date are mutually exclusive")
			}
			if hasDays && (*rule.Expiration.Days < 1 || *rule.Expiration.Days > 3650) {
				return NewInvalidArgumentError(fmt.Sprintf("Expiration Days must be between 1 and 3650, got %d", *rule.Expiration.Days))
			}
		}

		for _, t := range rule.Transitions {
			hasDays := t.Days != nil
			hasDate := t.Date != nil
			if hasDays && hasDate {
				return NewInvalidArgumentError("Transition Days and Date are mutually exclusive")
			}
			if hasDays && (*t.Days < 0 || *t.Days > 3650) {
				return NewInvalidArgumentError(fmt.Sprintf("Transition Days must be between 0 and 3650, got %d", *t.Days))
			}
			if t.StorageClass != "" && !validStorageClasses[t.StorageClass] {
				return NewInvalidArgumentError(fmt.Sprintf("invalid StorageClass: %s", t.StorageClass))
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			if t.NoncurrentDays != nil {
				nd := *t.NoncurrentDays
				if nd < 1 || nd > 3650 {
					return NewInvalidArgumentError(fmt.Sprintf("NoncurrentVersionTransition NoncurrentDays must be between 1 and 3650, got %d", nd))
				}
			}
			if t.StorageClass != "" && !validStorageClasses[t.StorageClass] {
				return NewInvalidArgumentError(fmt.Sprintf("invalid NoncurrentVersionTransition StorageClass: %s", t.StorageClass))
			}
		}

		if rule.AbortIncompleteMultipartUpload != nil && rule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
			d := *rule.AbortIncompleteMultipartUpload.DaysAfterInitiation
			if d < 1 || d > 3650 {
				return NewInvalidArgumentError(fmt.Sprintf("AbortIncompleteMultipartUpload DaysAfterInitiation must be between 1 and 3650, got %d", d))
			}
		}

		if rule.NoncurrentVersionExpiration != nil && rule.NoncurrentVersionExpiration.NoncurrentDays != nil {
			d := *rule.NoncurrentVersionExpiration.NoncurrentDays
			if d < 1 || d > 3650 {
				return NewInvalidArgumentError(fmt.Sprintf("NoncurrentVersionExpiration NoncurrentDays must be between 1 and 3650, got %d", d))
			}
		}
	}

	var rules []s3store.LifecycleRule
	for _, rule := range input.LifecycleConfiguration.Rules {
		lifecycleRule := s3store.LifecycleRule{
			ID:     rule.ID,
			Status: rule.Status,
		}

		if rule.Filter != nil {
			lifecycleRule.Filter = &s3store.LifecycleRuleFilter{
				Prefix:                rule.Filter.Prefix,
				ObjectSizeGreaterThan: rule.Filter.ObjectSizeGreaterThan,
				ObjectSizeLessThan:    rule.Filter.ObjectSizeLessThan,
			}
			if rule.Filter.Tag != nil {
				lifecycleRule.Filter.Tag = &types.Tag{Key: rule.Filter.Tag.Key, Value: rule.Filter.Tag.Value}
			}
			if rule.Filter.And != nil {
				lifecycleRule.Filter.And = &s3store.LifecycleRuleAndOperator{
					Prefix:                rule.Filter.And.Prefix,
					ObjectSizeGreaterThan: rule.Filter.And.ObjectSizeGreaterThan,
					ObjectSizeLessThan:    rule.Filter.And.ObjectSizeLessThan,
				}
				for _, t := range rule.Filter.And.Tags {
					lifecycleRule.Filter.And.Tags = append(lifecycleRule.Filter.And.Tags, types.Tag{Key: t.Key, Value: t.Value})
				}
			}
		}

		if rule.Expiration != nil {
			lifecycleRule.Expiration = &s3store.LifecycleExpiration{
				Date:                      rule.Expiration.Date,
				Days:                      rule.Expiration.Days,
				ExpiredObjectDeleteMarker: rule.Expiration.ExpiredObjectDeleteMarker,
			}
		}

		for _, t := range rule.Transitions {
			lifecycleRule.Transitions = append(lifecycleRule.Transitions, s3store.LifecycleTransition{
				Date:         t.Date,
				Days:         t.Days,
				StorageClass: s3store.ObjectStorageClass(t.StorageClass),
			})
		}

		if rule.NoncurrentVersionExpiration != nil {
			lifecycleRule.NoncurrentVersionExpiration = &s3store.NoncurrentVersionExpiration{
				NoncurrentDays:          rule.NoncurrentVersionExpiration.NoncurrentDays,
				NewerNoncurrentVersions: rule.NoncurrentVersionExpiration.NewerNoncurrentVersions,
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			lifecycleRule.NoncurrentVersionTransitions = append(lifecycleRule.NoncurrentVersionTransitions, s3store.NoncurrentVersionTransition{
				NoncurrentDays:          t.NoncurrentDays,
				NewerNoncurrentVersions: t.NewerNoncurrentVersions,
				StorageClass:            s3store.ObjectStorageClass(t.StorageClass),
			})
		}

		if rule.AbortIncompleteMultipartUpload != nil {
			lifecycleRule.AbortIncompleteMultipartUpload = &s3store.AbortIncompleteUpload{
				DaysAfterInitiation: rule.AbortIncompleteMultipartUpload.DaysAfterInitiation,
			}
		}

		rules = append(rules, lifecycleRule)
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	return store.buckets.SetLifecycleConfiguration(input.Bucket, &s3store.LifecycleConfiguration{Rules: rules})
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

	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.LifecycleConfiguration == nil {
		return nil, ErrNoSuchLifecycle
	}

	var rules []LifecycleRuleOutput
	for _, rule := range bucket.LifecycleConfiguration.Rules {
		outputRule := LifecycleRuleOutput{
			ID:     rule.ID,
			Status: rule.Status,
		}

		if rule.Filter != nil {
			outputRule.Filter = &LifecycleRuleFilterOutput{
				Prefix:                rule.Filter.Prefix,
				ObjectSizeGreaterThan: rule.Filter.ObjectSizeGreaterThan,
				ObjectSizeLessThan:    rule.Filter.ObjectSizeLessThan,
			}
			if rule.Filter.Tag != nil {
				outputRule.Filter.Tag = &Tag{Key: rule.Filter.Tag.Key, Value: rule.Filter.Tag.Value}
			}
			if rule.Filter.And != nil {
				outputRule.Filter.And = &LifecycleRuleAndOperatorOutput{
					Prefix:                rule.Filter.And.Prefix,
					ObjectSizeGreaterThan: rule.Filter.And.ObjectSizeGreaterThan,
					ObjectSizeLessThan:    rule.Filter.And.ObjectSizeLessThan,
				}
				for _, t := range rule.Filter.And.Tags {
					outputRule.Filter.And.Tags = append(outputRule.Filter.And.Tags, Tag{Key: t.Key, Value: t.Value})
				}
			}
		}

		if rule.Expiration != nil {
			outputRule.Expiration = &LifecycleExpirationOutput{
				Date:                      rule.Expiration.Date,
				Days:                      rule.Expiration.Days,
				ExpiredObjectDeleteMarker: rule.Expiration.ExpiredObjectDeleteMarker,
			}
		}

		for _, t := range rule.Transitions {
			outputRule.Transitions = append(outputRule.Transitions, LifecycleTransitionOutput{
				Date:         t.Date,
				Days:         t.Days,
				StorageClass: string(t.StorageClass),
			})
		}

		if rule.NoncurrentVersionExpiration != nil {
			outputRule.NoncurrentVersionExpiration = &NoncurrentVersionExpirationOutput{
				NoncurrentDays:          rule.NoncurrentVersionExpiration.NoncurrentDays,
				NewerNoncurrentVersions: rule.NoncurrentVersionExpiration.NewerNoncurrentVersions,
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			outputRule.NoncurrentVersionTransitions = append(outputRule.NoncurrentVersionTransitions, NoncurrentVersionTransitionOutput{
				NoncurrentDays:          t.NoncurrentDays,
				NewerNoncurrentVersions: t.NewerNoncurrentVersions,
				StorageClass:            string(t.StorageClass),
			})
		}

		if rule.AbortIncompleteMultipartUpload != nil {
			outputRule.AbortIncompleteMultipartUpload = &AbortIncompleteUploadOutput{
				DaysAfterInitiation: rule.AbortIncompleteMultipartUpload.DaysAfterInitiation,
			}
		}

		rules = append(rules, outputRule)
	}

	return &GetBucketLifecycleConfigurationOutput{Rules: rules}, nil
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
	return store.buckets.SetLifecycleConfiguration(input.Bucket, nil)
}
