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

	"github.com/google/uuid"
)

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
		slog.Error("Failed to cleanup multipart upload", "uploadID", uploadID, "error", err)
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
		slog.Error("Failed to remove upload dir", "error", err)
	}
	if err := s.storage.Bucket("blob_uploads").Delete([]byte(uploadID)); err != nil {
		slog.Error("Failed to delete upload metadata", "error", err)
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
