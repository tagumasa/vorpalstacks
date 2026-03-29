package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"

	"google.golang.org/protobuf/proto"
)

// CreateMultipartUpload initiates a multipart upload for an object.
func (s *ObjectStore) CreateMultipartUpload(ctx context.Context, bucket, key string, contentType string, metadata map[string]string, sseType SSEType, kmsKeyID, customerKeyMD5 string, sseMetadata *SSEObjectMetadata, plaintextDataKey []byte) (*MultipartUpload, error) {
	blobMeta := &storage.BlobMetadata{
		ContentType:   contentType,
		CustomHeaders: metadata,
	}

	uploadId, err := s.blobStore.CreateMultipartUpload(ctx, bucket, key, blobMeta)
	if err != nil {
		return nil, err
	}

	upload := &MultipartUpload{
		UploadID:         uploadId,
		Key:              key,
		BucketName:       bucket,
		Initiated:        time.Now().UTC(),
		ContentType:      contentType,
		Metadata:         metadata,
		SSEType:          sseType,
		KMSKeyID:         kmsKeyID,
		CustomerKeyMD5:   customerKeyMD5,
		SSEMetadata:      sseMetadata,
		PlaintextDataKey: plaintextDataKey,
	}

	data, err := proto.Marshal(MultipartUploadToProto(upload))
	if err != nil {
		return nil, err
	}

	if err := s.storage.Bucket(multipartBucketName(s.region)).Put([]byte(s.multipartKey(uploadId)), data); err != nil {
		return nil, err
	}

	indexKey := s.multipartIndexKey(bucket, key, uploadId)
	if err := s.storage.Bucket(multipartIndexBucketName(s.region)).Put([]byte(indexKey), []byte{}); err != nil {
		return nil, err
	}

	return upload, nil
}

// GetMultipartUpload retrieves a multipart upload by its upload ID.
func (s *ObjectStore) GetMultipartUpload(uploadId string) (*MultipartUpload, error) {
	data, err := s.storage.Bucket(multipartBucketName(s.region)).Get([]byte(s.multipartKey(uploadId)))
	if err != nil || data == nil {
		return nil, ErrUploadNotFound
	}

	var pbUpload pb.MultipartUpload
	if err := proto.Unmarshal(data, &pbUpload); err != nil {
		return nil, err
	}

	return ProtoToMultipartUpload(&pbUpload), nil
}

// UploadPart uploads a part of a multipart upload.
func (s *ObjectStore) UploadPart(ctx context.Context, bucket, key, uploadId string, partNumber int, reader io.Reader, encryptedSize int64, contentNonce, dataKey []byte) (*ObjectPart, error) {
	lockKey := "multipart#" + uploadId
	mutex := s.getVersionLock(lockKey)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.versionMutex.Delete(lockKey)
	}()

	upload, err := s.GetMultipartUpload(uploadId)
	if err != nil {
		return nil, err
	}

	if upload.BucketName != bucket || upload.Key != key {
		return nil, ErrUploadNotFound
	}

	etag, err := s.blobStore.UploadPart(ctx, bucket, key, uploadId, partNumber, reader)
	if err != nil {
		return nil, err
	}

	part := &ObjectPart{
		PartNumber:    partNumber,
		ETag:          etag,
		LastModified:  time.Now().UTC(),
		EncryptedSize: encryptedSize,
		ContentNonce:  contentNonce,
		DataKey:       dataKey,
	}

	upload.AddPart(partNumber, *part)

	data, err := proto.Marshal(MultipartUploadToProto(upload))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal upload metadata: %w", err)
	}
	if err := s.storage.Bucket(multipartBucketName(s.region)).Put([]byte(s.multipartKey(uploadId)), data); err != nil {
		return nil, fmt.Errorf("failed to save part metadata: %w", err)
	}

	return part, nil
}

// ListParts lists the parts of a multipart upload.
func (s *ObjectStore) ListParts(ctx context.Context, bucket, key, uploadId string, partNumberMarker int, maxParts int) ([]ObjectPart, int, bool, error) {
	if _, err := s.GetMultipartUpload(uploadId); err != nil {
		return nil, 0, false, fmt.Errorf("upload not found: %w", err)
	}

	parts, err := s.blobStore.ListParts(ctx, bucket, key, uploadId)
	if err != nil {
		return nil, 0, false, err
	}

	result := make([]ObjectPart, 0)
	nextPartNumberMarker := 0
	isTruncated := false

	for _, p := range parts {
		if p.PartNumber <= partNumberMarker {
			continue
		}
		result = append(result, ObjectPart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
			Size:       p.Size,
		})
		if len(result) >= maxParts {
			if len(parts) > len(result)+countSkipped(parts, partNumberMarker) {
				isTruncated = true
				nextPartNumberMarker = result[len(result)-1].PartNumber
			}
			break
		}
	}

	return result, nextPartNumberMarker, isTruncated, nil
}

