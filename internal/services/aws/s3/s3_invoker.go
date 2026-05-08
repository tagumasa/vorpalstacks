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

func (s *S3Service) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	objs := s.s3Store.Objects(region)
	_, err := objs.Put(ctx, bucket, key, bytes.NewReader(data), contentType, nil)
	if err != nil {
		return fmt.Errorf("s3 PutObject %s/%s: %w", bucket, key, err)
	}
	return nil
}

func (s *S3Service) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	objs := s.s3Store.Objects(region)
	if maxKeys <= 0 {
		maxKeys = 1000
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
