package s3

import (
	"context"
	"fmt"
	"io"
	"strings"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"

	"google.golang.org/protobuf/proto"
)

func (s *ObjectStore) putVersionedObject(bucket, key, versionId string, obj *Object) error {
	lockKey := bucket + "#" + key
	s.keyLocker.Lock(lockKey)
	defer func() {
		s.keyLocker.Unlock(lockKey)
		s.keyLocker.Delete(lockKey)
	}()

	latestKey := s.latestKeyStorageKey(bucket, key)
	var prevLatest pb.Object
	if err := s.BaseStore.GetProto(latestKey, &prevLatest); err == nil {
		prevLatest.IsLatest = false
		prevStorageKey := s.versionedStorageKey(bucket, key, prevLatest.VersionId)
		if err := s.BaseStore.PutProto(prevStorageKey, &prevLatest); err != nil {
			return fmt.Errorf("failed to update previous version: %w", err)
		}
	}

	versionedKey := s.versionedStorageKey(bucket, key, versionId)
	if err := s.BaseStore.PutProto(versionedKey, ObjectToProto(obj)); err != nil {
		return err
	}

	return s.BaseStore.PutProto(latestKey, ObjectToProto(obj))
}

// GetWithVersion retrieves a specific version of an object.
func (s *ObjectStore) GetWithVersion(ctx context.Context, bucket, key, versionId string) (io.ReadCloser, *Object, error) {
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

	obj.Size = blobMeta.Size
	obj.ETag = blobMeta.ETag
	obj.LastModified = blobMeta.LastModified
	obj.ContentType = blobMeta.ContentType
	obj.Metadata = blobMeta.CustomHeaders

	return reader, obj, nil
}

// PutWithVersioning stores an object with versioning support.
func (s *ObjectStore) PutWithVersioning(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
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
	if !isDeleteMarker && reader != nil {
		if isVersioned {
			blobMetaResult, err = s.blobStore.Put(ctx, bucket, key+"#"+versionId, reader, blobMeta)
		} else {
			blobMetaResult, err = s.blobStore.Put(ctx, bucket, key, reader, blobMeta)
		}
		if err != nil {
			return nil, err
		}
	}

	obj := newObject(key, bucket, contentType, metadata, versionId, isDeleteMarker, storageClass, sysMeta)

	if blobMetaResult != nil {
		obj.Size = blobMetaResult.Size
		obj.ETag = blobMetaResult.ETag
		obj.LastModified = blobMetaResult.LastModified
	}

	if isVersioned {
		if err := s.putVersionedObject(bucket, key, versionId, obj); err != nil {
			return nil, err
		}
	} else {
		lockKey := bucket + "#" + key
		s.keyLocker.Lock(lockKey)
		defer func() {
			s.keyLocker.Unlock(lockKey)
			s.keyLocker.Delete(lockKey)
		}()

		storageKey := s.storageKey(bucket, key)
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// DeleteWithVersion deletes a specific version of an object.
// If the version ID is provided, it deletes that specific version.
// If the object is versioned and the deleted version is the latest, it updates the latest pointer.
func (s *ObjectStore) DeleteWithVersion(ctx context.Context, bucket, key, versionId string) (*Object, error) {
	if versionId != "" {
		if s.isVersioningEnabled(bucket) {
			lockKey := bucket + "#" + key
			s.keyLocker.Lock(lockKey)
			defer func() {
				s.keyLocker.Unlock(lockKey)
				s.keyLocker.Delete(lockKey)
			}()

			latestKey := s.latestKeyStorageKey(bucket, key)
			var latestObj pb.Object
			isLatest := false
			if err := s.BaseStore.GetProto(latestKey, &latestObj); err == nil {
				isLatest = (latestObj.VersionId == versionId)
			}

			if err := s.blobStore.DeleteWithVersion(ctx, bucket, key, versionId); err != nil {
				return nil, err
			}

			storageKey := s.versionedStorageKey(bucket, key, versionId)
			if err := s.BaseStore.Delete(storageKey); err != nil {
				return nil, err
			}

			if isLatest {
				prefix := bucket + "#" + key + "#"
				var remainingVersions []*Object
				err := s.ScanPrefix(prefix, func(k string, v []byte) error {
					if strings.HasSuffix(k, "#_latest") {
						return nil
					}
					var pbObj pb.Object
					if err := proto.Unmarshal(v, &pbObj); err != nil {
						return err
					}
					remainingVersions = append(remainingVersions, ProtoToObject(&pbObj))
					return nil
				})
				if err != nil {
					return nil, err
				}

				if len(remainingVersions) == 0 {
					if err := s.BaseStore.Delete(latestKey); err != nil {
						return nil, err
					}
				} else {
					var newLatest *Object
					for _, v := range remainingVersions {
						if newLatest == nil || v.LastModified.After(newLatest.LastModified) {
							newLatest = v
						}
					}
					if newLatest != nil {
						newLatest.IsLatest = true
						newLatestKey := s.versionedStorageKey(bucket, key, newLatest.VersionID)
						if err := s.BaseStore.PutProto(newLatestKey, ObjectToProto(newLatest)); err != nil {
							return nil, err
						}
						if err := s.BaseStore.PutProto(latestKey, ObjectToProto(newLatest)); err != nil {
							return nil, err
						}
					}
				}
			}
			return nil, nil
		} else {
			if err := s.blobStore.Delete(ctx, bucket, key); err != nil {
				return nil, err
			}
			storageKey := s.versionedStorageKey(bucket, key, versionId)
			if err := s.BaseStore.Delete(storageKey); err != nil {
				return nil, err
			}
			return nil, nil
		}
	}

	if s.isVersioningEnabled(bucket) {
		deleteMarker, err := s.PutWithVersioning(ctx, bucket, key, nil, "", nil, true, StorageClassStandard, nil)
		if err != nil {
			return nil, err
		}
		return deleteMarker, nil
	}

	if err := s.blobStore.Delete(ctx, bucket, key); err != nil {
		return nil, err
	}
	storageKey := s.storageKey(bucket, key)
	if err := s.BaseStore.Delete(storageKey); err != nil {
		return nil, err
	}
	return nil, nil
}

// HeadWithVersion retrieves metadata for a specific version of an object.
func (s *ObjectStore) HeadWithVersion(ctx context.Context, bucket, key, versionId string) (*Object, error) {
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
			nullKey := s.versionedStorageKey(bucket, key, "null")
			if err = s.BaseStore.GetProto(nullKey, &pbObj); err != nil {
				fallbackKey := s.storageKey(bucket, key)
				if err = s.BaseStore.GetProto(fallbackKey, &pbObj); err != nil {
					return nil, ErrObjectNotFound
				}
			}
		}
	} else {
		storageKey := s.storageKey(bucket, key)
		if err = s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			return nil, ErrObjectNotFound
		}
	}

	obj := ProtoToObject(&pbObj)
	if obj.IsDeleteMarker {
		return nil, ErrObjectNotFound
	}

	blobVersionId := obj.VersionID
	if blobVersionId == "null" {
		blobVersionId = ""
	}

	var blobMeta *storage.BlobMetadata
	if blobVersionId != "" && s.isVersioningEnabled(bucket) {
		reader, meta, err := s.blobStore.GetWithVersion(ctx, bucket, key, blobVersionId)
		if err != nil {
			return nil, ErrObjectNotFound
		}
		reader.Close()
		blobMeta = meta
	} else {
		blobMeta, err = s.blobStore.Head(ctx, bucket, key)
		if err != nil {
			return nil, ErrObjectNotFound
		}
	}

	obj.Size = blobMeta.Size
	obj.ETag = blobMeta.ETag
	obj.LastModified = blobMeta.LastModified
	obj.ContentType = blobMeta.ContentType
	obj.Metadata = blobMeta.CustomHeaders

	return obj, nil
}

