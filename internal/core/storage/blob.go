// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"context"
	"io"
	"time"
)

// SmallObjectThreshold defines the size threshold (in bytes) for small objects.
const SmallObjectThreshold = 4 * 1024 * 1024

// BlobMetadata contains metadata for a blob.
type BlobMetadata struct {
	Key           string
	Size          int64
	ETag          string
	ContentType   string
	ContentMD5    string
	LastModified  time.Time
	CustomHeaders map[string]string
}

// BlobReader provides read access to a blob.
type BlobReader interface {
	io.ReadCloser
	Size() int64
	ETag() string
}

// PartInfo represents information about an uploaded part.
type PartInfo struct {
	PartNumber int
	ETag       string
	Size       int64
}

// MultipartUpload represents an in-progress multipart upload.
type MultipartUpload struct {
	UploadID  string
	Key       string
	Bucket    string
	Initiated time.Time
	Metadata  *BlobMetadata
}

// BlobStore provides an interface for storing and retrieving blobs.
type BlobStore interface {
	Put(ctx context.Context, bucket, key string, reader io.Reader, metadata *BlobMetadata) (*BlobMetadata, error)
	Get(ctx context.Context, bucket, key string) (BlobReader, *BlobMetadata, error)
	GetRange(ctx context.Context, bucket, key string, offset, length int64) (BlobReader, *BlobMetadata, error)
	Delete(ctx context.Context, bucket, key string) error
	Exists(ctx context.Context, bucket, key string) (bool, error)
	Head(ctx context.Context, bucket, key string) (*BlobMetadata, error)

	GetWithVersion(ctx context.Context, bucket, key, versionId string) (BlobReader, *BlobMetadata, error)
	GetRangeWithVersion(ctx context.Context, bucket, key, versionId string, offset, length int64) (BlobReader, *BlobMetadata, error)
	DeleteWithVersion(ctx context.Context, bucket, key, versionId string) error
	CopyWithVersion(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string) (*BlobMetadata, error)

	CreateMultipartUpload(ctx context.Context, bucket, key string, metadata *BlobMetadata) (uploadID string, err error)
	UploadPart(ctx context.Context, bucket, key, uploadID string, partNumber int, reader io.Reader) (etag string, err error)
	CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []PartInfo) (*BlobMetadata, error)
	AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error
	ListParts(ctx context.Context, bucket, key, uploadID string) ([]PartInfo, error)

	Copy(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) (*BlobMetadata, error)

	CreateBucket(ctx context.Context, name string) error
	DeleteBucket(ctx context.Context, name string) error
}
