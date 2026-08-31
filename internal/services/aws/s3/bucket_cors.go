package s3

import (
	"vorpalstacks/internal/common/request"
)

// PutBucketCORSInput is the input for PutBucketCORS.
type PutBucketCORSInput struct {
	Bucket            string
	CORSConfiguration *CORSConfigurationInput
}

// CORSConfigurationInput defines the CORS rules for a bucket.
type CORSConfigurationInput struct {
	CORSRules []CORSRuleInput `xml:"CORSRule"`
}

// CORSRuleInput defines a CORS rule with allowed methods, origins, and headers.
type CORSRuleInput struct {
	AllowedHeaders []string `xml:"AllowedHeader"`
	AllowedMethods []string `xml:"AllowedMethod"`
	AllowedOrigins []string `xml:"AllowedOrigin"`
	ExposeHeaders  []string `xml:"ExposeHeader,omitempty"`
	MaxAgeSeconds  *int32   `xml:"MaxAgeSeconds,omitempty"`
	ID             string   `xml:"ID,omitempty"`
}

// PutBucketCORS sets the CORS configuration for an S3 bucket.
func (o *BucketOperations) PutBucketCORS(ctx *request.RequestContext, input *PutBucketCORSInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketCORSCore(store.buckets, input)
}

// GetBucketCORSInput is the input for GetBucketCORS.
type GetBucketCORSInput struct {
	Bucket string
}

// GetBucketCORSOutput is the output of GetBucketCORS.
type GetBucketCORSOutput struct {
	CORSConfiguration *CORSConfigurationInput `xml:"CORSConfiguration"`
}

// GetBucketCORS retrieves the CORS configuration for an S3 bucket.
func (o *BucketOperations) GetBucketCORS(ctx *request.RequestContext, input *GetBucketCORSInput) (*GetBucketCORSOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketCORSCore(store.buckets, input)
}

// DeleteBucketCORSInput is the input for DeleteBucketCORS.
type DeleteBucketCORSInput struct {
	Bucket string
}

// DeleteBucketCORS removes the CORS configuration from an S3 bucket.
func (o *BucketOperations) DeleteBucketCORS(ctx *request.RequestContext, input *DeleteBucketCORSInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketCORSCore(store.buckets, input)
}
