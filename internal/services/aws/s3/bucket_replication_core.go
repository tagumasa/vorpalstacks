package s3

import (
	"context"
	"fmt"
	"io"
	"time"

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

// replicationCallTimeout bounds one async replication goroutine so a slow
// cross-region copy can never leak a goroutine forever.
const replicationCallTimeout = 30 * time.Second

func (s *S3Service) putBucketReplicationCore(stores *s3Stores, input *PutBucketReplicationInput) error {
	if input.ReplicationConfiguration == nil || len(input.ReplicationConfiguration.Rules) == 0 {
		return NewInvalidArgumentError("replication configuration must contain at least one rule")
	}

	if err := validateIAMRoleARN(input.ReplicationConfiguration.Role); err != nil {
		return err
	}

	// Replication requires versioning on both sides: the source bucket must
	// have versioning Enabled, and so must every destination bucket.
	sourceBucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}
	if sourceBucket.VersioningStatus != s3store.BucketVersioningEnabled {
		return NewInvalidRequestError("Versioning must be 'Enabled' on the bucket to apply a replication configuration")
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
		destBucketName := bucketNameFromArn(rule.Destination.Bucket)
		if destBucketName == "" {
			return NewInvalidArgumentError("rule Destination.Bucket is required")
		}
		destBucket, _ := s.findDestBucket(destBucketName, stores, "")
		if destBucket == nil {
			return NewInvalidRequestError("Destination bucket must exist")
		}
		if destBucket.VersioningStatus != s3store.BucketVersioningEnabled {
			return NewInvalidRequestError("Destination bucket must have versioning enabled")
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

// launchObjectReplication starts the async cross-region replication for a
// newly created object version (PutObject, CopyObject, and
// CompleteMultipartUpload all reach this entry). It resolves the source
// bucket metadata and returns immediately when no replication configuration
// is set, so the common no-replication path pays only one bucket read.
// Objects that are themselves replicas are never replicated onwards. An
// object matching at least one rule is marked PENDING synchronously — the
// status must be observable from the moment the creating operation returns —
// and the copy then runs on a bounded goroutine, so replication never
// delays or fails the object-creating operation beyond that status write.
func (s *S3Service) launchObjectReplication(reqCtx *request.RequestContext, stores *s3Stores, bucketName, key string, obj *s3store.Object) {
	bucket, err := stores.buckets.Get(bucketName)
	if err != nil {
		return
	}
	if bucket.ReplicationConfiguration == nil || len(bucket.ReplicationConfiguration.Rules) == 0 {
		return
	}
	if obj.ReplicationStatus == s3store.ReplicationStatusReplica {
		return
	}

	rules := matchingReplicationRules(bucket, key, obj)
	if len(rules) == 0 {
		return
	}

	if err := stores.objects.SetReplicationStatus(bucketName, key, obj.VersionID, s3store.ReplicationStatusPending); err != nil {
		logs.Warn("s3: failed to set replication status", logs.String("bucket", bucketName), logs.String("key", key), logs.Err(err))
	}

	s.goReplicationWorker(bucketName, key, "object", func(ctx context.Context) {
		s.replicateObject(ctx, reqCtx, stores, bucket, key, obj, rules)
	})
}

// goReplicationWorker runs fn on a goroutine bounded by
// replicationCallTimeout and guarded by panic recovery. The object and
// delete-marker replication launches share it so a slow or panicking copy
// can never leak the goroutine or fail the triggering operation.
func (s *S3Service) goReplicationWorker(bucketName, key, worker string, fn func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), replicationCallTimeout)
	go func() {
		defer cancel()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("s3: replication goroutine panic",
					logs.String("bucket", bucketName),
					logs.String("key", key),
					logs.String("worker", worker),
					logs.Any("panic", r))
			}
		}()
		fn(ctx)
	}()
}

// matchingReplicationRules returns the enabled rules whose filters match the
// object and whose destinations are usable. The same predicate gates the
// synchronous PENDING marking and the async copy, so an object is never
// marked for replication without a rule to replicate it under.
func matchingReplicationRules(bucket *s3store.Bucket, key string, obj *s3store.Object) []s3store.ReplicationRule {
	var matched []s3store.ReplicationRule
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
		if bucketNameFromArn(rule.Destination.Bucket) == "" {
			continue
		}
		matched = append(matched, rule)
	}
	return matched
}

