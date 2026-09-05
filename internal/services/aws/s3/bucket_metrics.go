package s3

import (
	"vorpalstacks/internal/common/request"
)

// PutBucketMetricsConfigurationInput is the input for PutBucketMetricsConfiguration.
type PutBucketMetricsConfigurationInput struct {
	Bucket               string
	Id                   string
	MetricsConfiguration *MetricsConfigurationInput
}

// MetricsConfigurationInput is the wire form of a metrics configuration.
type MetricsConfigurationInput struct {
	Id     string              `xml:"Id"`
	Filter *MetricsFilterInput `xml:"Filter"`
}

// MetricsFilterInput is the wire form of a metrics filter.
type MetricsFilterInput struct {
	Prefix         string                   `xml:"Prefix"`
	Tag            *Tag                     `xml:"Tag"`
	AccessPointArn string                   `xml:"AccessPointArn"`
	And            *MetricsAndOperatorInput `xml:"And"`
}

// MetricsAndOperatorInput is the wire form of a metrics filter conjunction.
type MetricsAndOperatorInput struct {
	Prefix         string `xml:"Prefix"`
	Tags           []Tag  `xml:"Tag"`
	AccessPointArn string `xml:"AccessPointArn"`
}

// PutBucketMetricsConfiguration stores a metrics configuration on a bucket.
func (o *BucketOperations) PutBucketMetricsConfiguration(ctx *request.RequestContext, input *PutBucketMetricsConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketMetricsConfigurationCore(store.buckets, input)
}

// GetBucketMetricsConfigurationInput is the input for GetBucketMetricsConfiguration.
type GetBucketMetricsConfigurationInput struct {
	Bucket string
	Id     string
}

// GetBucketMetricsConfigurationOutput is the output of GetBucketMetricsConfiguration.
type GetBucketMetricsConfigurationOutput struct {
	MetricsConfiguration *MetricsConfigurationOutput `xml:"MetricsConfiguration"`
}

// GetBucketMetricsConfiguration retrieves one metrics configuration.
func (o *BucketOperations) GetBucketMetricsConfiguration(ctx *request.RequestContext, input *GetBucketMetricsConfigurationInput) (*GetBucketMetricsConfigurationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketMetricsConfigurationCore(store.buckets, input)
}

// DeleteBucketMetricsConfigurationInput is the input for DeleteBucketMetricsConfiguration.
type DeleteBucketMetricsConfigurationInput struct {
	Bucket string
	Id     string
}

// DeleteBucketMetricsConfiguration removes one metrics configuration.
func (o *BucketOperations) DeleteBucketMetricsConfiguration(ctx *request.RequestContext, input *DeleteBucketMetricsConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketMetricsConfigurationCore(store.buckets, input)
}

// ListBucketMetricsConfigurationsInput is the input for ListBucketMetricsConfigurations.
type ListBucketMetricsConfigurationsInput struct {
	Bucket            string
	ContinuationToken string
}

// ListBucketMetricsConfigurationsOutput is the output of ListBucketMetricsConfigurations.
type ListBucketMetricsConfigurationsOutput struct {
	MetricsConfigurations []MetricsConfigurationOutput `xml:"MetricsConfiguration"`
	ContinuationToken     string                       `xml:"ContinuationToken,omitempty"`
	IsTruncated           bool                         `xml:"IsTruncated"`
	NextContinuationToken string                       `xml:"NextContinuationToken,omitempty"`
}

// ListBucketMetricsConfigurations lists a bucket's metrics configurations.
func (o *BucketOperations) ListBucketMetricsConfigurations(ctx *request.RequestContext, input *ListBucketMetricsConfigurationsInput) (*ListBucketMetricsConfigurationsOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.listBucketMetricsConfigurationsCore(store.buckets, input)
}

// MetricsConfigurationOutput is the wire form of a stored metrics
// configuration.
type MetricsConfigurationOutput struct {
	Id     string               `xml:"Id"`
	Filter *MetricsFilterOutput `xml:"Filter,omitempty"`
}

// MetricsFilterOutput is the wire form of a metrics filter.
type MetricsFilterOutput struct {
	Prefix         string                    `xml:"Prefix,omitempty"`
	Tag            *Tag                      `xml:"Tag,omitempty"`
	AccessPointArn string                    `xml:"AccessPointArn,omitempty"`
	And            *MetricsAndOperatorOutput `xml:"And,omitempty"`
}

// MetricsAndOperatorOutput is the wire form of a metrics filter conjunction.
type MetricsAndOperatorOutput struct {
	Prefix         string `xml:"Prefix,omitempty"`
	Tags           []Tag  `xml:"Tag,omitempty"`
	AccessPointArn string `xml:"AccessPointArn,omitempty"`
}
