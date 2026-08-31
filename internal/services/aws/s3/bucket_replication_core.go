package s3

import (
	"context"
	"fmt"
	"io"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/ptrutil"
)

// Core functions for the bucket replication operations. The
// BucketOperations methods in bucket_replication.go are thin adapters that
// receive the per-region store bundle and delegate all validation and
// persistence here. The async replication engine entry (replicateObject)
// also lives here: store access for cross-region copy belongs to the Core
// layer, and the goroutine callers in object_put.go and the object Core
// files reach it through this file.

func (s *S3Service) putBucketReplicationCore(stores *s3Stores, input *PutBucketReplicationInput) error {
	if input.ReplicationConfiguration == nil || len(input.ReplicationConfiguration.Rules) == 0 {
		return NewInvalidArgumentError("replication configuration must contain at least one rule")
	}

	if err := validateIAMRoleARN(input.ReplicationConfiguration.Role); err != nil {
		return err
	}

	config := &s3store.ReplicationConfiguration{
		Role: input.ReplicationConfiguration.Role,
	}

	for _, rule := range input.ReplicationConfiguration.Rules {
		if rule.Status != "Enabled" && rule.Status != "Disabled" {
			return NewInvalidArgumentError("rule Status must be Enabled or Disabled")
		}
		if rule.Destination == nil || rule.Destination.Bucket == "" {
			return NewInvalidArgumentError("rule Destination.Bucket is required")
		}
		if err := validateStorageClass(rule.Destination.StorageClass); err != nil {
			return err
		}
		if len(rule.ID) > maxReplicationIDLength {
			return NewInvalidArgumentError(fmt.Sprintf("rule ID exceeds maximum length of %d characters", maxReplicationIDLength))
		}
		if rule.Priority != nil && *rule.Priority < 0 {
			return NewInvalidArgumentError("rule Priority must be non-negative")
		}
		if rule.DeleteMarkerReplication != nil {
			if err := validateReplicationStatus(rule.DeleteMarkerReplication.Status); err != nil {
				return err
			}
		}
		if rule.Filter != nil {
			filterCount := 0
			if rule.Filter.Prefix != "" {
				filterCount++
			}
			if rule.Filter.Tag != nil {
				filterCount++
			}
			if rule.Filter.And != nil {
				filterCount++
			}
			if filterCount > 1 {
				return NewInvalidArgumentError("Filter must contain at most one of Prefix, Tag, or And")
			}
			if rule.Filter.And != nil && len(rule.Filter.And.Tags) == 0 && rule.Filter.And.Prefix == "" {
				return NewInvalidArgumentError("Filter.And must contain at least one Prefix or Tag")
			}
		}

		storeRule := s3store.ReplicationRule{
			ID:          rule.ID,
			Status:      rule.Status,
			Destination: &s3store.ReplicationDestination{Bucket: rule.Destination.Bucket, StorageClass: rule.Destination.StorageClass, Account: rule.Destination.Account},
		}
		if rule.Priority != nil {
			storeRule.Priority = *rule.Priority
		}
		if rule.Filter != nil {
			storeRule.Filter = &s3store.ReplicationFilter{Prefix: rule.Filter.Prefix}
			if rule.Filter.Tag != nil {
				storeRule.Filter.Tag = &s3store.ReplicationTagFilter{Key: rule.Filter.Tag.Key, Value: rule.Filter.Tag.Value}
			}
			if rule.Filter.And != nil {
				storeRule.Filter.AndOperator = &s3store.ReplicationAndOperator{Prefix: rule.Filter.And.Prefix}
				for _, t := range rule.Filter.And.Tags {
					storeRule.Filter.AndOperator.Tags = append(storeRule.Filter.AndOperator.Tags, s3store.ReplicationTagFilter{Key: t.Key, Value: t.Value})
				}
			}
		}
		if rule.DeleteMarkerReplication != nil && rule.DeleteMarkerReplication.Status == "Enabled" {
			storeRule.DeleteMarkerReplication = true
		}

		config.Rules = append(config.Rules, storeRule)
	}

	return stores.buckets.SetReplication(input.Bucket, config)
}

