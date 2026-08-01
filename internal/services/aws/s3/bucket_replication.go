package s3

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/ptrutil"
)

// PutBucketReplicationInput contains the input for PutBucketReplication.
type PutBucketReplicationInput struct {
	Bucket                   string
	ReplicationConfiguration *ReplicationConfigurationXML
}

// ReplicationConfigurationXML is the XML representation of the replication configuration.
type ReplicationConfigurationXML struct {
	XMLName xml.Name             `xml:"ReplicationConfiguration"`
	Role    string               `xml:"Role"`
	Rules   []ReplicationRuleXML `xml:"Rule"`
}

// ReplicationRuleXML is the XML representation of a replication rule.
type ReplicationRuleXML struct {
	ID                      string                      `xml:"ID,omitempty"`
	Priority                *int32                      `xml:"Priority,omitempty"`
	Status                  string                      `xml:"Status"`
	Filter                  *ReplicationFilterXML       `xml:"Filter,omitempty"`
	Destination             *ReplicationDestinationXML  `xml:"Destination"`
	DeleteMarkerReplication *DeleteMarkerReplicationXML `xml:"DeleteMarkerReplication,omitempty"`
}

// ReplicationFilterXML is the XML representation of a replication filter.
type ReplicationFilterXML struct {
	Prefix string             `xml:"Prefix,omitempty"`
	Tag    *ReplicationTagXML `xml:"Tag,omitempty"`
	And    *ReplicationAndXML `xml:"And,omitempty"`
}

// ReplicationTagXML is the XML representation of a tag filter.
type ReplicationTagXML struct {
	Key   string `xml:"Key,omitempty"`
	Value string `xml:"Value,omitempty"`
}

// ReplicationAndXML is the XML representation of an AND filter.
type ReplicationAndXML struct {
	Prefix string              `xml:"Prefix,omitempty"`
	Tags   []ReplicationTagXML `xml:"Tag,omitempty"`
}

// ReplicationDestinationXML is the XML representation of the replication destination.
type ReplicationDestinationXML struct {
	Bucket       string `xml:"Bucket"`
	StorageClass string `xml:"StorageClass,omitempty"`
	Account      string `xml:"Account,omitempty"`
}

// DeleteMarkerReplicationXML controls whether delete markers are replicated.
type DeleteMarkerReplicationXML struct {
	Status string `xml:"Status"`
}

// PutBucketReplication sets the replication configuration for a bucket.
func (o *BucketOperations) PutBucketReplication(ctx *request.RequestContext, input *PutBucketReplicationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}

	if input.ReplicationConfiguration == nil || len(input.ReplicationConfiguration.Rules) == 0 {
		return NewInvalidArgumentError("replication configuration must contain at least one rule")
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

	return store.buckets.SetReplication(input.Bucket, config)
}

// GetBucketReplicationInput contains the input for GetBucketReplication.
type GetBucketReplicationInput struct {
	Bucket string
}

// GetBucketReplicationOutput contains the result of GetBucketReplication.
type GetBucketReplicationOutput struct {
	ReplicationConfiguration *ReplicationConfigurationXML `xml:"ReplicationConfiguration"`
}

// GetBucketReplication retrieves the replication configuration for a bucket.
func (o *BucketOperations) GetBucketReplication(ctx *request.RequestContext, input *GetBucketReplicationInput) (*GetBucketReplicationOutput, error) {
	store, err := o.svc.store(ctx)
	if err != nil {
		return nil, err
	}

	config, err := store.buckets.GetReplication(input.Bucket)
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

// DeleteBucketReplicationInput contains the input for DeleteBucketReplication.
type DeleteBucketReplicationInput struct {
	Bucket string
}

// DeleteBucketReplication removes the replication configuration from a bucket.
func (o *BucketOperations) DeleteBucketReplication(ctx *request.RequestContext, input *DeleteBucketReplicationInput) error {
	store, err := o.svc.store(ctx)
	if err != nil {
		return err
	}
	return store.buckets.SetReplication(input.Bucket, nil)
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

		data, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			logs.Warn("s3: replication source read failed", logs.String("bucket", bucket.Name), logs.String("key", key), logs.Err(err))
			continue
		}

		if srcObj.SSEMetadata != nil {
			_, err = destStores.objects.PutEncrypted(ctx, destBucketName, key, data, srcObj.ContentType, srcObj.Metadata, srcObj.SSEMetadata, srcObj.StorageClass, nil)
		} else {
			_, err = destStores.objects.Put(ctx, destBucketName, key, bytes.NewReader(data), srcObj.ContentType, srcObj.Metadata)
		}
		if err != nil {
			logs.Warn("s3: replication write failed", logs.String("destBucket", destBucketName), logs.String("key", key), logs.Err(err))
			continue
		}
	}
}

