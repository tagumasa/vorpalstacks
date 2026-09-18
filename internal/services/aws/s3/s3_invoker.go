package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"vorpalstacks/internal/common/invokers"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// GetObject implements the invokers.S3Invoker interface. It retrieves the
// content of an object by region, bucket and key, returning the full byte
// content. The object store resolves to the bucket's owning region, so a
// read addressed to a region that does not hold the bucket still finds
// the objects the write path routed there. An object larger than maxBytes
// errors rather than returning a silent prefix — every bounded consumer
// parses the whole object, for which a truncated read is corruption.
func (s *S3Service) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	objs := s.objectStoreForBucket(region, bucket)
	reader, _, err := objs.Get(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject %s/%s: %w", bucket, key, err)
	}
	defer reader.Close()
	data, err := readBounded(reader, maxBytes)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject %s/%s: %w", bucket, key, err)
	}
	return data, nil
}

// GetObjectVersion implements the invokers.S3Invoker interface. It
// retrieves the content of a specific object version; an empty versionID
// reads the latest version, matching the store's version-aware read. As
// with GetObject, an object larger than maxBytes errors rather than
// returning a silent prefix.
func (s *S3Service) GetObjectVersion(ctx context.Context, region, bucket, key, versionID string, maxBytes int64) ([]byte, error) {
	objs := s.objectStoreForBucket(region, bucket)
	reader, _, err := objs.GetWithVersion(ctx, bucket, key, versionID)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObjectVersion %s/%s@%s: %w", bucket, key, versionID, err)
	}
	defer reader.Close()
	data, err := readBounded(reader, maxBytes)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObjectVersion %s/%s@%s: %w", bucket, key, versionID, err)
	}
	return data, nil
}

// readBounded reads the whole object under the maxBytes bound (zero or
// negative meaning the single-upload ceiling), erroring when the object
// exceeds it instead of truncating: the bounded invoker reads are
// whole-object parses, and a silent prefix would surface as a corrupt
// parse far from its cause.
func readBounded(reader io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		maxBytes = maxSingleUploadSize
	}
	data, err := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("object exceeds the %d-byte read bound", maxBytes)
	}
	return data, nil
}

// objectStoreForBucket resolves the object store that owns the named
// bucket. Bucket names are unique across the platform's regions
// (CreateBucket rejects a name another region already owns), so an
// operation addressed to a region that does not hold the bucket — a
// multi-region CloudTrail trail delivering or reading another region's
// files from its one bucket, for example — lands in the bucket's own
// region, where the bucket's readers and writers meet. A bucket no
// initialised region knows keeps the addressed region's store, preserving
// the plain per-region behaviour.
func (s *S3Service) objectStoreForBucket(region, bucket string) s3store.ObjectStoreInterface {
	if _, owner := s.s3Store.FindBucket(bucket); owner != "" && owner != region {
		return s.s3Store.Objects(owner)
	}
	return s.s3Store.Objects(region)
}

// bucketStoreForBucket resolves the bucket-record store that owns the
// named bucket, mirroring objectStoreForBucket for the bucket-level reads
// (existence, policy): a foreign-region bucket's records live in its
// owner's store, so a policy read addressed elsewhere must follow them or
// the bucket looks nonexistent.
func (s *S3Service) bucketStoreForBucket(region, bucket string) s3store.BucketStoreInterface {
	if _, owner := s.s3Store.FindBucket(bucket); owner != "" && owner != region {
		return s.s3Store.Buckets(owner)
	}
	return s.s3Store.Buckets(region)
}

// PutObject stores an object in S3 via the cross-service invoker.
func (s *S3Service) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	objs := s.objectStoreForBucket(region, bucket)
	_, err := objs.Put(ctx, bucket, key, bytes.NewReader(data), contentType, nil)
	if err != nil {
		return fmt.Errorf("s3 PutObject %s/%s: %w", bucket, key, err)
	}
	return nil
}

// PutObjectWithMetadata stores an object carrying S3 object metadata
// (x-amz-meta-*), e.g. CloudTrail digest files whose signature travels as
// object metadata.
func (s *S3Service) PutObjectWithMetadata(ctx context.Context, region, bucket, key string, data []byte, contentType string, metadata map[string]string) error {
	objs := s.objectStoreForBucket(region, bucket)
	_, err := objs.Put(ctx, bucket, key, bytes.NewReader(data), contentType, metadata)
	if err != nil {
		return fmt.Errorf("s3 PutObjectWithMetadata %s/%s: %w", bucket, key, err)
	}
	return nil
}

// maxListAllKeys is the safety cap when listing all objects (maxKeys <= 0).
// This prevents unbounded memory consumption on extremely large buckets
// while accommodating legitimate bulk-load workloads that list every file
// under a prefix.
const maxListAllKeys = 100000

// BucketExists implements the invokers.S3Invoker interface. It reports
// whether the bucket exists so cross-service consumers can tell a missing
// source bucket apart from an empty one. Existence follows the bucket's
// owning region — a foreign-region address still finds the bucket.
func (s *S3Service) BucketExists(ctx context.Context, region, bucket string) (bool, error) {
	return s.bucketStoreForBucket(region, bucket).Exists(bucket), nil
}