func (s *S3Service) getBucketReplicationCore(stores *s3Stores, input *GetBucketReplicationInput) (*GetBucketReplicationOutput, error) {
	config, err := stores.buckets.GetReplication(input.Bucket)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, ErrNoSuchReplication
	}

	result := &GetBucketReplicationOutput{
		ReplicationConfiguration: &ReplicationConfigurationXML{
			Role: config.Role,
		},
	}

	for _, rule := range config.Rules {
		xmlRule := ReplicationRuleXML{
			ID:       rule.ID,
			Status:   rule.Status,
			Priority: ptrutil.PtrNonZero(rule.Priority),
			Destination: &ReplicationDestinationXML{
				Bucket:       rule.Destination.Bucket,
				StorageClass: rule.Destination.StorageClass,
				Account:      rule.Destination.Account,
			},
		}
		if rule.DeleteMarkerReplication {
			xmlRule.DeleteMarkerReplication = &DeleteMarkerReplicationXML{Status: "Enabled"}
		} else {
			xmlRule.DeleteMarkerReplication = &DeleteMarkerReplicationXML{Status: "Disabled"}
		}
		if rule.Filter != nil {
			xmlRule.Filter = &ReplicationFilterXML{Prefix: rule.Filter.Prefix}
			if rule.Filter.Tag != nil {
				xmlRule.Filter.Tag = &ReplicationTagXML{Key: rule.Filter.Tag.Key, Value: rule.Filter.Tag.Value}
			}
			if rule.Filter.AndOperator != nil {
				xmlRule.Filter.And = &ReplicationAndXML{Prefix: rule.Filter.AndOperator.Prefix}
				for _, t := range rule.Filter.AndOperator.Tags {
					xmlRule.Filter.And.Tags = append(xmlRule.Filter.And.Tags, ReplicationTagXML{Key: t.Key, Value: t.Value})
				}
			}
		}
		result.ReplicationConfiguration.Rules = append(result.ReplicationConfiguration.Rules, xmlRule)
	}

	return result, nil
}

func (s *S3Service) deleteBucketReplicationCore(stores *s3Stores, input *DeleteBucketReplicationInput) error {
	return stores.buckets.SetReplication(input.Bucket, nil)
}

// replicateObject performs cross-region replication for a newly created object.
// It reads the bucket's replication configuration, evaluates each rule's filter
// against the object key, and copies the object to the destination bucket.
// Errors are logged but do not fail the original PUT.
func (s *S3Service) replicateObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, bucket *s3store.Bucket, key string, obj *s3store.Object) {
	if bucket.ReplicationConfiguration == nil || len(bucket.ReplicationConfiguration.Rules) == 0 {
		return
	}

	sourceRegion := reqCtx.GetRegion()

	for _, rule := range bucket.ReplicationConfiguration.Rules {
		if rule.Status != "Enabled" {
			continue
		}
		if !matchReplicationFilter(key, obj.Tags, rule.Filter) {
			continue
		}
		if rule.Destination == nil || rule.Destination.Bucket == "" {
			continue
		}

		destBucketName := bucketNameFromArn(rule.Destination.Bucket)
		if destBucketName == "" {
			continue
		}

		destBucket, destRegion := s.findDestBucket(destBucketName, stores, sourceRegion)
		if destBucket == nil {
			logs.Warn("s3: replication destination bucket not found", logs.String("bucket", destBucketName))
			continue
		}

		destStores, err := s.storeForReplication(reqCtx, destRegion)
		if err != nil {
			logs.Warn("s3: replication store lookup failed", logs.String("destRegion", destRegion), logs.Err(err))
			continue
		}

		reader, srcObj, err := stores.objects.Get(ctx, bucket.Name, key)
		if err != nil {
			logs.Warn("s3: replication source read failed", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(err))
			continue
		}

		if srcObj.SSEMetadata != nil {
			// SSE objects must be buffered because PutEncrypted accepts
			// []byte, not io.Reader.  A streaming PutEncrypted would
			// require a store interface change.
			data, readErr := io.ReadAll(reader)
			reader.Close()
			if readErr != nil {
				logs.Warn("s3: replication source read failed", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(readErr))
				if err := stores.objects.SetReplicationStatus(bucket.Name, key, obj.VersionID, "FAILED"); err != nil {
					logs.Warn("s3: failed to set replication status", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(err))
				}
				continue
			}
			_, err = destStores.objects.PutEncrypted(ctx, destBucketName, key, data, srcObj.ContentType, srcObj.Metadata, srcObj.SSEMetadata, srcObj.StorageClass, nil)
		} else {
			// Non-SSE objects can be streamed directly from source to
			// destination without buffering the entire body.
			_, err = destStores.objects.Put(ctx, destBucketName, key, reader, srcObj.ContentType, srcObj.Metadata)
			reader.Close()
		}
		if err != nil {
			logs.Warn("s3: replication write failed", logs.String("destBucket", destBucketName), logs.String("key", key), logs.Err(err))
			if statusErr := stores.objects.SetReplicationStatus(bucket.Name, key, obj.VersionID, "FAILED"); statusErr != nil {
				logs.Warn("s3: failed to set replication status", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(statusErr))
			}
			continue
		}
		if err := stores.objects.SetReplicationStatus(bucket.Name, key, obj.VersionID, "COMPLETED"); err != nil {
			logs.Warn("s3: failed to set replication status", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(err))
		}
	}
}
