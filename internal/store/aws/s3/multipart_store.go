package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"

	"google.golang.org/protobuf/proto"
)

// CreateMultipartUpload initiates a multipart upload for an object.
func (s *ObjectStore) CreateMultipartUpload(ctx context.Context, bucket, key string, contentType string, metadata map[string]string, sseType SSEType, kmsKeyID, customerKeyMD5 string, sseMetadata *SSEObjectMetadata, plaintextDataKey []byte, storageClass ObjectStorageClass, acl *AccessControlPolicy) (*MultipartUpload, error) {
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
		StorageClass:     storageClass,
		ContentType:      contentType,
		Metadata:         metadata,
		SSEType:          sseType,
		KMSKeyID:         kmsKeyID,
		CustomerKeyMD5:   customerKeyMD5,
		SSEMetadata:      sseMetadata,
		PlaintextDataKey: plaintextDataKey,
		ACL:              acl,
	}

	data, err := proto.Marshal(MultipartUploadToProto(upload))
	if err != nil {
		s.blobStore.AbortMultipartUpload(ctx, bucket, key, uploadId)
		return nil, err
	}

	if err := s.storage.Bucket(multipartBucketName(s.region)).Put([]byte(uploadId), data); err != nil {
		s.blobStore.AbortMultipartUpload(ctx, bucket, key, uploadId)
		return nil, err
	}

	indexKey := s.multipartIndexKey(bucket, key, uploadId)
	if err := s.storage.Bucket(multipartIndexBucketName(s.region)).Put([]byte(indexKey), []byte{}); err != nil {
		s.storage.Bucket(multipartBucketName(s.region)).Delete([]byte(uploadId))
		s.blobStore.AbortMultipartUpload(ctx, bucket, key, uploadId)
		return nil, err
	}

	return upload, nil
}

// GetMultipartUpload retrieves a multipart upload by its upload ID.
func (s *ObjectStore) GetMultipartUpload(uploadId string) (*MultipartUpload, error) {
	data, err := s.storage.Bucket(multipartBucketName(s.region)).Get([]byte(uploadId))
	if err != nil {
		// A read failure is infrastructure trouble, not a missing upload;
		// surfacing it as not-found would mask it behind 404s.
		return nil, err
	}
	if data == nil {
		return nil, ErrUploadNotFound
	}

	var pbUpload pb.MultipartUpload
	if err := proto.Unmarshal(data, &pbUpload); err != nil {
		return nil, err
	}

	return ProtoToMultipartUpload(&pbUpload), nil
}

// UploadPart uploads a part of a multipart upload.
func (s *ObjectStore) UploadPart(ctx context.Context, bucket, key, uploadId string, partNumber int, reader io.Reader, encryptedSize int64, plainSize int64, contentNonce, dataKey []byte) (*ObjectPart, error) {
	lockKey := "multipart#" + uploadId
	s.keyLocker.Lock(lockKey)
	defer s.keyLocker.Unlock(lockKey)

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
		Size:          plainSize,
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
	if err := s.storage.Bucket(multipartBucketName(s.region)).Put([]byte(uploadId), data); err != nil {
		return nil, fmt.Errorf("failed to save part metadata: %w", err)
	}

	return part, nil
}

// ListParts lists the parts of a multipart upload and returns them together
// with the upload record they were read from.
func (s *ObjectStore) ListParts(ctx context.Context, bucket, key, uploadId string, partNumberMarker int, maxParts int) (*ListPartsResult, error) {
	// The record is resolved before the page-limit branch so a missing or
	// misaddressed upload reports ErrUploadNotFound whatever the page size.
	upload, err := s.GetMultipartUpload(uploadId)
	if err != nil {
		return nil, err
	}

	if upload.BucketName != bucket || upload.Key != key {
		return nil, ErrUploadNotFound
	}

	// Callers resolve the default page size; a limit of zero means an
	// empty, non-truncated page and only negative values are clamped here.
	if maxParts <= 0 {
		return &ListPartsResult{Upload: upload}, nil
	}

	if len(upload.Parts) == 0 {
		return s.listPartsFromBlob(ctx, bucket, key, uploadId, upload, partNumberMarker, maxParts)
	}

	parts, nextPartNumberMarker, isTruncated := listPartsFromUpload(upload.Parts, partNumberMarker, maxParts)
	return &ListPartsResult{
		Upload:               upload,
		Parts:                parts,
		NextPartNumberMarker: nextPartNumberMarker,
		IsTruncated:          isTruncated,
	}, nil
}

