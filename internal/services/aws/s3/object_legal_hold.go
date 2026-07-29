package s3

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// PutObjectLegalHoldInput contains the parameters for setting an object's legal hold.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally specifies a specific version of the object.
// LegalHold specifies the legal hold status to apply.
type PutObjectLegalHoldInput struct {
	Bucket    string
	Key       string
	VersionId string
	LegalHold *LegalHoldInput
}

// LegalHoldInput represents the legal hold status to apply to an object.
// Status must be "ON" or "OFF".
type LegalHoldInput struct {
	Status string `xml:"Status"`
}

// PutObjectLegalHold applies a legal hold to an object.
// The bucket must have Object Lock enabled. Legal hold prevents object deletion or overwriting.
func (o *ObjectOperations) PutObjectLegalHold(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectLegalHoldInput) error {
	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return ErrObjectLockNotEnabled
	}

	if input.LegalHold == nil {
		return NewInvalidArgumentError("LegalHold is required")
	}

	var status s3store.ObjectLockLegalHoldStatus
	switch input.LegalHold.Status {
	case "ON":
		status = s3store.ObjectLockLegalHoldOn
	case "OFF":
		status = s3store.ObjectLockLegalHoldOff
	default:
		return NewInvalidArgumentError(fmt.Sprintf("invalid legal hold status: %s (must be ON or OFF)", input.LegalHold.Status))
	}

	legalHold := &s3store.ObjectLockLegalHold{Status: status}

	return stores.objects.SetObjectLegalHold(ctx, input.Bucket, input.Key, input.VersionId, legalHold)
}

// GetObjectLegalHoldInput contains the parameters for retrieving an object's legal hold.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally specifies a specific version of the object.
type GetObjectLegalHoldInput struct {
	Bucket    string
	Key       string
	VersionId string
}

// GetObjectLegalHoldOutput contains the result of retrieving an object's legal hold status.
type GetObjectLegalHoldOutput struct {
	LegalHold *LegalHoldOutput `xml:"LegalHold"`
}

// LegalHoldOutput represents the legal hold status of an object.
// Status is either "ON" or "OFF".
type LegalHoldOutput struct {
	Status string `xml:"Status"`
}

// GetObjectLegalHold retrieves the legal hold status for an object.
// Returns the current legal hold configuration for the specified object version.
func (o *ObjectOperations) GetObjectLegalHold(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *GetObjectLegalHoldInput) (*GetObjectLegalHoldOutput, error) {
	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if !bucket.ObjectLockEnabled {
		return nil, ErrObjectLockNotEnabled
	}

	legalHold, err := stores.objects.GetObjectLegalHold(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	return &GetObjectLegalHoldOutput{
		LegalHold: &LegalHoldOutput{
			Status: string(legalHold.Status),
		},
	}, nil
}

// PutObjectRetentionInput contains the parameters for setting an object's retention period.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally specifies a specific version of the object.
// Retention specifies the retention mode and retain until date.
type PutObjectRetentionInput struct {
	Bucket    string
	Key       string
	VersionId string
	Retention *RetentionInput
}

// RetentionInput represents the retention configuration for an object.
// Mode specifies either GOVERNANCE or COMPLIANCE retention.
// RetainUntilDate specifies when the retention period expires.
type RetentionInput struct {
	Mode            string     `xml:"Mode"`
	RetainUntilDate *time.Time `xml:"RetainUntilDate,omitempty"`
}

// PutObjectRetention applies a retention period to an object.
// The bucket must have Object Lock enabled. COMPLIANCE mode prevents deletion until retention expires.
func (o *ObjectOperations) PutObjectRetention(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectRetentionInput) error {
	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return ErrObjectLockNotEnabled
	}

	if input.Retention == nil {
		return NewInvalidArgumentError("retention is required")
	}

	retention := &s3store.ObjectLockRetention{
		Mode:            s3store.ObjectLockRetentionMode(input.Retention.Mode),
		RetainUntilDate: input.Retention.RetainUntilDate,
	}

	return stores.objects.SetObjectRetention(ctx, input.Bucket, input.Key, input.VersionId, retention)
}

// GetObjectRetentionInput contains the parameters for retrieving an object's retention configuration.
// Bucket is the name of the S3 bucket.
// Key is the object key within the bucket.
// VersionId optionally specifies a specific version of the object.
type GetObjectRetentionInput struct {
	Bucket    string
	Key       string
	VersionId string
}

// GetObjectRetentionOutput contains the result of retrieving an object's retention configuration.
type GetObjectRetentionOutput struct {
	Retention *RetentionOutput `xml:"Retention"`
}

// RetentionOutput represents the retention configuration of an object.
// Mode is either GOVERNANCE or COMPLIANCE.
// RetainUntilDate is the date when the retention period expires.
type RetentionOutput struct {
	Mode            string     `xml:"Mode"`
	RetainUntilDate *time.Time `xml:"RetainUntilDate,omitempty"`
}

// GetObjectRetention retrieves the retention configuration for an object.
// Returns the retention mode and retain until date for the specified object version.
func (o *ObjectOperations) GetObjectRetention(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *GetObjectRetentionInput) (*GetObjectRetentionOutput, error) {
	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if !bucket.ObjectLockEnabled {
		return nil, ErrObjectLockNotEnabled
	}

	retention, err := stores.objects.GetObjectRetention(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	return &GetObjectRetentionOutput{
		Retention: &RetentionOutput{
			Mode:            string(retention.Mode),
			RetainUntilDate: retention.RetainUntilDate,
		},
	}, nil
}