// replicateDeleteMarker propagates a delete marker to replication destination
// buckets when DeleteMarkerReplication is Enabled for a matching rule.
// Errors are logged but do not fail the original DELETE.
func (s *S3Service) replicateDeleteMarker(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, bucket *s3store.Bucket, key string) {
	if bucket.ReplicationConfiguration == nil || len(bucket.ReplicationConfiguration.Rules) == 0 {
		return
	}

	sourceRegion := reqCtx.GetRegion()

	for _, rule := range bucket.ReplicationConfiguration.Rules {
		if rule.Status != "Enabled" || !rule.DeleteMarkerReplication {
			continue
		}
		if !matchReplicationFilter(key, nil, rule.Filter) {
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
			logs.Warn("s3: delete-marker replication destination bucket not found", logs.String("bucket", destBucketName))
			continue
		}

		destStores, err := s.storeForReplication(reqCtx, destRegion)
		if err != nil {
			logs.Warn("s3: delete-marker replication store lookup failed", logs.String("destRegion", destRegion), logs.Err(err))
			continue
		}

		marker, err := destStores.objects.DeleteWithVersion(ctx, destBucketName, key, "")
		if err != nil {
			logs.Warn("s3: delete-marker replication write failed", logs.String("destBucket", destBucketName), logs.String("key", key), logs.Err(err))
			continue
		}
		_ = marker
	}
}

// matchReplicationFilter evaluates whether an object key matches the replication filter.
func matchReplicationFilter(key string, tags []types.Tag, filter *s3store.ReplicationFilter) bool {
	if filter == nil {
		return true
	}

	// When AndOperator is present, all conditions (prefix + tags) must match.
	// When AndOperator is absent, Filter-level Prefix and Tag are evaluated
	// as independent conditions (both must match if both are set).
	if filter.AndOperator != nil {
		if filter.AndOperator.Prefix != "" && !startsWith(key, filter.AndOperator.Prefix) {
			return false
		}
		for _, filterTag := range filter.AndOperator.Tags {
			if !objectHasTag(tags, filterTag.Key, filterTag.Value) {
				return false
			}
		}
		return true
	}

	if filter.Prefix != "" && !startsWith(key, filter.Prefix) {
		return false
	}
	if filter.Tag != nil {
		if !objectHasTag(tags, filter.Tag.Key, filter.Tag.Value) {
			return false
		}
	}
	return true
}

// objectHasTag checks whether the object's tag set contains a tag with the
// given key and value.
func objectHasTag(tags []types.Tag, key, value string) bool {
	for _, t := range tags {
		if t.Key == key && t.Value == value {
			return true
		}
	}
	return false
}

func startsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// bucketNameFromArn extracts the bucket name from an S3 ARN or bucket path.
// Accepts "arn:aws:s3:::bucket-name", "bucket-name", or "/bucket-name".
func bucketNameFromArn(bucketArn string) string {
	bucketArn = strings.TrimPrefix(bucketArn, "/")
	bucketArn = strings.TrimPrefix(bucketArn, "arn:aws:s3:::")
	bucketArn = strings.TrimPrefix(bucketArn, "arn:aws-cn:s3:::")
	bucketArn = strings.TrimPrefix(bucketArn, "arn:aws-us-gov:s3:::")
	return bucketArn
}

// findDestBucket resolves a destination bucket by name, searching across all
// regions when the multi-region store is available. Falls back to the source
// region store when multi-region is not configured. Returns the bucket
// metadata and the region in which it was found.
func (s *S3Service) findDestBucket(name string, sourceStores *s3Stores, sourceRegion string) (*s3store.Bucket, string) {
	if s.s3Store != nil {
		if bucket, region := s.s3Store.FindBucket(name); bucket != nil {
			return bucket, region
		}
		return nil, ""
	}
	bucket, err := sourceStores.buckets.Get(name)
	if err != nil {
		return nil, ""
	}
	region := bucket.Region
	if region == "" {
		region = sourceRegion
	}
	return bucket, region
}

// storeForReplication returns the stores for a destination region.
func (s *S3Service) storeForReplication(reqCtx *request.RequestContext, region string) (*s3Stores, error) {
	if s.s3Store != nil {
		buckets := s.s3Store.Buckets(region)
		objects := s.s3Store.Objects(region)
		if buckets == nil || objects == nil {
			return nil, fmt.Errorf("failed to get stores for region %s", region)
		}
		return &s3Stores{buckets: buckets, objects: objects}, nil
	}
	return nil, fmt.Errorf("multi-region store not available")
}
