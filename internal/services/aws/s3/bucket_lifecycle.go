package s3

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
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

// lifecycleDayUnit is the duration one lifecycle day occupies. TEST_MODE
// compresses it to one second so rule windows are observable within a test
// run; production keeps the real day unit. Both the x-amz-expiration
// projection and the enforcement worker derive their timing from this single
// unit, so the promised expiry and the executed one can never disagree.
var lifecycleDayUnit = 24 * time.Hour

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		lifecycleDayUnit = time.Second
	}
}

// lifecycleDaysInstant returns the instant a Days-based lifecycle action
// occurs, following the AWS timing calculation: the days are added to base —
// the object's creation, or the successor's creation for noncurrent windows —
// and the result is rounded up to the next midnight UTC. Midnights are the
// multiples of the day unit from the epoch, so the compressed test unit
// keeps the same calculation shape.
func lifecycleDaysInstant(base time.Time, days int32) time.Time {
	t := base.UTC().Add(time.Duration(days) * lifecycleDayUnit)
	if rem := t.Sub(t.Truncate(lifecycleDayUnit)); rem > 0 {
		t = t.Add(lifecycleDayUnit - rem)
	}
	return t
}

// objectExpiration projects a bucket's lifecycle configuration onto one
// object: the earliest instant an Enabled expiration rule will expire it,
// with the rule's ID. Only the current version carries a projected expiry,
// and rules whose only expiration member is ExpiredObjectDeleteMarker
// remove markers, not data objects. A Days window follows the AWS timing
// calculation: the days are added to the object's LastModified and the
// result is rounded up to the next midnight UTC; a Date rule expires at
// the configured date. Objects whose replication has not succeeded carry
// no expiry — S3 Lifecycle prevents expiration and transition actions on
// them until replication succeeds.
func objectExpiration(config *s3store.LifecycleConfiguration, obj *s3store.Object) (time.Time, string, bool) {
	if config == nil || !obj.IsLatest || obj.IsDeleteMarker {
		return time.Time{}, "", false
	}
	if !lifecycleReplicationEligible(obj) {
		return time.Time{}, "", false
	}
	var earliest time.Time
	ruleID := ""
	for _, rule := range config.Rules {
		if rule.Status != "Enabled" || rule.Expiration == nil {
			continue
		}
		exp := rule.Expiration
		var expiry time.Time
		switch {
		case exp.Days != nil && *exp.Days > 0:
			expiry = lifecycleDaysInstant(obj.LastModified, *exp.Days)
		case exp.Date != nil:
			expiry = exp.Date.UTC()
		default:
			continue
		}
		if !matchesLifecycleFilter(obj, rule.Filter) {
			continue
		}
		if ruleID == "" || expiry.Before(earliest) {
			earliest, ruleID = expiry, rule.ID
		}
	}
	return earliest, ruleID, ruleID != ""
}

// objectExpirationHeaderValue renders the x-amz-expiration response value
// AWS documents: expiry-date at HTTP-date, with the rule-id URL-encoded.
func objectExpirationHeaderValue(config *s3store.LifecycleConfiguration, obj *s3store.Object) string {
	expiry, ruleID, ok := objectExpiration(config, obj)
	if !ok {
		return ""
	}
	return fmt.Sprintf("expiry-date=%q, rule-id=%q", expiry.Format(http.TimeFormat), url.QueryEscape(ruleID))
}

// objectExpirationHeaderFor renders the x-amz-expiration value for an
// object read from a bucket, looking the bucket's lifecycle configuration
// up through the cross-region bucket finder; a missing bucket yields no
// header (the read planes validate existence before this runs).
func (s *S3Service) objectExpirationHeaderFor(bucket string, obj *s3store.Object) string {
	b, _ := s.s3FindBucket(bucket)
	if b == nil {
		return ""
	}
	return objectExpirationHeaderValue(b.LifecycleConfiguration, obj)
}
