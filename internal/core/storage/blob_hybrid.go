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
	"log/slog"
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
	return s.putAt(s.plainAddress(bucket, key), key, reader, metadata)
}

// putVersionedUnlock streams a versioned write to the versioned address —
// the same locations GetWithVersion and GetRangeWithVersion read.
func (s *HybridBlobStore) putVersionedUnlock(bucket, key, versionId string, reader io.Reader, metadata *BlobMetadata) (*BlobMetadata, error) {
	return s.putAt(s.versionedAddress(bucket, key, versionId), key, reader, metadata)
}

// putAt streams reader into the tier the peek selects at one address. key is
// the logical object key recorded in the metadata; the address, not the key,
// decides where the bytes land.
func (s *HybridBlobStore) putAt(addr blobAddress, key string, reader io.Reader, metadata *BlobMetadata) (*BlobMetadata, error) {
	// Peek at the first threshold+1 bytes to determine whether this is a
	// small or large object without loading arbitrarily large content into
	// memory.
	peekBuf := make([]byte, s.threshold+1)
	n, peekErr := io.ReadFull(reader, peekBuf)
	peekData := peekBuf[:n]

	if metadata == nil {
		metadata = &BlobMetadata{}
	}
	metadata.Key = key

	s.cleanupAllTiers(addr)

	// io.ReadFull returns nil only when exactly len(buf) bytes were read,
	// meaning there is at least one more byte in reader → large object.
	isLarge := peekErr == nil

	if !isLarge {
		// io.ReadFull returns io.EOF (zero bytes read) or io.ErrUnexpectedEOF
		// (partial read) when the object is smaller than the buffer. Any
		// other error (network timeout, disk failure, etc.) must be propagated
		// to prevent storing truncated data as a complete object.
		if peekErr != io.EOF && peekErr != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("read error during peek: %w", peekErr)
		}
		data := peekData
		metadata.Size = int64(len(data))
		if metadata.ETag == "" {
			metadata.ETag = s.calculateETag(data)
		}
		if metadata.LastModified.IsZero() {
			metadata.LastModified = time.Now().UTC()
		}

		if err := s.putSmallObject(addr.storageKey, data, metadata); err != nil {
			return nil, err
		}
	} else {
		// Large object: combine peeked bytes with remaining reader and stream to disk.
		fullReader := io.MultiReader(bytes.NewReader(peekData), reader)
		if err := s.putLargeStreaming(addr, fullReader, metadata); err != nil {
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

func (s *HybridBlobStore) putLargeStreaming(addr blobAddress, reader io.Reader, meta *BlobMetadata) error {
	path := addr.path
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil { // #nosec G301
		return fmt.Errorf("failed to create directory: %w", err)
	}
	tmpPath := path + ".tmp." + uuid.New().String()
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	hash := md5.New()
	writer := io.MultiWriter(f, hash)
	size, copyErr := io.Copy(writer, reader)
	closeErr := f.Close()
	if copyErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write temp file: %w", copyErr)
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", closeErr)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	meta.Size = size
	if meta.ETag == "" {
		meta.ETag = hex.EncodeToString(hash.Sum(nil))
	}
	if meta.LastModified.IsZero() {
		meta.LastModified = time.Now().UTC()
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return s.storage.Bucket("blob_meta").Put([]byte(addr.storageKey), metaBytes)
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

	if obj, ok := s.getSmallObject(storageKey); ok {
		return newMemoryReader(obj.Data, obj.Metadata), obj.Metadata, nil
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

	if obj, ok := s.getSmallObject(storageKey); ok {
		start, end := clampRange(offset, length, int64(len(obj.Data)))
		return newMemoryReader(obj.Data[start:end], obj.Metadata), obj.Metadata, nil
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
	var firstErr error

	if err := s.storage.Bucket("blob_small").Delete([]byte(storageKey)); err != nil {
		firstErr = err
	}
	if err := s.storage.Bucket("blob_meta").Delete([]byte(storageKey)); err != nil && firstErr == nil {
		firstErr = err
	}

	path := s.filePath(bucket, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) && firstErr == nil {
		firstErr = err
	}
	pruneEmptyDirs(filepath.Dir(path), s.blobRoot())

	return firstErr
}

func (s *HybridBlobStore) cleanupAllTiers(addr blobAddress) {
	for _, op := range []struct {
		name string
		fn   func() error
	}{
		{"blob_small", func() error { return s.storage.Bucket("blob_small").Delete([]byte(addr.storageKey)) }},
		{"blob_meta", func() error { return s.storage.Bucket("blob_meta").Delete([]byte(addr.storageKey)) }},
		{"file", func() error {
			// A fresh overwrite target has no file-tier entry; that is
			// the common case, not a cleanup failure.
			err := os.Remove(addr.path)
			if os.IsNotExist(err) {
				return nil
			}
			if err == nil {
				// Match the delete paths: a reclaimed file prunes its
				// directory chain so a tier-switching overwrite leaves
				// no skeleton directories behind.
				pruneEmptyDirs(filepath.Dir(addr.path), s.blobRoot())
			}
			return err
		}},
	} {
		if err := op.fn(); err != nil {
			slog.Error("blob cleanup failed", "tier", op.name, "error", err)
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

	if obj, ok := s.getSmallObject(storageKey); ok {
		return obj.Metadata, nil
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
	srcStorageKey := s.storageKey(srcBucket, srcKey)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check small object store first.
	if obj, ok := s.getSmallObject(srcStorageKey); ok {
		return s.putUnlock(ctx, dstBucket, dstKey, bytes.NewReader(obj.Data), obj.Metadata)
	}

	// Large object: stream directly from file without loading entire content into memory.
	path := s.filePath(srcBucket, srcKey)
	f, info, fileErr := openAndStat(path, srcBucket+"/"+srcKey)
	if fileErr != nil {
		return nil, fileErr
	}
	defer f.Close()

	meta, _ := s.getMetadata(srcStorageKey)
	if meta == nil {
		meta = &BlobMetadata{
			Key:          srcKey,
			Size:         info.Size(),
			LastModified: info.ModTime().UTC(),
		}
	}

	return s.putUnlock(ctx, dstBucket, dstKey, f, meta)
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

// blobAddress names one logical object's storage locations: the Pebble key
// shared by the blob_small and blob_meta buckets and the file-tier path
// large objects stream to. Composing both halves in one place is what keeps
// a write and the reads that follow it on the same locations — a versioned
// write must not reach the file tier as a "key#versionId"-munged plain key,
// because sanitisation would embed the version in the final segment and
// diverge from the read path (which appends the version after sanitisation)
// for every key whose final segment rewrites (trailing "/", "/.", "/..").
type blobAddress struct {
	storageKey string
	path       string
}

// plainAddress is the address of the unversioned copy of key: the write
// target on never-versioned and suspended buckets, and the copy the
// "null"-version reads fall back to.
func (s *HybridBlobStore) plainAddress(bucket, key string) blobAddress {
	return blobAddress{
		storageKey: s.storageKey(bucket, key),
		path:       s.filePath(bucket, key),
	}
}

// versionedAddress is the address of one explicit version of key, matching
// the WithVersion read paths byte for byte.
func (s *HybridBlobStore) versionedAddress(bucket, key, versionId string) blobAddress {
	return blobAddress{
		storageKey: s.storageKeyWithVersion(bucket, key, versionId),
		path:       s.filePathWithVersion(bucket, key, versionId),
	}
}

func (s *HybridBlobStore) storageKey(bucket, key string) string {
	return bucket + "#" + escapeBlobKeyComponent(key)
}

// escapeBlobKeyComponent neutralises the "#" record separator inside a raw
// key component of the composed Pebble keys ("bucket#key[#versionId]").
// Without the escape, a legal user key such as "a#b" would compose the same
// key as the versioned address of version "b" of object "a" and the two
// objects would alias each other in the small tier. "%" is escaped first so
// the mapping stays injective (a literal "%23" in a key cannot impersonate
// an escaped "#").
func escapeBlobKeyComponent(component string) string {
	component = strings.ReplaceAll(component, "%", "%25")
	return strings.ReplaceAll(component, "#", "%23")
}

// blobRoot is the root directory every bucket's blob tree lives under.
func (s *HybridBlobStore) blobRoot() string {
	return filepath.Join(s.dataDir, "blobs")
}

// pruneEmptyDirs removes the directory chain above a deleted blob file,
// but only while the directories are empty, and never at or above stop.
// Object writes create directory nodes (MkdirAll in putLargeStreaming);
// without this walk every deleted object leaves its prefix chain behind
// and the skeletons accumulate without bound. A concurrent writer
// repopulating a directory stops the walk — ReadDir sees its file or
// in-flight temp file — so pruning cannot remove a directory another
// blob still needs.
func pruneEmptyDirs(dir, stop string) {
	for strings.HasPrefix(dir, stop+string(os.PathSeparator)) {
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
		dir = filepath.Dir(dir)
	}
}

func (s *HybridBlobStore) storageKeyWithVersion(bucket, key, versionId string) string {
	if versionId == "" {
		versionId = "null"
	}
	return bucket + "#" + escapeBlobKeyComponent(key) + "#" + versionId
}

// sanitizeKey maps an object key to a filesystem path that is both
// traversal-safe and injective: distinct keys always map to distinct
// paths, so legal AWS keys such as "a/../b" and "b", or "a//b" and
// "a/b", never share a blob file (the previous filepath.Clean-based
// mapping aliased them). Percent signs are escaped first so the mapping
// stays reversible; the path-meaningful segments (".", "..", and the
// empty segment a doubled or trailing slash produces) are rewritten to
// forms that carry no filesystem meaning; and "#" is escaped so a key
// such as "a#b" cannot collide with the version suffix of another key's
// versioned file ("f:a" + "#versionId" — see filePathWithVersion).
//
// Each segment also carries a role marker: the final segment (the blob
// file) is prefixed "f:" and every other segment (the directories it
// lives in) "d:". The S3 key namespace is flat — a key and its prefix
// extension ("a" and "a/b") are independent objects — but the
// filesystem namespace distinguishes files from directories at one
// path. The markers keep file names and directory names disjoint, so
// one key's file can never sit where another key's directory must be.
func (s *HybridBlobStore) sanitizeKey(key string) string {
	segments := strings.Split(key, "/")
	for i, seg := range segments {
		seg = strings.ReplaceAll(seg, "%", "%25")
		seg = strings.ReplaceAll(seg, "#", "%23")
		switch seg {
		case "":
			seg = "%2F"
		case ".":
			seg = "%2E"
		case "..":
			seg = "%2E%2E"
		}
		if i == len(segments)-1 {
			seg = "f:" + seg
		} else {
			seg = "d:" + seg
		}
		segments[i] = seg
	}
	return strings.Join(segments, string(os.PathSeparator))
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

// getSmallObject reads the small-tier record for a storage key. A record
// that fails to unmarshal is corruption, not absence: it is logged so the
// gap stays diagnosable, and treated as absent — the file tier holds no
// copy of a small object, so the caller's not-found follows traceably.
func (s *HybridBlobStore) getSmallObject(storageKey string) (*smallObject, bool) {
	data, err := s.storage.Bucket("blob_small").Get([]byte(storageKey))
	if err != nil || data == nil {
		return nil, false
	}
	var obj smallObject
	if err := json.Unmarshal(data, &obj); err != nil {
		slog.Error("blob small-tier record is corrupt; treating as absent", "storageKey", storageKey, "error", err)
		return nil, false
	}
	return &obj, true
}
