package s3

import (
	"time"

	"vorpalstacks/internal/common/bucketname"
	s3store "vorpalstacks/internal/store/aws/s3"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// inventoryOptionalFields is the set of OptionalFields values AWS accepts in
// an inventory configuration, in the canonical report column order.
var inventoryOptionalFields = map[string]bool{
	"Size":                         true,
	"LastModifiedDate":             true,
	"StorageClass":                 true,
	"ETag":                         true,
	"IsMultipartUploaded":          true,
	"ReplicationStatus":            true,
	"EncryptionStatus":             true,
	"ObjectLockRetainUntilDate":    true,
	"ObjectLockMode":               true,
	"ObjectLockLegalHoldStatus":    true,
	"IntelligentTieringAccessTier": true,
	"BucketKeyStatus":              true,
	"ChecksumAlgorithm":            true,
	"ObjectAccessControlList":      true,
	"ObjectOwner":                  true,
	"LifecycleExpirationDate":      true,
}

// putBucketInventoryConfigurationCore validates and stores one inventory
// configuration under the request's id, replacing any existing configuration
// with that id.
func (s *S3Service) putBucketInventoryConfigurationCore(bucketStore s3store.BucketStoreInterface, in *PutBucketInventoryConfigurationInput) error {
	if in.Id == "" {
		return NewInvalidArgumentError("the id query parameter is required")
	}
	config := in.InventoryConfiguration
	if config == nil || config.Id == "" || config.IsEnabled == nil || config.Destination == nil ||
		config.Destination.S3BucketDestination == nil || config.Schedule == nil {
		return ErrMalformedXML
	}
	dest := config.Destination.S3BucketDestination
	if dest.Bucket == "" || dest.Format == "" {
		return ErrMalformedXML
	}
	// The destination bucket field carries an S3 bucket ARN, so a value
	// whose resource part is not a valid bucket name could never receive a
	// delivery — the worker would fail to resolve it on every schedule
	// window. "arn:aws:s3:::bucket/path" parses as an ARN but names no
	// bucket, so the resource part is validated with the shared bucket
	// naming rules.
	if parsedDest, err := svcarn.ParseARN(dest.Bucket); err != nil || parsedDest.Service != "s3" || !bucketname.Validate(parsedDest.Resource) {
		return NewInvalidArgumentError("the destination bucket must be a valid S3 bucket ARN (arn:aws:s3:::bucket-name): " + dest.Bucket)
	}
	switch dest.Format {
	case "CSV", "Parquet", "ORC":
	default:
		return ErrMalformedXML
	}
	if config.IncludedObjectVersions != "All" && config.IncludedObjectVersions != "Current" {
		return ErrMalformedXML
	}
	if config.Schedule.Frequency != "Daily" && config.Schedule.Frequency != "Weekly" {
		return ErrMalformedXML
	}
	for _, field := range config.OptionalFields {
		if !inventoryOptionalFields[field] {
			return ErrMalformedXML
		}
	}

	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return err
	}
	if _, exists := bucket.InventoryConfigurations[in.Id]; !exists &&
		len(bucket.InventoryConfigurations) >= s3store.MaxBucketConfigurations {
		return ErrTooManyConfigurations
	}

	stored := &s3store.InventoryConfiguration{
		ID:                     in.Id,
		IsEnabled:              *config.IsEnabled,
		IncludedObjectVersions: config.IncludedObjectVersions,
		OptionalFields:         config.OptionalFields,
		Schedule:               &s3store.InventorySchedule{Frequency: config.Schedule.Frequency},
		// The report worker treats LastDelivery as the schedule anchor: the
		// configuration time before the first delivery, the delivery time
		// after it.
		LastDelivery: time.Now(),
	}
	if config.Filter != nil {
		stored.Filter = &s3store.InventoryFilter{Prefix: config.Filter.Prefix}
	}
	if dest != nil {
		storedEncryption := &s3store.InventoryEncryption{}
		if dest.Encryption != nil {
			if dest.Encryption.SSES3 != nil {
				storedEncryption.SSES3 = true
			}
			if dest.Encryption.SSEKMS != nil {
				storedEncryption.SSEKMS = &s3store.InventorySSEKMS{KeyID: dest.Encryption.SSEKMS.KeyID}
			}
		}
		stored.Destination = &s3store.InventoryDestination{
			S3BucketDestination: &s3store.InventoryS3BucketDestination{
				AccountID:  dest.AccountID,
				Bucket:     dest.Bucket,
				Format:     dest.Format,
				Prefix:     dest.Prefix,
				Encryption: storedEncryption,
			},
		}
	}
	return bucketStore.SetInventoryConfiguration(in.Bucket, in.Id, stored)
}

