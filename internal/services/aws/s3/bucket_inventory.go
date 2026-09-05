package s3

import (
	"vorpalstacks/internal/common/request"
)

// PutBucketInventoryConfigurationInput is the input for PutBucketInventoryConfiguration.
type PutBucketInventoryConfigurationInput struct {
	Bucket                 string
	Id                     string
	InventoryConfiguration *InventoryConfigurationInput
}

// InventoryConfigurationInput is the wire form of an inventory configuration.
type InventoryConfigurationInput struct {
	Id                     string                     `xml:"Id"`
	Destination            *InventoryDestinationInput `xml:"Destination"`
	IsEnabled              *bool                      `xml:"IsEnabled"`
	Filter                 *InventoryFilterInput      `xml:"Filter"`
	IncludedObjectVersions string                     `xml:"IncludedObjectVersions"`
	OptionalFields         []string                   `xml:"OptionalFields>Field"`
	Schedule               *InventoryScheduleInput    `xml:"Schedule"`
}

// InventoryDestinationInput is the wire form of an inventory destination.
type InventoryDestinationInput struct {
	S3BucketDestination *InventoryS3BucketDestinationInput `xml:"S3BucketDestination"`
}

// InventoryS3BucketDestinationInput is the wire form of the S3 report destination.
type InventoryS3BucketDestinationInput struct {
	AccountID  string                    `xml:"AccountId"`
	Bucket     string                    `xml:"Bucket"`
	Format     string                    `xml:"Format"`
	Prefix     string                    `xml:"Prefix"`
	Encryption *InventoryEncryptionInput `xml:"Encryption"`
}

// InventoryEncryptionInput is the wire form of the report-file encryption.
// Exactly one of SSE-S3 / SSE-KMS is set.
type InventoryEncryptionInput struct {
	SSES3  *struct{}             `xml:"SSE-S3"`
	SSEKMS *InventorySSEKMSInput `xml:"SSE-KMS"`
}

// InventorySSEKMSInput is the wire form of the SSE-KMS report encryption.
type InventorySSEKMSInput struct {
	KeyID string `xml:"KeyId"`
}

// InventoryFilterInput is the wire form of an inventory filter.
type InventoryFilterInput struct {
	Prefix string `xml:"Prefix"`
}

// InventoryScheduleInput is the wire form of an inventory schedule.
type InventoryScheduleInput struct {
	Frequency string `xml:"Frequency"`
}

// PutBucketInventoryConfiguration stores an inventory configuration on a bucket.
func (o *BucketOperations) PutBucketInventoryConfiguration(ctx *request.RequestContext, input *PutBucketInventoryConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.putBucketInventoryConfigurationCore(store.buckets, input)
}

// GetBucketInventoryConfigurationInput is the input for GetBucketInventoryConfiguration.
type GetBucketInventoryConfigurationInput struct {
	Bucket string
	Id     string
}

// GetBucketInventoryConfigurationOutput is the output of GetBucketInventoryConfiguration.
type GetBucketInventoryConfigurationOutput struct {
	InventoryConfiguration *InventoryConfigurationOutput `xml:"InventoryConfiguration"`
}

// GetBucketInventoryConfiguration retrieves one inventory configuration.
func (o *BucketOperations) GetBucketInventoryConfiguration(ctx *request.RequestContext, input *GetBucketInventoryConfigurationInput) (*GetBucketInventoryConfigurationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.getBucketInventoryConfigurationCore(store.buckets, input)
}

// DeleteBucketInventoryConfigurationInput is the input for DeleteBucketInventoryConfiguration.
type DeleteBucketInventoryConfigurationInput struct {
	Bucket string
	Id     string
}

// DeleteBucketInventoryConfiguration removes one inventory configuration.
func (o *BucketOperations) DeleteBucketInventoryConfiguration(ctx *request.RequestContext, input *DeleteBucketInventoryConfigurationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return o.svc.deleteBucketInventoryConfigurationCore(store.buckets, input)
}

// ListBucketInventoryConfigurationsInput is the input for ListBucketInventoryConfigurations.
type ListBucketInventoryConfigurationsInput struct {
	Bucket            string
	ContinuationToken string
}

// ListBucketInventoryConfigurationsOutput is the output of ListBucketInventoryConfigurations.
type ListBucketInventoryConfigurationsOutput struct {
	InventoryConfigurations []InventoryConfigurationOutput `xml:"InventoryConfiguration"`
	ContinuationToken       string                         `xml:"ContinuationToken,omitempty"`
	IsTruncated             bool                           `xml:"IsTruncated"`
	NextContinuationToken   string                         `xml:"NextContinuationToken,omitempty"`
}

// ListBucketInventoryConfigurations lists a bucket's inventory configurations.
func (o *BucketOperations) ListBucketInventoryConfigurations(ctx *request.RequestContext, input *ListBucketInventoryConfigurationsInput) (*ListBucketInventoryConfigurationsOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}
	return o.svc.listBucketInventoryConfigurationsCore(store.buckets, input)
}

// InventoryConfigurationOutput is the wire form of a stored inventory
// configuration. The conversion lives with the Core functions.
type InventoryConfigurationOutput struct {
	Id                     string                      `xml:"Id"`
	Destination            *InventoryDestinationOutput `xml:"Destination"`
	IsEnabled              bool                        `xml:"IsEnabled"`
	Filter                 *InventoryFilterOutput      `xml:"Filter,omitempty"`
	IncludedObjectVersions string                      `xml:"IncludedObjectVersions"`
	OptionalFields         []string                    `xml:"OptionalFields>Field,omitempty"`
	Schedule               *InventoryScheduleOutput    `xml:"Schedule"`
}

// InventoryDestinationOutput is the wire form of an inventory destination.
type InventoryDestinationOutput struct {
	S3BucketDestination *InventoryS3BucketDestinationOutput `xml:"S3BucketDestination"`
}

// InventoryS3BucketDestinationOutput is the wire form of the S3 report destination.
type InventoryS3BucketDestinationOutput struct {
	AccountID  string                     `xml:"AccountId,omitempty"`
	Bucket     string                     `xml:"Bucket"`
	Format     string                     `xml:"Format"`
	Prefix     string                     `xml:"Prefix,omitempty"`
	Encryption *InventoryEncryptionOutput `xml:"Encryption,omitempty"`
}

// InventoryEncryptionOutput is the wire form of the report-file encryption.
type InventoryEncryptionOutput struct {
	SSES3  *struct{}              `xml:"SSE-S3,omitempty"`
	SSEKMS *InventorySSEKMSOutput `xml:"SSE-KMS,omitempty"`
}

// InventorySSEKMSOutput is the wire form of the SSE-KMS report encryption.
type InventorySSEKMSOutput struct {
	KeyID string `xml:"KeyId"`
}

// InventoryFilterOutput is the wire form of an inventory filter.
type InventoryFilterOutput struct {
	Prefix string `xml:"Prefix"`
}

// InventoryScheduleOutput is the wire form of an inventory schedule.
type InventoryScheduleOutput struct {
	Frequency string `xml:"Frequency"`
}
