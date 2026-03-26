package s3

import (
	"fmt"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
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
	MaxAgeSeconds  *int     `xml:"MaxAgeSeconds,omitempty"`
	ID             string   `xml:"ID,omitempty"`
}

// PutBucketCORS sets the CORS configuration for an S3 bucket.
func (o *BucketOperations) PutBucketCORS(ctx *request.RequestContext, input *PutBucketCORSInput) error {
	var rules []s3store.CORSRule
	for _, r := range input.CORSConfiguration.CORSRules {
		if len(r.AllowedMethods) == 0 {
			return fmt.Errorf("CORS rule must have at least one AllowedMethod")
		}
		if len(r.AllowedOrigins) == 0 {
			return fmt.Errorf("CORS rule must have at least one AllowedOrigin")
		}
		for _, method := range r.AllowedMethods {
			if !validCORSMethods[method] {
				return fmt.Errorf("invalid CORS method: %s (must be GET, PUT, HEAD, POST, or DELETE)", method)
			}
		}
		rules = append(rules, s3store.CORSRule{
			AllowedHeaders: r.AllowedHeaders,
			AllowedMethods: r.AllowedMethods,
			AllowedOrigins: r.AllowedOrigins,
			ExposeHeaders:  r.ExposeHeaders,
			MaxAgeSeconds:  r.MaxAgeSeconds,
			ID:             r.ID,
		})
	}

	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetCORS(input.Bucket, &s3store.CORSConfiguration{
		CORSRules: rules,
	})
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
	bucket, err := store.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if bucket.CORSConfiguration == nil {
		return nil, ErrNoSuchCORS
	}

	var rules []CORSRuleInput
	for _, r := range bucket.CORSConfiguration.CORSRules {
		rules = append(rules, CORSRuleInput{
			AllowedHeaders: r.AllowedHeaders,
			AllowedMethods: r.AllowedMethods,
			AllowedOrigins: r.AllowedOrigins,
			ExposeHeaders:  r.ExposeHeaders,
			MaxAgeSeconds:  r.MaxAgeSeconds,
			ID:             r.ID,
		})
	}

	return &GetBucketCORSOutput{
		CORSConfiguration: &CORSConfigurationInput{
			CORSRules: rules,
		},
	}, nil
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
	return store.buckets.SetCORS(input.Bucket, nil)
}
