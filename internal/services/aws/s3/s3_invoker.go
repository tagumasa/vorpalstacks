package s3

import (
	"context"
	"fmt"
	"io"

	"vorpalstacks/internal/eventbus"
)

// GetObject implements the eventbus.S3Invoker interface. It retrieves the
// content of an object by region, bucket and key, returning the full byte
// content.
func (s *S3Service) GetObject(ctx context.Context, region, bucket, key string) ([]byte, error) {
	objs := s.s3Store.Objects(region)
	reader, _, err := objs.Get(ctx, bucket, key)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject %s/%s: %w", bucket, key, err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("s3 GetObject read %s/%s: %w", bucket, key, err)
	}
	return data, nil
}

// ListObjects implements the eventbus.S3Invoker interface. It returns object
// keys under the given region, bucket and prefix.
func (s *S3Service) ListObjects(ctx context.Context, region, bucket, prefix string) ([]string, error) {
	objs := s.s3Store.Objects(region)
	result, err := objs.List(bucket, prefix, "", "", 1000)
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