// GetBucketPolicy implements the invokers.S3Invoker interface. It returns
// the bucket's policy document JSON, or an empty string when the bucket
// carries no policy. The bucket-record store follows the owning region,
// so a policy read addressed to a foreign region still reaches the
// bucket's records.
func (s *S3Service) GetBucketPolicy(ctx context.Context, region, bucket string) (string, error) {
	b, err := s.bucketStoreForBucket(region, bucket).Get(bucket)
	if err != nil {
		return "", fmt.Errorf("s3 GetBucketPolicy %s: %w", bucket, err)
	}
	return b.Policy, nil
}

// EnsureBucket implements the invokers.S3Invoker interface. It creates the
// bucket when it does not exist yet so services that own an internal bucket
// do not depend on manual provisioning. A concurrent creator winning the
// race is fine: the bucket exists, which is all this method guarantees.
func (s *S3Service) EnsureBucket(ctx context.Context, region, bucket string) error {
	buckets := s.s3Store.Buckets(region)
	if buckets.Exists(bucket) {
		return nil
	}
	// The store returns its own already-exists sentinel when the race is
	// lost, so compare against that sentinel rather than the service-level
	// error shape.
	if _, err := buckets.Create(bucket, region); err != nil && !errors.Is(err, s3store.ErrBucketAlreadyExists) {
		return fmt.Errorf("s3 EnsureBucket %s: %w", bucket, err)
	}
	return nil
}

// DeleteObject implements the invokers.S3Invoker interface. It removes the
// object so transient payloads can be purged after use, following the
// bucket's owning region like the writes it undoes.
func (s *S3Service) DeleteObject(ctx context.Context, region, bucket, key string) error {
	if err := s.objectStoreForBucket(region, bucket).Delete(ctx, bucket, key); err != nil {
		return fmt.Errorf("s3 DeleteObject %s/%s: %w", bucket, key, err)
	}
	return nil
}

// ListObjects lists objects in an S3 bucket via the cross-service invoker.
// When maxKeys <= 0, all objects are returned by paginating through the
// full result set, up to maxListAllKeys. When maxKeys > 0, at most maxKeys
// objects are returned from a single page. The object store follows the
// bucket's owning region, so a foreign-region address still lists the
// bucket's objects.
func (s *S3Service) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	objs := s.objectStoreForBucket(region, bucket)

	if maxKeys <= 0 {
		var allKeys []string
		marker := ""
		for {
			result, err := objs.List(bucket, prefix, "", marker, s3MaxKeys)
			if err != nil {
				return nil, fmt.Errorf("s3 ListObjects %s/%s: %w", bucket, prefix, err)
			}
			for _, o := range result.Objects {
				allKeys = append(allKeys, o.Key)
				if len(allKeys) >= maxListAllKeys {
					return nil, fmt.Errorf("s3 ListObjects %s/%s: exceeded safety cap of %d keys", bucket, prefix, maxListAllKeys)
				}
			}
			if !result.IsTruncated {
				break
			}
			marker = result.NextMarker
		}
		return allKeys, nil
	}

	result, err := objs.List(bucket, prefix, "", "", maxKeys)
	if err != nil {
		return nil, fmt.Errorf("s3 ListObjects %s/%s: %w", bucket, prefix, err)
	}
	keys := make([]string, 0, len(result.Objects))
	for _, o := range result.Objects {
		keys = append(keys, o.Key)
	}
	return keys, nil
}

// ListObjectEntries implements the invokers.S3Invoker interface. It follows
// the same pagination and safety-cap rules as ListObjects but returns the
// full metadata records the ListObjectsV2 item shape exposes (Step Functions
// Distributed Map ItemReader datasets).
func (s *S3Service) ListObjectEntries(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]invokers.S3ObjectEntry, error) {
	objs := s.objectStoreForBucket(region, bucket)

	collect := func(listResult *s3store.ObjectListResult) []invokers.S3ObjectEntry {
		entries := make([]invokers.S3ObjectEntry, 0, len(listResult.Objects))
		for _, o := range listResult.Objects {
			entries = append(entries, invokers.S3ObjectEntry{
				Key:          o.Key,
				ETag:         o.ETag,
				LastModified: o.LastModified.Unix(),
				Size:         o.Size,
				StorageClass: string(o.StorageClass),
			})
		}
		return entries
	}

	if maxKeys <= 0 {
		var all []invokers.S3ObjectEntry
		marker := ""
		for {
			result, err := objs.List(bucket, prefix, "", marker, s3MaxKeys)
			if err != nil {
				return nil, fmt.Errorf("s3 ListObjectEntries %s/%s: %w", bucket, prefix, err)
			}
			all = append(all, collect(result)...)
			if len(all) >= maxListAllKeys {
				return nil, fmt.Errorf("s3 ListObjectEntries %s/%s: exceeded safety cap of %d keys", bucket, prefix, maxListAllKeys)
			}
			if !result.IsTruncated {
				break
			}
			marker = result.NextMarker
		}
		return all, nil
	}

	result, err := objs.List(bucket, prefix, "", "", maxKeys)
	if err != nil {
		return nil, fmt.Errorf("s3 ListObjectEntries %s/%s: %w", bucket, prefix, err)
	}
	return collect(result), nil
}

var _ invokers.S3Invoker = (*S3Service)(nil)
