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
			blobMetaResult, err = s.blobStore.Put(ctx, bucket, key+"#"+versionId, reader, blobMeta)
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
		storageKey := s.storageKey(bucket, key)
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// GetEncrypted retrieves encrypted data for an object.
func (s *ObjectStore) GetEncrypted(ctx context.Context, bucket, key, versionId string) ([]byte, *Object, error) {
	var pbObj pb.Object
	var err error

	isVersioned := s.isVersioningEnabled(bucket)
	effectiveVersionId := versionId
	if !isVersioned && versionId == "null" {
		effectiveVersionId = ""
	}

	if effectiveVersionId != "" || isVersioned {
		storageKey := s.versionedStorageKey(bucket, key, effectiveVersionId)
		if effectiveVersionId == "" {
			storageKey = s.latestKeyStorageKey(bucket, key)
		}
		if err = s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			return nil, nil, ErrObjectNotFound
		}
	} else {
		storageKey := s.storageKey(bucket, key)
		if err = s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			return nil, nil, ErrObjectNotFound
		}
	}

	obj := ProtoToObject(&pbObj)
	if obj.IsDeleteMarker {
		return nil, obj, ErrObjectNotFound
	}

	var reader io.ReadCloser
	var blobMeta *storage.BlobMetadata

	if isVersioned {
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
