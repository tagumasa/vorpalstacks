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
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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
	if err := os.WriteFile(path, data, 0644); err != nil { // #nosec G306
		return fmt.Errorf("failed to write file: %w", err)
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
	// #nosec G304
	fileData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("object not found: %s/%s", bucket, key)
		}
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	meta, err := s.getMetadata(storageKey)
	if err != nil {
		meta = &BlobMetadata{
			Key:          key,
			Size:         int64(len(fileData)),
			ETag:         s.calculateETag(fileData),
			LastModified: time.Now().UTC(),
		}
	}

	return newMemoryReader(fileData, meta), meta, nil
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
	reader, meta, err := s.Get(ctx, bucket, key)
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()

	mr := reader.(*memoryReader)
	start := offset
	if start < 0 {
		start = 0
	}
	end := start + length
	if end > int64(len(mr.data)) {
		end = int64(len(mr.data))
	}

	rangedData := mr.data[start:end]
	return newMemoryReader(rangedData, meta), meta, nil
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
	reader, meta, err := s.GetWithVersion(ctx, bucket, key, versionId)
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()

	mr := reader.(*memoryReader)
	start := offset
	if start < 0 {
		start = 0
	}
	end := start + length
	if end > int64(len(mr.data)) {
		end = int64(len(mr.data))
	}

	rangedData := mr.data[start:end]
	return newMemoryReader(rangedData, meta), meta, nil
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
		log.Printf("Failed to delete from blob_small: %v", err)
	}
	if err := s.storage.Bucket("blob_meta").Delete([]byte(storageKey)); err != nil {
		log.Printf("Failed to delete from blob_meta: %v", err)
	}

	path := s.filePath(bucket, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove file %s: %v", path, err)
	}

	return nil
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
	fileData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("object not found: %s/%s (version %s)", bucket, key, versionId)
		}
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	meta, err := s.getMetadata(storageKey)
	if err != nil {
		meta = &BlobMetadata{
			Key:          key,
			Size:         int64(len(fileData)),
			ETag:         s.calculateETag(fileData),
			LastModified: time.Now().UTC(),
		}
	}

	return newMemoryReader(fileData, meta), meta, nil
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
		log.Printf("Failed to delete from blob_small: %v", err)
	}
	if err := s.storage.Bucket("blob_meta").Delete([]byte(storageKey)); err != nil {
		log.Printf("Failed to delete from blob_meta: %v", err)
	}

	path := s.filePathWithVersion(bucket, key, versionId)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove file %s: %v", path, err)
	}

	return nil
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

// CreateMultipartUpload initiates a multipart upload for an object.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - metadata: Optional metadata for the object
//
// Returns:
//   - string: The upload ID
//   - error: An error if the operation fails
func (s *HybridBlobStore) CreateMultipartUpload(ctx context.Context, bucket, key string, metadata *BlobMetadata) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	uploadID := uuid.New().String()
	uploadDir := s.uploadDir(uploadID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil { // #nosec G301
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	upload := &multipartUploadState{
		Bucket:   bucket,
		Key:      key,
		Metadata: metadata,
	}
	data, err := json.Marshal(upload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal upload state: %w", err)
	}
	if err := s.storage.Bucket("blob_uploads").Put([]byte(uploadID), data); err != nil {
		return "", fmt.Errorf("failed to store upload state: %w", err)
	}

	return uploadID, nil
}

// UploadPart uploads a part of a multipart upload.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - uploadID: The upload ID from CreateMultipartUpload
//   - partNumber: The part number (1-based)
//   - reader: The part data reader
//
// Returns:
//   - string: The ETag for the uploaded part
//   - error: An error if the operation fails
func (s *HybridBlobStore) UploadPart(ctx context.Context, bucket, key, uploadID string, partNumber int, reader io.Reader) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read part: %w", err)
	}

	etag := s.calculateETag(data)
	path := s.partPath(uploadID, partNumber)
	if err := os.WriteFile(path, data, 0644); err != nil { // #nosec G306
		return "", fmt.Errorf("failed to write part: %w", err)
	}

	return etag, nil
}

