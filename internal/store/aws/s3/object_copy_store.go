package s3

import (
	"context"
	"io"
)

// copyObjectReader validates source and destination keys, opens a reader for
// the source object (optionally at a specific version), and delegates to
// PutWithVersioning. It is the shared implementation for all four public
// Copy variants, which differ only in whether they override content type,
// metadata, or version.
func (s *ObjectStore) copyObjectReader(ctx context.Context, reader io.ReadCloser, dstBucket, dstKey string, contentType string, metadata map[string]string, storageClass ObjectStorageClass) (*Object, error) {
	defer reader.Close()
	if storageClass == "" {
		storageClass = StorageClassStandard
	}
	return s.PutWithVersioning(ctx, dstBucket, dstKey, reader, contentType, metadata, false, storageClass, nil)
}

// openObjectReader returns a reader and source object metadata for the given
// bucket/key. When versionId is non-empty and the bucket has versioning
// enabled the versioned blob is fetched directly; otherwise the standard
// (latest) path is used.
func (s *ObjectStore) openObjectReader(ctx context.Context, bucket, key, versionId string) (io.ReadCloser, *Object, error) {
	if versionId != "" && s.isVersioningEnabled(bucket) {
		obj, err := s.HeadWithVersion(ctx, bucket, key, versionId)
		if err != nil {
			return nil, nil, err
		}
		blobReader, _, rErr := s.blobStore.GetWithVersion(ctx, bucket, key, versionId)
		if rErr != nil {
			return nil, nil, rErr
		}
		return blobReader, obj, nil
	}
	return s.Get(ctx, bucket, key)
}

// Copy copies an object from one location to another, preserving the source
// content type, metadata, and storage class.
func (s *ObjectStore) Copy(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}
	reader, srcObj, err := s.Get(ctx, srcBucket, srcKey)
	if err != nil {
		return nil, err
	}
	return s.copyObjectReader(ctx, reader, dstBucket, dstKey, srcObj.ContentType, srcObj.Metadata, srcObj.StorageClass)
}

// CopyWithMetadata copies an object with custom content type and metadata.
func (s *ObjectStore) CopyWithMetadata(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string, contentType string, metadata map[string]string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}
	reader, srcObj, err := s.Get(ctx, srcBucket, srcKey)
	if err != nil {
		return nil, err
	}
	ct := contentType
	if ct == "" {
		ct = srcObj.ContentType
	}
	return s.copyObjectReader(ctx, reader, dstBucket, dstKey, ct, metadata, srcObj.StorageClass)
}

// CopyWithVersion copies a specific version of an object.
func (s *ObjectStore) CopyWithVersion(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}
	reader, srcObj, err := s.openObjectReader(ctx, srcBucket, srcKey, srcVersionId)
	if err != nil {
		return nil, err
	}
	return s.copyObjectReader(ctx, reader, dstBucket, dstKey, srcObj.ContentType, srcObj.Metadata, srcObj.StorageClass)
}

// CopyWithVersionAndMetadata copies a specific version of an object with custom content type and metadata.
func (s *ObjectStore) CopyWithVersionAndMetadata(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string, contentType string, metadata map[string]string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}
	reader, srcObj, err := s.openObjectReader(ctx, srcBucket, srcKey, srcVersionId)
	if err != nil {
		return nil, err
	}
	ct := contentType
	if ct == "" {
		ct = srcObj.ContentType
	}
	return s.copyObjectReader(ctx, reader, dstBucket, dstKey, ct, metadata, srcObj.StorageClass)
}