func (s *ObjectStore) listPartsFromBlob(ctx context.Context, bucket, key, uploadId string, upload *MultipartUpload, partNumberMarker int, maxParts int) (*ListPartsResult, error) {
	parts, err := s.blobStore.ListParts(ctx, bucket, key, uploadId)
	if err != nil {
		return nil, err
	}

	result := &ListPartsResult{Upload: upload, Parts: make([]ObjectPart, 0)}
	skipped := 0

	for _, p := range parts {
		if p.PartNumber <= partNumberMarker {
			skipped++
			continue
		}
		result.Parts = append(result.Parts, ObjectPart{
			PartNumber: p.PartNumber,
			ETag:       p.ETag,
			Size:       p.Size,
		})
		if len(result.Parts) >= maxParts {
			if len(parts) > len(result.Parts)+skipped {
				result.IsTruncated = true
				result.NextPartNumberMarker = result.Parts[len(result.Parts)-1].PartNumber
			}
			break
		}
	}

	return result, nil
}

func listPartsFromUpload(parts []ObjectPart, partNumberMarker int, maxParts int) ([]ObjectPart, int, bool) {
	sorted := make([]ObjectPart, len(parts))
	copy(sorted, parts)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].PartNumber < sorted[j].PartNumber
	})

	result := make([]ObjectPart, 0, maxParts)
	nextPartNumberMarker := 0

	for _, p := range sorted {
		if p.PartNumber <= partNumberMarker {
			continue
		}
		if len(result) >= maxParts {
			nextPartNumberMarker = result[len(result)-1].PartNumber
			return result, nextPartNumberMarker, true
		}
		result = append(result, p)
	}

	return result, 0, false
}

// CompleteMultipartUpload completes a multipart upload by assembling the parts.
func (s *ObjectStore) CompleteMultipartUpload(ctx context.Context, bucket, key, uploadId string, parts []ObjectPart) (*Object, error) {
	lockKey := "multipart#" + uploadId
	s.keyLocker.Lock(lockKey)
	defer s.keyLocker.Unlock(lockKey)

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

	versionId := "null"
	if s.isVersioningEnabled(bucket) {
		versionId = s.generateVersionId()
	}

	// The blob tier derives the storage locations from bucket/key/versionId
	// itself (plain address for the null version, versioned address
	// otherwise), so the completed object lands exactly where the
	// WithVersion reads look.
	blobMeta, err := s.blobStore.CompleteMultipartUpload(ctx, bucket, key, versionId, uploadId, blobParts)
	if err != nil {
		return nil, err
	}

	obj := newObject(key, bucket, upload.ContentType, upload.Metadata, versionId, false, upload.StorageClass, nil)
	applyBlobFacts(obj, blobMeta)

	// Keep the plain-size part boundaries so partNumber reads can resolve
	// individual parts of the completed object.
	sortedParts := make([]ObjectPart, len(upload.Parts))
	copy(sortedParts, upload.Parts)
	sort.Slice(sortedParts, func(i, j int) bool { return sortedParts[i].PartNumber < sortedParts[j].PartNumber })
	for _, p := range sortedParts {
		obj.Parts = append(obj.Parts, ObjectPartBoundary{PartNumber: p.PartNumber, Size: p.Size})
	}

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

		var partInfos []PartEncryptionInfo
		for _, p := range upload.Parts {
			if p.EncryptedSize > 0 || p.DataKey != nil {
				partInfos = append(partInfos, PartEncryptionInfo{
					EncryptedSize: p.EncryptedSize,
					PlainSize:     p.Size,
					ContentNonce:  p.ContentNonce,
					DataKey:       p.DataKey,
				})
			}
		}
		sseMetadata.PartEncryptionInfos = partInfos

		var totalPlain int64
		for _, p := range upload.Parts {
			totalPlain += p.Size
		}
		sseMetadata.UnencryptedSize = totalPlain

		// The assembled blob is ciphertext; the object's reported size is
		// the plaintext length, exactly like the single-put path.
		if totalPlain > 0 {
			obj.Size = totalPlain
		}

		obj.SSEMetadata = sseMetadata
	}

	if s.isVersioningEnabled(bucket) {
		if err := s.putVersionedObject(bucket, key, versionId, obj); err != nil {
			return nil, err
		}
	} else {
		storageKey := s.versionedStorageKey(bucket, key, "null")
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return nil, err
		}
	}

	// The upload is torn down only after the record is persisted: the parts
	// stay recoverable for a client retry until the completed object itself
	// is durably stored (the blob tier's keep-parts-for-retry contract).
	if err := s.AbortMultipartUpload(ctx, bucket, key, uploadId); err != nil {
		logs.Error("Failed to cleanup multipart upload after complete", logs.Err(err))
	}

	return obj, nil
}