// CompleteMultipartUpload completes a multipart upload by combining all parts.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - uploadID: The upload ID from CreateMultipartUpload
//   - parts: The list of part information
//
// Returns:
//   - *BlobMetadata: The final object metadata
//   - error: An error if the operation fails
func (s *HybridBlobStore) CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []PartInfo) (*BlobMetadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var combined bytes.Buffer
	var partHashes bytes.Buffer
	for _, part := range parts {
		data, err := os.ReadFile(s.partPath(uploadID, part.PartNumber))
		if err != nil {
			return nil, fmt.Errorf("failed to read part %d: %w", part.PartNumber, err)
		}
		combined.Write(data)
		partHash := md5.Sum(data)
		partHashes.Write(partHash[:])
	}

	if err := s.abortMultipartUploadUnlock(uploadID); err != nil {
		log.Printf("Failed to cleanup multipart upload %s: %v", uploadID, err)
	}

	var metadata *BlobMetadata
	if len(parts) > 0 {
		combinedHash := md5.Sum(partHashes.Bytes())
		metadata = &BlobMetadata{
			ETag: hex.EncodeToString(combinedHash[:]) + "-" + fmt.Sprintf("%d", len(parts)),
		}
	}

	return s.putUnlock(ctx, bucket, key, &combined, metadata)
}

// AbortMultipartUpload aborts a multipart upload, removing all uploaded parts.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - uploadID: The upload ID from CreateMultipartUpload
//
// Returns:
//   - error: An error if the operation fails
func (s *HybridBlobStore) AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.abortMultipartUploadUnlock(uploadID)
}

func (s *HybridBlobStore) abortMultipartUploadUnlock(uploadID string) error {
	if err := os.RemoveAll(s.uploadDir(uploadID)); err != nil {
		log.Printf("Failed to remove upload dir: %v", err)
	}
	if err := s.storage.Bucket("blob_uploads").Delete([]byte(uploadID)); err != nil {
		log.Printf("Failed to delete upload metadata: %v", err)
	}
	return nil
}

// ListParts lists all parts of a multipart upload.
//
// Parameters:
//   - ctx: The context for the operation
//   - bucket: The bucket name
//   - key: The object key
//   - uploadID: The upload ID from CreateMultipartUpload
//
// Returns:
//   - []PartInfo: The list of uploaded parts
//   - error: An error if the operation fails
func (s *HybridBlobStore) ListParts(ctx context.Context, bucket, key, uploadID string) ([]PartInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uploadDir := s.uploadDir(uploadID)
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return nil, err
	}

	var parts []PartInfo
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "part.") {
			var partNum int
			if _, err := fmt.Sscanf(entry.Name(), "part.%d", &partNum); err != nil {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			partPath := filepath.Join(uploadDir, entry.Name())
			data, err := os.ReadFile(partPath)
			var etag string
			if err == nil {
				etag = s.calculateETag(data)
			}

			parts = append(parts, PartInfo{
				PartNumber: partNum,
				Size:       info.Size(),
				ETag:       etag,
			})
		}
	}
	return parts, nil
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

	return os.MkdirAll(filepath.Join(s.dataDir, "blobs", name), 0755) // #nosec G301
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

	return os.RemoveAll(filepath.Join(s.dataDir, "blobs", name))
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

func (s *HybridBlobStore) filePath(bucket, key string) string {
	sanitized := strings.ReplaceAll(key, "/", string(os.PathSeparator))
	sanitized = filepath.Clean(sanitized)
	if strings.Contains(sanitized, "..") {
		sanitized = filepath.Base(sanitized)
	}
	return filepath.Join(s.dataDir, "blobs", bucket, sanitized)
}

func (s *HybridBlobStore) filePathWithVersion(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	sanitized := strings.ReplaceAll(key, "/", string(os.PathSeparator))
	sanitized = filepath.Clean(sanitized)
	if strings.Contains(sanitized, "..") {
		sanitized = filepath.Base(sanitized)
	}
	return filepath.Join(s.dataDir, "blobs", bucket, sanitized+"#"+versionId)
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

type memoryReader struct {
	*bytes.Reader
	data []byte
	meta *BlobMetadata
}

func newMemoryReader(data []byte, meta *BlobMetadata) *memoryReader {
	return &memoryReader{
		Reader: bytes.NewReader(data),
		data:   data,
		meta:   meta,
	}
}

// Size returns the size of the stored data in bytes.
func (r *memoryReader) Size() int64 { return int64(len(r.data)) }

// ETag returns the ETag of the stored data.
func (r *memoryReader) ETag() string { return r.meta.ETag }

// Close closes the memory reader. This is a no-op since the data is held in memory.
func (r *memoryReader) Close() error { return nil }