func countSkipped(parts []storage.PartInfo, partNumberMarker int) int {
	count := 0
	for _, p := range parts {
		if p.PartNumber <= partNumberMarker {
			count++
		}
	}
	return count
}

// CompleteMultipartUpload completes a multipart upload by assembling the parts.
func (s *ObjectStore) CompleteMultipartUpload(ctx context.Context, bucket, key, uploadId string, parts []ObjectPart) (*Object, error) {
	upload, err := s.GetMultipartUpload(uploadId)
	if err != nil {
		return nil, err
	}

	if upload.BucketName != bucket || upload.Key != key {
		return nil, ErrUploadNotFound
	}

	var blobParts []storage.PartInfo
	for _, p := range parts {
		blobParts = append(blobParts, storage.PartInfo{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
		})
	}

	blobMeta, err := s.blobStore.CompleteMultipartUpload(ctx, bucket, key, uploadId, blobParts)
	if err != nil {
		return nil, err
	}

	if err := s.AbortMultipartUpload(ctx, bucket, key, uploadId); err != nil {
		logs.Error("Failed to cleanup multipart upload after complete", logs.Err(err))
	}

	versionId := "null"
	if s.isVersioningEnabled(bucket) {
		versionId = s.generateVersionId()
	}

	obj := newObject(key, bucket, upload.ContentType, upload.Metadata, versionId, false)
	obj.Size = blobMeta.Size
	obj.ETag = blobMeta.ETag
	obj.LastModified = blobMeta.LastModified

	if upload.SSEType != "" {
		var sseMetadata *SSEObjectMetadata
		if upload.SSEMetadata != nil && upload.SSEMetadata.EncryptedDataKey != nil {
			sseMetadata = upload.SSEMetadata
		} else if len(upload.Parts) > 0 && upload.Parts[0].DataKey != nil {
			sseMetadata = &SSEObjectMetadata{
				EncryptionType:   upload.SSEType,
				EncryptedDataKey: upload.Parts[0].DataKey,
				ContentNonce:     upload.Parts[0].ContentNonce,
				KMSKeyID:         upload.KMSKeyID,
			}
		} else {
			sseMetadata = &SSEObjectMetadata{
				EncryptionType: upload.SSEType,
				KMSKeyID:       upload.KMSKeyID,
			}
		}
		obj.SSEMetadata = sseMetadata
	}

	if s.isVersioningEnabled(bucket) {
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

// AbortMultipartUpload aborts a multipart upload.
func (s *ObjectStore) AbortMultipartUpload(ctx context.Context, bucket, key, uploadId string) error {
	if _, err := s.GetMultipartUpload(uploadId); err != nil {
		return fmt.Errorf("upload not found: %w", err)
	}

	err1 := s.blobStore.AbortMultipartUpload(ctx, bucket, key, uploadId)
	err2 := s.storage.Bucket(multipartBucketName(s.region)).Delete([]byte(s.multipartKey(uploadId)))
	err3 := s.storage.Bucket(multipartIndexBucketName(s.region)).Delete([]byte(s.multipartIndexKey(bucket, key, uploadId)))
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return err3
}

// ListMultipartUploads lists the in-progress multipart uploads for a bucket.
func (s *ObjectStore) ListMultipartUploads(bucket, prefix, keyMarker, uploadIdMarker string, maxUploads int) (*MultipartUploadListResult, error) {
	if maxUploads <= 0 {
		maxUploads = 1000
	}

	var uploads []*MultipartUpload
	count := 0
	started := keyMarker == "" && uploadIdMarker == ""
	hasMore := false

	indexPrefix := bucket + "#"
	iter := s.storage.Bucket(multipartIndexBucketName(s.region)).ScanPrefix([]byte(indexPrefix))
	defer iter.Close()

	for iter.Next() {
		indexKey := string(iter.Key())

		parts := strings.SplitN(indexKey, "#", 4)
		if len(parts) < 3 {
			continue
		}
		key := parts[1]
		uploadId := parts[2]

		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}

		if !started {
			if key == keyMarker && (uploadIdMarker == "" || uploadId == uploadIdMarker) {
				started = true
			}
			continue
		}

		upload, err := s.GetMultipartUpload(uploadId)
		if err != nil {
			continue
		}

		if count < maxUploads {
			uploads = append(uploads, upload)
			count++
		} else {
			hasMore = true
			break
		}
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	result := &MultipartUploadListResult{
		Uploads:     uploads,
		IsTruncated: hasMore,
	}

	if len(uploads) > 0 {
		result.NextKeyMarker = uploads[len(uploads)-1].Key
		result.NextUploadIDMarker = uploads[len(uploads)-1].UploadID
	}

	return result, nil
}
