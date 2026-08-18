package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

// PutEncrypted stores encrypted data for an object.
func (s *ObjectStore) PutEncrypted(ctx context.Context, bucket, key string, encryptedData []byte, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
	return s.PutEncryptedWithVersioning(ctx, bucket, key, encryptedData, contentType, metadata, sseMetadata, false, storageClass, sysMeta)
}

// PutEncryptedWithVersioning stores encrypted data for an object with versioning support.
func (s *ObjectStore) PutEncryptedWithVersioning(ctx context.Context, bucket, key string, encryptedData []byte, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
	versionId := "null"
	isVersioned := s.isVersioningEnabled(bucket)

	if isVersioned {
		versionId = s.generateVersionId()
	}

	blobMeta := &storage.BlobMetadata{
		ContentType:   contentType,
		CustomHeaders: metadata,
	}

	var blobMetaResult *storage.BlobMetadata
	var err error
	if !isDeleteMarker && encryptedData != nil {
		reader := io.NopCloser(bytes.NewReader(encryptedData))
		if isVersioned {
			blobMetaResult, err = s.blobStore.PutWithVersion(ctx, bucket, key, versionId, reader, blobMeta)
		} else {
			blobMetaResult, err = s.blobStore.Put(ctx, bucket, key, reader, blobMeta)
		}
		if err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC()
	size := int64(len(encryptedData))
	if sseMetadata != nil && sseMetadata.UnencryptedSize > 0 {
		size = sseMetadata.UnencryptedSize
	}
	sc := storageClass
	if sc == "" {
		sc = StorageClassStandard
	}
	obj := &Object{
		Key:            key,
		BucketName:     bucket,
		Size:           size,
		ETag:           "",
		LastModified:   now,
		ContentType:    contentType,
		Metadata:       metadata,
		StorageClass:   sc,
		IsLatest:       true,
		IsDeleteMarker: isDeleteMarker,
		VersionID:      versionId,
		SSEMetadata:    sseMetadata,
	}
	if sysMeta != nil {
		obj.ContentEncoding = sysMeta.ContentEncoding
		obj.ContentLanguage = sysMeta.ContentLanguage
		obj.ContentDisposition = sysMeta.ContentDisposition
		obj.CacheControl = sysMeta.CacheControl
	}

	if blobMetaResult != nil {
		obj.ETag = blobMetaResult.ETag
		obj.LastModified = blobMetaResult.LastModified
	}

	if sseMetadata != nil {
		obj.ServerSideEncryption = string(sseMetadata.EncryptionType)
		if sseMetadata.KMSKeyID != "" {
			obj.SSEKMSKeyID = sseMetadata.KMSKeyID
		}
	}

	if isVersioned {
		if err := s.putVersionedObject(bucket, key, versionId, obj); err != nil {
			return nil, err
		}
	} else {
		storageKey := s.versionedStorageKey(bucket, key, "null")
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// GetEncrypted retrieves encrypted data for an object.
func (s *ObjectStore) GetEncrypted(ctx context.Context, bucket, key, versionId string) ([]byte, *Object, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		return nil, nil, err
	}

	obj := ProtoToObject(pbObj)

	var reader io.ReadCloser
	var blobMeta *storage.BlobMetadata

	if s.isVersioningEnabled(bucket) {
		blobVersionId := obj.VersionID
		if blobVersionId == "" {
			blobVersionId = "null"
		}
		reader, blobMeta, err = s.blobStore.GetWithVersion(ctx, bucket, key, blobVersionId)
	} else {
		reader, blobMeta, err = s.blobStore.Get(ctx, bucket, key)
	}
	if err != nil {
		return nil, nil, ErrObjectNotFound
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read encrypted data: %w", err)
	}

	obj.Size = blobMeta.Size
	obj.ETag = blobMeta.ETag
	obj.LastModified = blobMeta.LastModified
	obj.ContentType = blobMeta.ContentType
	obj.Metadata = blobMeta.CustomHeaders

	return data, obj, nil
}

// UpdateObjectEncryption rewrites the stored ciphertext and SSE metadata of
// an existing object version in place. The blob is rewritten carrying the
// original ETag so the update stays invisible to conditional requests, and
// every other object field (content metadata, timestamps, storage class,
// tags, ACL, lock state, version identifier) is preserved unchanged. No new
// version is created; this is the storage counterpart of the
// UpdateObjectEncryption API.
func (s *ObjectStore) UpdateObjectEncryption(ctx context.Context, bucket, key, versionId string, encryptedData []byte, sseMetadata *SSEObjectMetadata) (*Object, error) {
	isVersioned := s.isVersioningEnabled(bucket)
	effectiveVersionId := versionId
	if !isVersioned && versionId == "null" {
		effectiveVersionId = ""
	}

	var pbObj pb.Object
	switch {
	case effectiveVersionId != "":
		if err := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, effectiveVersionId), &pbObj); err != nil {
			return nil, ErrObjectNotFound
		}
	case isVersioned:
		if err := s.BaseStore.GetProto(s.latestKeyStorageKey(bucket, key), &pbObj); err != nil {
			if err2 := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, "null"), &pbObj); err2 != nil {
				return nil, ErrObjectNotFound
			}
		}
	default:
		if err := s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, "null"), &pbObj); err != nil {
			return nil, ErrObjectNotFound
		}
	}

	obj := ProtoToObject(&pbObj)
	if obj.IsDeleteMarker {
		return nil, ErrObjectNotFound
	}

	blobVersionId := obj.VersionID
	if blobVersionId == "" {
		blobVersionId = "null"
	}
	blobMeta := &storage.BlobMetadata{
		ContentType:   obj.ContentType,
		CustomHeaders: obj.Metadata,
		ETag:          obj.ETag,
		LastModified:  obj.LastModified,
	}
	reader := io.NopCloser(bytes.NewReader(encryptedData))
	var err error
	if isVersioned {
		_, err = s.blobStore.PutWithVersion(ctx, bucket, key, blobVersionId, reader, blobMeta)
	} else {
		_, err = s.blobStore.Put(ctx, bucket, key, reader, blobMeta)
	}
	if err != nil {
		return nil, err
	}

	obj.SSEMetadata = sseMetadata
	obj.ServerSideEncryption = string(sseMetadata.EncryptionType)
	obj.SSEKMSKeyID = sseMetadata.KMSKeyID
	if sseMetadata.UnencryptedSize > 0 {
		obj.Size = sseMetadata.UnencryptedSize
	}

	if isVersioned {
		var latest pb.Object
		if s.BaseStore.GetProto(s.latestKeyStorageKey(bucket, key), &latest) == nil && latest.VersionId == obj.VersionID {
			if err := s.putVersionedObject(bucket, key, obj.VersionID, obj); err != nil {
				return nil, err
			}
		} else if err := s.BaseStore.PutProto(s.versionedStorageKey(bucket, key, obj.VersionID), ObjectToProto(obj)); err != nil {
			return nil, err
		}
	} else {
		if err := s.BaseStore.PutProto(s.versionedStorageKey(bucket, key, "null"), ObjectToProto(obj)); err != nil {
			return nil, err
		}
	}

	return obj, nil
}
