package s3

import (
	"context"
	"io"
)

// GetRange retrieves a range of bytes from an object.
func (s *ObjectStore) GetRange(ctx context.Context, bucket, key string, offset, length int64) (io.ReadCloser, *Object, error) {
	return s.GetRangeWithVersion(ctx, bucket, key, "", offset, length)
}
