package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

// GetWithVersion retrieves a specific version of an object.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - versionId: The version ID
//
// Returns:
//   - BlobReader: A reader for the object data
//   - *BlobMetadata: The object metadata
//   - error: An error if the object is not found or the operation fails
func (s *HybridBlobStore) GetWithVersion(ctx context.Context, bucket, key, versionId string) (BlobReader, *BlobMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storageKey := s.storageKeyWithVersion(bucket, key, versionId)

	data, err := s.storage.Bucket("blob_small").Get([]byte(storageKey))
	if err == nil && data != nil {
		var obj smallObject
		if err := json.Unmarshal(data, &obj); err == nil {
			return newMemoryReader(obj.Data, obj.Metadata), obj.Metadata, nil
		}
	}

	path := s.filePathWithVersion(bucket, key, versionId)
	f, info, err := openAndStat(path, fmt.Sprintf("%s/%s (version %s)", bucket, key, versionId))
	if err != nil {
		return nil, nil, err
	}

	meta, err := s.getMetadata(storageKey)
	if err != nil {
		meta = &BlobMetadata{
			Key:          key,
			Size:         info.Size(),
			LastModified: info.ModTime().UTC(),
		}
	}

	return newFileReader(f, info.Size(), meta), meta, nil
}

// GetRangeWithVersion retrieves a range of bytes from a specific version of an object.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - versionId: The version ID
//   - offset: The starting offset
//   - length: The number of bytes to read
//
// Returns:
//   - BlobReader: A reader for the object data
//   - *BlobMetadata: The object metadata
//   - error: An error if the object is not found or the operation fails
func (s *HybridBlobStore) GetRangeWithVersion(ctx context.Context, bucket, key, versionId string, offset, length int64) (BlobReader, *BlobMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storageKey := s.storageKeyWithVersion(bucket, key, versionId)

	data, err := s.storage.Bucket("blob_small").Get([]byte(storageKey))
	if err == nil && data != nil {
		var obj smallObject
		if err := json.Unmarshal(data, &obj); err == nil {
			start, end := clampRange(offset, length, int64(len(obj.Data)))
			return newMemoryReader(obj.Data[start:end], obj.Metadata), obj.Metadata, nil
		}
	}

	path := s.filePathWithVersion(bucket, key, versionId)
	f, info, err := openAndStat(path, fmt.Sprintf("%s/%s (version %s)", bucket, key, versionId))
	if err != nil {
		return nil, nil, err
	}

	start, end := clampRange(offset, length, info.Size())

	meta, err := s.getMetadata(storageKey)
	if err != nil {
		meta = &BlobMetadata{
			Key:          key,
			Size:         info.Size(),
			LastModified: info.ModTime().UTC(),
		}
	}

	return newSectionFileReader(f, start, end-start, meta), meta, nil
}

// DeleteWithVersion removes a specific version of an object.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - versionId: The version ID
//
// Returns:
//   - error: An error if the operation fails
func (s *HybridBlobStore) DeleteWithVersion(ctx context.Context, bucket, key, versionId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	storageKey := s.storageKeyWithVersion(bucket, key, versionId)

	if err := s.storage.Bucket("blob_small").Delete([]byte(storageKey)); err != nil {
		slog.Error("Failed to delete from blob_small", "error", err)
	}
	if err := s.storage.Bucket("blob_meta").Delete([]byte(storageKey)); err != nil {
		slog.Error("Failed to delete from blob_meta", "error", err)
	}

	path := s.filePathWithVersion(bucket, key, versionId)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		slog.Error("Failed to remove file", "path", path, "error", err)
	}

	return nil
}

// CopyWithVersion copies a specific version of an object to a new location.
//
// Parameters:
//   - ctx: The context for the operation
//   - srcBucket: The source bucket name
//   - srcKey: The source object key
//   - srcVersionId: The source version ID
//   - dstBucket: The destination bucket name
//   - dstKey: The destination object key
//
// Returns:
//   - *BlobMetadata: The destination object metadata
//   - error: An error if the operation fails
func (s *HybridBlobStore) CopyWithVersion(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string) (*BlobMetadata, error) {
	reader, meta, err := s.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return s.Put(ctx, dstBucket, dstKey, reader, meta)
}