// replicateObject performs cross-region replication for a newly created
// object version under the rules that matched it at launch time. Each
// destination is re-validated at replication time — the destination bucket
// must still exist and keep versioning enabled — and the source status is
// written exactly once after the whole pass: COMPLETED only when every
// matched rule succeeded, FAILED as soon as one failed. Errors are logged
// but do not fail the original PUT.
func (s *S3Service) replicateObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, bucket *s3store.Bucket, key string, obj *s3store.Object, rules []s3store.ReplicationRule) {
	sourceRegion := reqCtx.GetRegion()
	failed := 0

	for _, rule := range rules {
		destBucketName := bucketNameFromArn(rule.Destination.Bucket)
		destBucket, destRegion := s.findDestBucket(destBucketName, stores, sourceRegion)
		if destBucket == nil {
			logs.Warn("s3: replication destination bucket not found", logs.String("bucket", destBucketName))
			failed++
			continue
		}
		// Destination versioning is required at replication time, not only
		// when the configuration was written: if it has been suspended
		// since, the rule fails instead of degrading the replica to a null
		// version on the suspended bucket.
		if destBucket.VersioningStatus != s3store.BucketVersioningEnabled {
			logs.Warn("s3: replication destination versioning is not enabled", logs.String("bucket", destBucketName))
			failed++
			continue
		}

		destStores, err := s.storeForReplication(reqCtx, destRegion)
		if err != nil {
			logs.Warn("s3: replication store lookup failed", logs.String("destRegion", destRegion), logs.Err(err))
			failed++
			continue
		}

		// The read is pinned to the version that triggered replication so a
		// newer version uploaded in the interim is never copied under this
		// version's replication status. A fresh reader is opened per rule
		// so a streaming copy never consumes another rule's data.
		reader, srcObj, err := stores.objects.GetWithVersion(ctx, bucket.Name, key, obj.VersionID)
		if err != nil {
			logs.Warn("s3: replication source read failed", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(err))
			failed++
			continue
		}

		// The replica inherits the source object's storage class unless the
		// rule's Destination.StorageClass overrides it.
		storageClass := srcObj.StorageClass
		if rule.Destination.StorageClass != "" {
			storageClass = s3store.ObjectStorageClass(rule.Destination.StorageClass)
		}
		if storageClass == "" {
			storageClass = s3store.StorageClassStandard
		}

		var replica *s3store.Object
		if srcObj.SSEMetadata != nil {
			// SSE objects must be buffered because PutEncrypted accepts
			// []byte, not io.Reader.  A streaming PutEncrypted would
			// require a store interface change.
			data, readErr := io.ReadAll(reader)
			reader.Close()
			if readErr != nil {
				logs.Warn("s3: replication source read failed", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(readErr))
				failed++
				continue
			}
			replica, err = destStores.objects.PutEncrypted(ctx, destBucketName, key, data, srcObj.ContentType, srcObj.Metadata, srcObj.SSEMetadata, storageClass, nil)
		} else {
			// Non-SSE objects can be streamed directly from source to
			// destination without buffering the entire body.
			replica, err = destStores.objects.PutWithVersioning(ctx, destBucketName, key, reader, srcObj.ContentType, srcObj.Metadata, false, storageClass, nil)
			reader.Close()
		}
		if err != nil {
			logs.Warn("s3: replication write failed", logs.String("destBucket", destBucketName), logs.String("key", key), logs.Err(err))
			failed++
			continue
		}
		s.applyReplicaMetadata(ctx, destStores, destBucketName, key, replica, srcObj)
	}

	status := s3store.ReplicationStatusCompleted
	if failed > 0 {
		status = s3store.ReplicationStatusFailed
	}
	if err := stores.objects.SetReplicationStatus(bucket.Name, key, obj.VersionID, status); err != nil {
		logs.Warn("s3: failed to set replication status", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(err))
	}
}

// applyReplicaMetadata carries the replicated metadata that the destination
// write path does not take as arguments: the source tag set, ACL, Object
// Lock retention and legal hold, and the REPLICA replication status that
// distinguishes a destination copy from an original upload. Failures are
// logged but never demote the replication result — the object data itself
// has already been copied.
func (s *S3Service) applyReplicaMetadata(ctx context.Context, destStores *s3Stores, destBucket, key string, replica *s3store.Object, srcObj *s3store.Object) {
	versionId := ""
	if replica != nil {
		versionId = replica.VersionID
	}
	if len(srcObj.Tags) > 0 {
		if err := destStores.objects.SetTags(destBucket, key, versionId, srcObj.Tags); err != nil {
			logs.Warn("s3: replica tag write failed", logs.String("destBucket", destBucket), logs.String("key", key), logs.Err(err))
		}
	}
	if srcObj.ACL != nil {
		if err := destStores.objects.SetACLWithVersion(destBucket, key, versionId, srcObj.ACL); err != nil {
			logs.Warn("s3: replica ACL write failed", logs.String("destBucket", destBucket), logs.String("key", key), logs.Err(err))
		}
	}
	if srcObj.ObjectLockRetention != nil {
		if err := destStores.objects.SetObjectRetention(ctx, destBucket, key, versionId, srcObj.ObjectLockRetention); err != nil {
			logs.Warn("s3: replica object-lock retention write failed", logs.String("destBucket", destBucket), logs.String("key", key), logs.Err(err))
		}
	}
	if srcObj.ObjectLockLegalHold != nil {
		if err := destStores.objects.SetObjectLegalHold(ctx, destBucket, key, versionId, srcObj.ObjectLockLegalHold); err != nil {
			logs.Warn("s3: replica object-lock legal-hold write failed", logs.String("destBucket", destBucket), logs.String("key", key), logs.Err(err))
		}
	}
	if err := destStores.objects.SetReplicationStatus(destBucket, key, versionId, s3store.ReplicationStatusReplica); err != nil {
		logs.Warn("s3: failed to mark replica", logs.String("destBucket", destBucket), logs.String("key", key), logs.Err(err))
	}
}
