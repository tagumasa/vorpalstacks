package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// GetObject implements the eventbus.S3Invoker interface. It retrieves the
// content of an object by region, bucket and key, returning the full byte
// content.
func (s *S3Service) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	objs := s.s3Store.Objects(region)
	reader, _, err := objs.Get(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject %s/%s: %w", bucket, key, err)
	}
	defer reader.Close()
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024 * 1024
	}
	data, err := io.ReadAll(io.LimitReader(reader, maxBytes))
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject read %s/%s: %w", bucket, key, err)
	}
	return data, nil
}

// GetObjectVersion implements the eventbus.S3Invoker interface. It
// retrieves the content of a specific object version; an empty versionID
// reads the latest version, matching the store's version-aware read.
func (s *S3Service) GetObjectVersion(ctx context.Context, region, bucket, key, versionID string, maxBytes int64) ([]byte, error) {
	objs := s.s3Store.Objects(region)
	reader, _, err := objs.GetWithVersion(ctx, bucket, key, versionID)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObjectVersion %s/%s@%s: %w", bucket, key, versionID, err)
	}
	defer reader.Close()
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024 * 1024
	}
	data, err := io.ReadAll(io.LimitReader(reader, maxBytes))
	if err != nil {
		return nil, fmt.Errorf("s3 GetObjectVersion read %s/%s@%s: %w", bucket, key, versionID, err)
	}
	return data, nil
}

// PutObject stores an object in S3 via the cross-service invoker.
func (s *S3Service) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	objs := s.s3Store.Objects(region)
	_, err := objs.Put(ctx, bucket, key, bytes.NewReader(data), contentType, nil)
	if err != nil {
		return fmt.Errorf("s3 PutObject %s/%s: %w", bucket, key, err)
	}
	return nil
}

// maxListAllKeys is the safety cap when listing all objects (maxKeys <= 0).
// This prevents unbounded memory consumption on extremely large buckets
// while accommodating legitimate bulk-load workloads that list every file
// under a prefix.
const maxListAllKeys = 100000

// BucketExists implements the eventbus.S3Invoker interface. It reports
// whether the bucket exists so cross-service consumers can tell a missing
// source bucket apart from an empty one.
func (s *S3Service) BucketExists(ctx context.Context, region, bucket string) (bool, error) {
	return s.s3Store.Buckets(region).Exists(bucket), nil
}

// EnsureBucket implements the eventbus.S3Invoker interface. It creates the
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

// DeleteObject implements the eventbus.S3Invoker interface. It removes the
// object so transient payloads can be purged after use.
func (s *S3Service) DeleteObject(ctx context.Context, region, bucket, key string) error {
	if err := s.s3Store.Objects(region).Delete(ctx, bucket, key); err != nil {
		return fmt.Errorf("s3 DeleteObject %s/%s: %w", bucket, key, err)
	}
	return nil
}

// ListObjects lists objects in an S3 bucket via the cross-service invoker.
// When maxKeys <= 0, all objects are returned by paginating through the
// full result set, up to maxListAllKeys. When maxKeys > 0, at most maxKeys
// objects are returned from a single page.
func (s *S3Service) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	objs := s.s3Store.Objects(region)

	if maxKeys <= 0 {
		var allKeys []string
		marker := ""
		for {
			result, err := objs.List(bucket, prefix, "", marker, 1000)
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

var _ eventbus.S3Invoker = (*S3Service)(nil)
