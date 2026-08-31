package s3

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
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
	return o.svc.putBucketReplicationCore(store, input)
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
	return o.svc.getBucketReplicationCore(store, input)
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
	return o.svc.deleteBucketReplicationCore(store, input)
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
