package s3

import (
	"context"
	"io"
)

// Copy copies an object from one location to another.
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
	defer reader.Close()

	return s.PutWithVersioning(ctx, dstBucket, dstKey, reader, srcObj.ContentType, srcObj.Metadata, false, StorageClassStandard, nil)
}

// CopyWithMetadata copies an object with custom content type and metadata.
func (s *ObjectStore) CopyWithMetadata(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string, contentType string, metadata map[string]string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}
	reader, _, err := s.Get(ctx, srcBucket, srcKey)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return s.PutWithVersioning(ctx, dstBucket, dstKey, reader, contentType, metadata, false, StorageClassStandard, nil)
}

// CopyWithVersion copies a specific version of an object.
func (s *ObjectStore) CopyWithVersion(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}

	_, srcObj, err := s.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	if s.isVersioningEnabled(srcBucket) {
		reader, _, err = s.blobStore.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		reader, _, err = s.blobStore.Get(ctx, srcBucket, srcKey)
	}
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return s.PutWithVersioning(ctx, dstBucket, dstKey, reader, srcObj.ContentType, srcObj.Metadata, false, StorageClassStandard, nil)
}

// CopyWithVersionAndMetadata copies a specific version of an object with custom content type and metadata.
func (s *ObjectStore) CopyWithVersionAndMetadata(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string, contentType string, metadata map[string]string) (*Object, error) {
	if err := validateS3Key(srcKey); err != nil {
		return nil, err
	}
	if err := validateS3Key(dstKey); err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	var err error
	if s.isVersioningEnabled(srcBucket) && srcVersionId != "" {
		reader, _, err = s.blobStore.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		reader, _, err = s.blobStore.Get(ctx, srcBucket, srcKey)
	}
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return s.PutWithVersioning(ctx, dstBucket, dstKey, reader, contentType, metadata, false, StorageClassStandard, nil)
}