// getBucketInventoryConfigurationCore reads one inventory configuration.
func (s *S3Service) getBucketInventoryConfigurationCore(bucketStore s3store.BucketStoreInterface, in *GetBucketInventoryConfigurationInput) (*GetBucketInventoryConfigurationOutput, error) {
	stored, err := bucketStore.GetInventoryConfiguration(in.Bucket, in.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, ErrNoSuchConfiguration
	}
	return &GetBucketInventoryConfigurationOutput{InventoryConfiguration: inventoryConfigurationToOutput(stored)}, nil
}

// deleteBucketInventoryConfigurationCore removes one inventory
// configuration; an absent id deletes nothing, matching the idempotent 204
// the AWS API documents.
func (s *S3Service) deleteBucketInventoryConfigurationCore(bucketStore s3store.BucketStoreInterface, in *DeleteBucketInventoryConfigurationInput) error {
	if in.Id == "" {
		return NewInvalidArgumentError("the id query parameter is required")
	}
	return bucketStore.DeleteInventoryConfiguration(in.Bucket, in.Id)
}

// listBucketInventoryConfigurationsCore returns one page of at most 100
// configurations ordered by id; the continuation token is the id of the last
// configuration returned.
func (s *S3Service) listBucketInventoryConfigurationsCore(bucketStore s3store.BucketStoreInterface, in *ListBucketInventoryConfigurationsInput) (*ListBucketInventoryConfigurationsOutput, error) {
	configs, err := bucketStore.ListInventoryConfigurations(in.Bucket)
	if err != nil {
		return nil, err
	}
	const pageSize = 100
	start := 0
	if in.ContinuationToken != "" {
		for start < len(configs) && configs[start].ID <= in.ContinuationToken {
			start++
		}
	}
	end := start + pageSize
	if end > len(configs) {
		end = len(configs)
	}
	out := &ListBucketInventoryConfigurationsOutput{
		ContinuationToken: in.ContinuationToken,
	}
	for _, stored := range configs[start:end] {
		out.InventoryConfigurations = append(out.InventoryConfigurations, *inventoryConfigurationToOutput(stored))
	}
	if end < len(configs) {
		out.IsTruncated = true
		out.NextContinuationToken = configs[end-1].ID
	}
	return out, nil
}

func inventoryConfigurationToOutput(stored *s3store.InventoryConfiguration) *InventoryConfigurationOutput {
	out := &InventoryConfigurationOutput{
		Id:                     stored.ID,
		IsEnabled:              stored.IsEnabled,
		IncludedObjectVersions: stored.IncludedObjectVersions,
		OptionalFields:         stored.OptionalFields,
	}
	if stored.Filter != nil {
		out.Filter = &InventoryFilterOutput{Prefix: stored.Filter.Prefix}
	}
	if stored.Schedule != nil {
		out.Schedule = &InventoryScheduleOutput{Frequency: stored.Schedule.Frequency}
	}
	if stored.Destination != nil && stored.Destination.S3BucketDestination != nil {
		dest := stored.Destination.S3BucketDestination
		var encryption *InventoryEncryptionOutput
		if dest.Encryption != nil {
			encryption = &InventoryEncryptionOutput{}
			if dest.Encryption.SSES3 {
				encryption.SSES3 = &struct{}{}
			}
			if dest.Encryption.SSEKMS != nil {
				encryption.SSEKMS = &InventorySSEKMSOutput{KeyID: dest.Encryption.SSEKMS.KeyID}
			}
		}
		out.Destination = &InventoryDestinationOutput{
			S3BucketDestination: &InventoryS3BucketDestinationOutput{
				AccountID:  dest.AccountID,
				Bucket:     dest.Bucket,
				Format:     dest.Format,
				Prefix:     dest.Prefix,
				Encryption: encryption,
			},
		}
	}
	return out
}
