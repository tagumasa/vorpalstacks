// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"vorpalstacks/internal/utils/naming"
)

// HybridBlobStore implements a hybrid blob storage strategy that stores small objects
// in a key-value store and large objects as files on disk. It provides automatic
// tiering based on object size.
type HybridBlobStore struct {
	storage   BasicStorage
	dataDir   string
	threshold int64
	mu        sync.RWMutex
}

// NewHybridBlobStore creates a new hybrid blob store with the given storage backend
// and data directory for file-based storage of large objects.
func NewHybridBlobStore(storage BasicStorage, dataDir string) (*HybridBlobStore, error) {
	for _, dir := range []string{
		filepath.Join(dataDir, "blobs"),
		filepath.Join(dataDir, "uploads"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil { // #nosec G301
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return &HybridBlobStore{
		storage:   storage,
		dataDir:   dataDir,
		threshold: SmallObjectThreshold,
	}, nil
}

// Put stores an object in the hybrid blob store. Small objects are stored
// in the key-value store, while large objects are stored as files on disk.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - reader: The data reader
//   - metadata: Optional metadata for the object
//
// Returns:
//   - *BlobMetadata: The stored metadata
//   - error: An error if the operation fails
func (s *HybridBlobStore) Put(ctx context.Context, bucket, key string, reader io.Reader, metadata *BlobMetadata) (*BlobMetadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.putUnlock(ctx, bucket, key, reader, metadata)
}

func (s *HybridBlobStore) putUnlock(ctx context.Context, bucket, key string, reader io.Reader, metadata *BlobMetadata) (*BlobMetadata, error) {
	var buf bytes.Buffer
	size, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	data := buf.Bytes()

	if metadata == nil {
		metadata = &BlobMetadata{}
	}
	metadata.Key = key
	metadata.Size = size
	if metadata.ETag == "" {
		metadata.ETag = s.calculateETag(data)
	}
	metadata.LastModified = time.Now().UTC()

	storageKey := s.storageKey(bucket, key)

	s.cleanupAllTiers(storageKey, bucket, key)

	if size < s.threshold {
		if err := s.putSmallObject(storageKey, data, metadata); err != nil {
			return nil, err
		}
	} else {
		if err := s.putLargeObject(bucket, key, data, metadata); err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

func (s *HybridBlobStore) putSmallObject(key string, data []byte, meta *BlobMetadata) error {
	obj := &smallObject{
		Data:     data,
		Metadata: meta,
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal small object: %w", err)
	}
	return s.storage.Bucket("blob_small").Put([]byte(key), bytes)
}

func (s *HybridBlobStore) putLargeObject(bucket, key string, data []byte, meta *BlobMetadata) error {
	path := s.filePath(bucket, key)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil { // #nosec G301
		return fmt.Errorf("failed to create directory: %w", err)
	}
	tmpPath := path + ".tmp." + uuid.New().String()
	if err := os.WriteFile(tmpPath, data, 0644); err != nil { // #nosec G306
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	storageKey := s.storageKey(bucket, key)
	return s.storage.Bucket("blob_meta").Put([]byte(storageKey), metaBytes)
}

// Get retrieves an object from the hybrid blob store.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//
// Returns:
//   - BlobReader: A reader for the object data
//   - *BlobMetadata: The object metadata
//   - error: An error if the object is not found or the operation fails
func (s *HybridBlobStore) Get(ctx context.Context, bucket, key string) (BlobReader, *BlobMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storageKey := s.storageKey(bucket, key)

	data, err := s.storage.Bucket("blob_small").Get([]byte(storageKey))
	if err == nil && data != nil {
		var obj smallObject
		if err := json.Unmarshal(data, &obj); err == nil {
			return newMemoryReader(obj.Data, obj.Metadata), obj.Metadata, nil
		}
	}

	path := s.filePath(bucket, key)
	f, info, err := openAndStat(path, bucket+"/"+key)
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

func (s *HybridBlobStore) getMetadata(storageKey string) (*BlobMetadata, error) {
	data, err := s.storage.Bucket("blob_meta").Get([]byte(storageKey))
	if err != nil || data == nil {
		return nil, fmt.Errorf("metadata not found")
	}
	var meta BlobMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// GetRange retrieves a range of bytes from an object.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - offset: The starting offset
//   - length: The number of bytes to read
//
// Returns:
//   - BlobReader: A reader for the object data
//   - *BlobMetadata: The object metadata
//   - error: An error if the object is not found or the operation fails
func (s *HybridBlobStore) GetRange(ctx context.Context, bucket, key string, offset, length int64) (BlobReader, *BlobMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storageKey := s.storageKey(bucket, key)

	data, err := s.storage.Bucket("blob_small").Get([]byte(storageKey))
	if err == nil && data != nil {
		var obj smallObject
		if err := json.Unmarshal(data, &obj); err == nil {
			start, end := clampRange(offset, length, int64(len(obj.Data)))
			return newMemoryReader(obj.Data[start:end], obj.Metadata), obj.Metadata, nil
		}
	}

	path := s.filePath(bucket, key)
	f, info, err := openAndStat(path, bucket+"/"+key)
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

// Delete removes an object from the hybrid blob store.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//
// Returns:
//   - error: An error if the operation fails
func (s *HybridBlobStore) Delete(ctx context.Context, bucket, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.deleteUnlock(bucket, key)
}

func (s *HybridBlobStore) deleteUnlock(bucket, key string) error {
	storageKey := s.storageKey(bucket, key)

	if err := s.storage.Bucket("blob_small").Delete([]byte(storageKey)); err != nil {
		return err
	}
	if err := s.storage.Bucket("blob_meta").Delete([]byte(storageKey)); err != nil {
		return err
	}

	path := s.filePath(bucket, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (s *HybridBlobStore) cleanupAllTiers(storageKey, bucket, key string) {
	for _, op := range []struct {
		name string
		fn   func() error
	}{
		{"blob_small", func() error { return s.storage.Bucket("blob_small").Delete([]byte(storageKey)) }},
		{"blob_meta", func() error { return s.storage.Bucket("blob_meta").Delete([]byte(storageKey)) }},
		{"file", func() error { return os.Remove(s.filePath(bucket, key)) }},
	} {
		if err := op.fn(); err != nil {
			fmt.Printf("[WARN] blob cleanup failed: tier=%s err=%v\n", op.name, err)
		}
	}
}

// Exists checks if an object exists in the hybrid blob store.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//
// Returns:
//   - bool: True if the object exists
//   - error: An error if the operation fails
func (s *HybridBlobStore) Exists(ctx context.Context, bucket, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storageKey := s.storageKey(bucket, key)

	if s.storage.Bucket("blob_small").Has([]byte(storageKey)) {
		return true, nil
	}

	path := s.filePath(bucket, key)
	_, err := os.Stat(path)
	return err == nil, nil
}

// Head retrieves metadata for an object without returning the object data.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//
// Returns:
//   - *BlobMetadata: The object metadata
//   - error: An error if the object is not found or the operation fails
func (s *HybridBlobStore) Head(ctx context.Context, bucket, key string) (*BlobMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storageKey := s.storageKey(bucket, key)

	data, err := s.storage.Bucket("blob_small").Get([]byte(storageKey))
	if err == nil && data != nil {
		var obj smallObject
		if err := json.Unmarshal(data, &obj); err == nil {
			return obj.Metadata, nil
		}
	}

	meta, err := s.getMetadata(storageKey)
	if err == nil {
		return meta, nil
	}

	return nil, fmt.Errorf("object not found: %s/%s", bucket, key)
}

// Copy copies an object from one location to another.
//
// Parameters:
//   - ctx: The context for the operation
//   - srcBucket: The source bucket name
//   - srcKey: The source object key
//   - dstBucket: The destination bucket name
//   - dstKey: The destination object key
//
// Returns:
//   - *BlobMetadata: The destination object metadata
//   - error: An error if the operation fails
func (s *HybridBlobStore) Copy(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) (*BlobMetadata, error) {
	reader, meta, err := s.Get(ctx, srcBucket, srcKey)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return s.Put(ctx, dstBucket, dstKey, reader, meta)
}

// CreateBucket creates a new bucket in the hybrid blob store.
//
// Parameters:
//   - ctx: The context for the operation
//   - name: The bucket name
//
// Returns:
//   - error: An error if the operation fails
func (s *HybridBlobStore) CreateBucket(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucketDir := filepath.Join(s.dataDir, "blobs", name)
	if _, err := naming.ValidatePathWithinDir(filepath.Join(s.dataDir, "blobs"), name); err != nil {
		return fmt.Errorf("invalid bucket name: %w", err)
	}
	return os.MkdirAll(bucketDir, 0755) // #nosec G301
}

// DeleteBucket deletes a bucket from the hybrid blob store.
//
// Parameters:
//   - ctx: The context for the operation
//   - name: The bucket name
//
// Returns:
//   - error: An error if the operation fails
func (s *HybridBlobStore) DeleteBucket(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucketDir := filepath.Join(s.dataDir, "blobs", name)
	if _, err := naming.ValidatePathWithinDir(filepath.Join(s.dataDir, "blobs"), name); err != nil {
		return fmt.Errorf("invalid bucket name: %w", err)
	}
	return os.RemoveAll(bucketDir)
}

func (s *HybridBlobStore) storageKey(bucket, key string) string {
	return bucket + "#" + key
}

func (s *HybridBlobStore) storageKeyWithVersion(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	return bucket + "#" + key + "#" + versionId
}

func (s *HybridBlobStore) sanitizeKey(key string) string {
	sanitized := strings.ReplaceAll(key, "/", string(os.PathSeparator))
	sanitized = filepath.Clean(sanitized)
	if strings.Contains(sanitized, "..") {
		sanitized = filepath.Base(sanitized)
	}
	return sanitized
}

func (s *HybridBlobStore) filePath(bucket, key string) string {
	return filepath.Join(s.dataDir, "blobs", bucket, s.sanitizeKey(key))
}

func (s *HybridBlobStore) filePathWithVersion(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	return filepath.Join(s.dataDir, "blobs", bucket, s.sanitizeKey(key)+"#"+versionId)
}

func (s *HybridBlobStore) uploadDir(uploadID string) string {
	return filepath.Join(s.dataDir, "uploads", uploadID)
}

func (s *HybridBlobStore) partPath(uploadID string, partNum int) string {
	return filepath.Join(s.uploadDir(uploadID), fmt.Sprintf("part.%d", partNum))
}

func (s *HybridBlobStore) calculateETag(data []byte) string {
	// #nosec G401 - MD5 is standard for S3 ETags
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

type smallObject struct {
	Data     []byte        `json:"data"`
	Metadata *BlobMetadata `json:"metadata"`
}

type multipartUploadState struct {
	Bucket   string        `json:"bucket"`
	Key      string        `json:"key"`
	Metadata *BlobMetadata `json:"metadata"`
}