// GetRangeWithVersion retrieves a range of bytes from a specific version of an object.
func (s *ObjectStore) GetRangeWithVersion(ctx context.Context, bucket, key, versionId string, offset, length int64) (io.ReadCloser, *Object, error) {
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
		reader, blobMeta, err = s.blobStore.GetRangeWithVersion(ctx, bucket, key, blobVersionId, offset, length)
	} else {
		reader, blobMeta, err = s.blobStore.GetRange(ctx, bucket, key, offset, length)
	}
	if err != nil {
		return nil, nil, ErrObjectNotFound
	}

	obj.Size = blobMeta.Size
	obj.ETag = blobMeta.ETag
	obj.LastModified = blobMeta.LastModified
	obj.ContentType = blobMeta.ContentType
	obj.Metadata = blobMeta.CustomHeaders

	return reader, obj, nil
}

// SetACLWithVersion sets the access control list for a specific version of an object.
func (s *ObjectStore) SetACLWithVersion(bucket, key, versionId string, acp *AccessControlPolicy) error {
	lockKey := bucket + "#" + key
	s.keyLocker.Lock(lockKey)
	defer func() {
		s.keyLocker.Unlock(lockKey)
		s.keyLocker.Delete(lockKey)
	}()

	var pbObj pb.Object
	var storageKey string

	isVersioned := s.isVersioningEnabled(bucket)
	effectiveVersionId := versionId
	if !isVersioned && versionId == "null" {
		effectiveVersionId = ""
	}

	if effectiveVersionId != "" || isVersioned {
		storageKey = s.versionedStorageKey(bucket, key, effectiveVersionId)
		if effectiveVersionId == "" {
			storageKey = s.latestKeyStorageKey(bucket, key)
		}
	} else {
		storageKey = s.storageKey(bucket, key)
	}

	if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
		return ErrObjectNotFound
	}

	obj := ProtoToObject(&pbObj)
	obj.ACL = acp
	return s.BaseStore.PutProto(storageKey, ObjectToProto(obj))
}

// GetACLWithVersion retrieves the access control list for a specific version of an object.
func (s *ObjectStore) GetACLWithVersion(bucket, key, versionId string) (*AccessControlPolicy, error) {
	var pbObj pb.Object
	var storageKey string

	isVersioned := s.isVersioningEnabled(bucket)
	effectiveVersionId := versionId
	if !isVersioned && versionId == "null" {
		effectiveVersionId = ""
	}

	if effectiveVersionId != "" || isVersioned {
		storageKey = s.versionedStorageKey(bucket, key, effectiveVersionId)
		if effectiveVersionId == "" {
			storageKey = s.latestKeyStorageKey(bucket, key)
		}
	} else {
		storageKey = s.storageKey(bucket, key)
	}

	if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
		return nil, ErrObjectNotFound
	}

	obj := ProtoToObject(&pbObj)
	return obj.ACL, nil
}
