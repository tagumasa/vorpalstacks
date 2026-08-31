package s3

import (
	"vorpalstacks/internal/common/request"
)

// RequestPaymentConfigurationInput defines the request payment configuration for a bucket.
type RequestPaymentConfigurationInput struct {
	Payer string `xml:"Payer"`
}

// PutBucketRequestPaymentInput represents the input for setting bucket request payment configuration.
type PutBucketRequestPaymentInput struct {
	Bucket                      string
	RequestPaymentConfiguration *RequestPaymentConfigurationInput
}

// PutBucketRequestPayment configures request payment for an S3 bucket.
func (o *BucketOperations) PutBucketRequestPayment(ctx *request.RequestContext, input *PutBucketRequestPaymentInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketRequestPaymentCore(store.buckets, input)
}

// GetBucketRequestPaymentInput represents the input for retrieving bucket request payment configuration.
type GetBucketRequestPaymentInput struct {
	Bucket string
}

// GetBucketRequestPaymentOutput represents the output from retrieving bucket request payment configuration.
type GetBucketRequestPaymentOutput struct {
	RequestPaymentConfiguration *RequestPaymentConfigurationInput `xml:"RequestPaymentConfiguration"`
}

// GetBucketRequestPayment retrieves request payment configuration for an S3 bucket.
func (o *BucketOperations) GetBucketRequestPayment(ctx *request.RequestContext, input *GetBucketRequestPaymentInput) (*GetBucketRequestPaymentOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketRequestPaymentCore(store.buckets, input)
}
