package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"vorpalstacks/internal/core/storage"
)

// PutEncrypted stores encrypted data for an object.
func (s *ObjectStore) PutEncrypted(ctx context.Context, bucket, key string, encryptedData []byte, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
	return s.PutEncryptedWithVersioning(ctx, bucket, key, encryptedData, contentType, metadata, sseMetadata, false, storageClass, sysMeta)
}

// PutEncryptedWithVersioning stores encrypted data with
// versioning support, through the same put path as a plain write: the
// record carries the SSE fields, and the shared persistence tail keeps the
// null record, its blob slot and any "_latest" pointer consistent.
func (s *ObjectStore) PutEncryptedWithVersioning(ctx context.Context, bucket, key string, encryptedData []byte, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
	var reader io.Reader
	if encryptedData != nil {
		reader = io.NopCloser(bytes.NewReader(encryptedData))
	}
	return s.putObjectVersioned(ctx, bucket, key, reader, contentType, metadata, isDeleteMarker, storageClass, sysMeta, sseMetadata)
}

// PutEncryptedStreaming stores an encrypted object from an encrypting
// reader: the blob write consumes the stream, and the SSE metadata the
// stream accumulated becomes the record's SSE fields. sseMetadata is
// filled by the caller's reader as it is drained; putObjectVersioned
// reads it after the blob write, when the digest, unencrypted size and
// part table are final — so the whole encrypted object is never held in
// memory, only the reader's bounded chunk.
func (s *ObjectStore) PutEncryptedStreaming(ctx context.Context, bucket, key string, encryptedReader io.Reader, contentType string, metadata map[string]string, sseMetadata *SSEObjectMetadata, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
	return s.putObjectVersioned(ctx, bucket, key, encryptedReader, contentType, metadata, false, storageClass, sysMeta, sseMetadata)
}

// GetEncrypted retrieves encrypted data for an object.
func (s *ObjectStore) GetEncrypted(ctx context.Context, bucket, key, versionId string) ([]byte, *Object, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		return nil, nil, err
	}

	obj := ProtoToObject(pbObj)

	// The blob tier mirrors the record's layout (see readRecordBlob); the
	// null lookup falls back across both blob locations.
	reader, blobMeta, err := s.readRecordBlob(ctx, bucket, key, obj.VersionID)
	if err != nil {
		return nil, nil, blobReadErr(err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read encrypted data: %w", err)
	}

	applyBlobFacts(obj, blobMeta)

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
	var updated *Object
	err := s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		// The blob mirrors the record's layout: a real version writes the
		// versioned blob, and the null version writes the null-versioned
		// slot only when the key carries versioned layout (this update is
		// the one writer of a null-versioned blob); otherwise the plain
		// slot is the null blob's home.
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
		if blobVersionId != "null" || s.BaseStore.Exists(s.latestKeyStorageKey(bucket, key)) {
			_, err = s.blobStore.PutWithVersion(ctx, bucket, key, blobVersionId, reader, blobMeta)
		} else {
			_, err = s.blobStore.Put(ctx, bucket, key, reader, blobMeta)
		}
		if err != nil {
			return err
		}

		obj.SSEMetadata = sseMetadata
		obj.ServerSideEncryption = string(sseMetadata.EncryptionType)
		obj.SSEKMSKeyID = sseMetadata.KMSKeyID
		if sseMetadata.UnencryptedSize > 0 {
			obj.Size = sseMetadata.UnencryptedSize
		}
		updated = obj
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}
