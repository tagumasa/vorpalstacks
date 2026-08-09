package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"vorpalstacks/internal/eventbus"
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