// AbortMultipartUpload aborts a multipart upload. A missing upload — or one
// whose bucket/key does not match the upload's own address, which AWS
// reports identically as NoSuchUpload — returns the ErrUploadNotFound
// sentinel unwrapped so every caller maps it with errors.Is.
func (s *ObjectStore) AbortMultipartUpload(ctx context.Context, bucket, key, uploadId string) error {
	upload, err := s.GetMultipartUpload(uploadId)
	if err != nil {
		return err
	}
	if upload.BucketName != bucket || upload.Key != key {
		return ErrUploadNotFound
	}

	err1 := s.blobStore.AbortMultipartUpload(ctx, bucket, key, uploadId)
	err2 := s.storage.Bucket(multipartBucketName(s.region)).Delete([]byte(uploadId))
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
	// Callers resolve the default page size; a limit of zero means an
	// empty, non-truncated page and only negative values are clamped here.
	if maxUploads < 0 {
		maxUploads = 0
	}
	if maxUploads == 0 {
		return &MultipartUploadListResult{}, nil
	}

	var uploads []*MultipartUpload
	count := 0
	started := keyMarker == "" && uploadIdMarker == ""
	hasMore := false

	indexPrefix := bucket + keySep
	prefixLen := len(indexPrefix)
	iter := s.storage.Bucket(multipartIndexBucketName(s.region)).ScanPrefix([]byte(indexPrefix))
	defer iter.Close()

	for iter.Next() {
		indexKey := string(iter.Key())

		lastSep := strings.LastIndex(indexKey, keySep)
		if lastSep <= prefixLen-1 {
			continue
		}
		key := indexKey[prefixLen:lastSep]
		uploadId := indexKey[lastSep+1:]

		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}

		// Marker semantics per the API contract: without an upload-id-marker
		// only keys lexicographically greater than the key-marker are
		// included; with one, equal-key uploads whose upload ID is
		// lexicographically greater than the upload-id-marker are included
		// too (the marker entry itself is never returned).
		if !started {
			if key > keyMarker {
				started = true
			} else if key == keyMarker && uploadIdMarker != "" && uploadId > uploadIdMarker {
				started = true
			} else {
				continue
			}
		}

		upload, err := s.GetMultipartUpload(uploadId)
		if err != nil {
			// A missing upload record makes the index entry stale: the
			// abort/complete paths remove both together, so a lone index
			// key is residue from a partial failure — reconcile it here
			// instead of skipping it on every listing forever. Any other
			// error is infrastructure trouble and stops the listing.
			if !errors.Is(err, ErrUploadNotFound) {
				return nil, err
			}
			if delErr := s.storage.Bucket(multipartIndexBucketName(s.region)).Delete([]byte(indexKey)); delErr != nil {
				logs.Warn("s3: stale multipart index entry could not be removed", logs.String("indexKey", indexKey), logs.Err(delErr))
			}
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

	// The Next markers name the first entry NOT returned and are emitted
	// only on truncation — AWS defines them for exactly that case.
	if hasMore && len(uploads) > 0 {
		result.NextKeyMarker = uploads[len(uploads)-1].Key
		result.NextUploadIDMarker = uploads[len(uploads)-1].UploadID
	}

	return result, nil
}
