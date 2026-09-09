package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_s3"

	"google.golang.org/protobuf/proto"
)

// applyBlobFacts copies the blob-authoritative fields onto the object
// record. The blob layer owns only the physical facts of the stored bytes
// (size, ETag, modification time). ContentType and user metadata are
// semantic fields of the object record itself: the multipart write path
// never supplies them to the blob layer, so copying them here would blank
// the values requested at upload on every read.
func applyBlobFacts(obj *Object, blobMeta *storage.BlobMetadata) {
	obj.Size = blobMeta.Size
	obj.ETag = blobMeta.ETag
	obj.LastModified = blobMeta.LastModified
}

// putVersionedObject writes a new object version, updates the _latest
// pointer and clears the previous latest's IsLatest flag in one atomic
// batch — a crash between the writes can never leave a version that lists
// but is not current, or two records flagged latest.
func (s *ObjectStore) putVersionedObject(bucket, key, versionId string, obj *Object) error {
	lockKey := bucket + keySep + key
	s.keyLocker.Lock(lockKey)
	defer s.keyLocker.Unlock(lockKey)

	latestKey := s.latestKeyStorageKey(bucket, key)

	var prevLatest pb.Object
	prevLatestFound := s.BaseStore.GetProto(latestKey, &prevLatest) == nil
	if !prevLatestFound {
		// Objects created before versioning was enabled have no "_latest"
		// pointer; their "null" record is the previous latest and must be
		// flipped like any other, or it lists alongside the new version.
		prevLatestFound = s.BaseStore.GetProto(s.versionedStorageKey(bucket, key, "null"), &prevLatest) == nil
	}

	batchBucket, ok := s.BaseStore.Bucket().(storage.BatchBucket)
	if !ok {
		return fmt.Errorf("s3: storage bucket does not support atomic batches")
	}
	batch := batchBucket.NewBatch()
	defer batch.Close()

	objBytes, err := proto.Marshal(ObjectToProto(obj))
	if err != nil {
		return err
	}
	if err := batch.Put([]byte(s.versionedStorageKey(bucket, key, versionId)), objBytes); err != nil {
		return err
	}
	if err := batch.Put([]byte(latestKey), objBytes); err != nil {
		return err
	}

	if prevLatestFound && prevLatest.VersionId != versionId {
		prevLatest.IsLatest = false
		prevBytes, err := proto.Marshal(&prevLatest)
		if err != nil {
			return err
		}
		if err := batch.Put([]byte(s.versionedStorageKey(bucket, key, prevLatest.VersionId)), prevBytes); err != nil {
			return fmt.Errorf("failed to update previous version: %w", err)
		}
	}

	return batch.Commit()
}

// findDeleteMarkerRecord returns the delete-marker record a failed
// resolution should report alongside not-found, mirroring the layout
// rules of resolveObjectMetaPB, or nil when the failure is a plain miss.
func (s *ObjectStore) findDeleteMarkerRecord(bucket, key, versionId string) *pb.Object {
	tryKeys := []string{}
	if versionId != "" {
		tryKeys = append(tryKeys, s.versionedStorageKey(bucket, key, versionId))
	} else {
		tryKeys = append(tryKeys, s.latestKeyStorageKey(bucket, key), s.versionedStorageKey(bucket, key, "null"))
	}
	for _, storageKey := range tryKeys {
		var marker pb.Object
		if err := s.BaseStore.GetProto(storageKey, &marker); err == nil && marker.IsDeleteMarker {
			return &marker
		}
	}
	return nil
}

// readRecordBlob opens the blob matching an object record's version
// layout. A real version ID addresses the versioned blob. The null
// version may live at either location — never-enabled and suspended
// writes use the plain one, encryption updates the null-versioned one —
// and the blob layer's null lookup falls back across both tiers.
func (s *ObjectStore) readRecordBlob(ctx context.Context, bucket, key, versionID string) (io.ReadCloser, *storage.BlobMetadata, error) {
	blobVersion := versionID
	if blobVersion == "" {
		blobVersion = "null"
	}
	return s.blobStore.GetWithVersion(ctx, bucket, key, blobVersion)
}

// readRecordBlobRange is the range-read twin of readRecordBlob; the blob
// layer's range lookup has no null fallback, so the plain location is
// tried explicitly.
func (s *ObjectStore) readRecordBlobRange(ctx context.Context, bucket, key, versionID string, offset, length int64) (io.ReadCloser, *storage.BlobMetadata, error) {
	blobVersion := versionID
	if blobVersion == "" {
		blobVersion = "null"
	}
	reader, meta, err := s.blobStore.GetRangeWithVersion(ctx, bucket, key, blobVersion, offset, length)
	if err == nil || blobVersion != "null" {
		return reader, meta, err
	}
	return s.blobStore.GetRange(ctx, bucket, key, offset, length)
}

// blobReadErr maps a blob-tier read error to its caller-facing form: a
// blob that addresses no stored object is the not-found sentinel, while
// any other failure (I/O, permissions) propagates so a store fault is
// never rendered as a missing object.
func blobReadErr(err error) error {
	if errors.Is(err, storage.ErrBlobNotFound) {
		return ErrObjectNotFound
	}
	return err
}

// GetWithVersion retrieves a specific version of an object.
func (s *ObjectStore) GetWithVersion(ctx context.Context, bucket, key, versionId string) (io.ReadCloser, *Object, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		// A delete marker reports the marker record alongside not-found so
		// callers can surface the marker's metadata.
		if marker := s.findDeleteMarkerRecord(bucket, key, versionId); marker != nil {
			return nil, ProtoToObject(marker), ErrObjectNotFound
		}
		return nil, nil, err
	}

	obj := ProtoToObject(pbObj)

	// The blob tier mirrors the record's layout, not the bucket's current
	// status (see readRecordBlob).
	reader, blobMeta, err := s.readRecordBlob(ctx, bucket, key, obj.VersionID)
	if err != nil {
		return nil, nil, blobReadErr(err)
	}

	applyBlobFacts(obj, blobMeta)

	return reader, obj, nil
}

// PutWithVersioning stores an object with versioning support.
func (s *ObjectStore) PutWithVersioning(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata) (*Object, error) {
	return s.putObjectVersioned(ctx, bucket, key, reader, contentType, metadata, isDeleteMarker, storageClass, sysMeta, nil)
}

// putObjectVersioned is the single put path for plain and encrypted
// writes: it assigns the version identity from the bucket's layout,
// writes the blob, builds the record through newObject, and persists it
// through the shared versioned/unversioned tail. A non-nil sseMetadata
// marks an encrypted write — the record carries the SSE fields and its
// size prefers the unencrypted size the encryption recorded.
func (s *ObjectStore) putObjectVersioned(ctx context.Context, bucket, key string, reader io.Reader, contentType string, metadata map[string]string, isDeleteMarker bool, storageClass ObjectStorageClass, sysMeta *SystemMetadata, sseMetadata *SSEObjectMetadata) (*Object, error) {
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
	if !isDeleteMarker && reader != nil {
		var err error
		if isVersioned {
			blobMetaResult, err = s.blobStore.PutWithVersion(ctx, bucket, key, versionId, reader, blobMeta)
		} else {
			blobMetaResult, err = s.blobStore.Put(ctx, bucket, key, reader, blobMeta)
		}
		if err != nil {
			return nil, err
		}
	}

	obj := newObject(key, bucket, contentType, metadata, versionId, isDeleteMarker, storageClass, sysMeta)
	if blobMetaResult != nil {
		applyBlobFacts(obj, blobMetaResult)
	}
	if sseMetadata != nil {
		obj.SSEMetadata = sseMetadata
		obj.ServerSideEncryption = string(sseMetadata.EncryptionType)
		if sseMetadata.KMSKeyID != "" {
			obj.SSEKMSKeyID = sseMetadata.KMSKeyID
		}
		if sseMetadata.UnencryptedSize > 0 {
			obj.Size = sseMetadata.UnencryptedSize
		}
	}

	if isVersioned {
		if err := s.putVersionedObject(bucket, key, versionId, obj); err != nil {
			return nil, err
		}
		return obj, nil
	}
	if err := s.persistUnversionedObject(ctx, bucket, key, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// persistUnversionedObject is the unversioned persistence tail shared by
// every put flavour: the null record becomes the current version, a stale
// null-versioned blob from an earlier encryption update is removed so it
// cannot shadow the plain blob, and any "_latest" pointer left by an
// earlier enabled period is demoted — the previous latest is flagged
// non-current and the pointer deleted, or the stale pointer would shadow
// the unversioned write.
func (s *ObjectStore) persistUnversionedObject(ctx context.Context, bucket, key string, obj *Object) error {
	lockKey := bucket + keySep + key
	s.keyLocker.Lock(lockKey)
	defer s.keyLocker.Unlock(lockKey)

	storageKey := s.versionedStorageKey(bucket, key, "null")
	// The outgoing record decides the stale-blob cleanup: an encryption
	// update of the null version is the only writer that leaves its blob
	// at the null-versioned location, and it also stamps SSE metadata on
	// the record. Without that marker the blob already lives at the plain
	// location and the versioned delete would be a wasted miss.
	var prevNull pb.Object
	hadNullSSEBlob := s.BaseStore.GetProto(storageKey, &prevNull) == nil && prevNull.SseMetadata != nil

	// The plain blob write is already current (the caller persists the
	// record only after the blob write succeeded); a stale null-versioned
	// blob from an earlier encryption update would shadow it in the null
	// lookup order. It is removed BEFORE the records commit: a failure
	// then leaves only a removed stale blob, while committing the record
	// first could leave the new record shadowed by the stale blob.
	if hadNullSSEBlob {
		if err := s.blobStore.DeleteWithVersion(ctx, bucket, key, "null"); err != nil {
			return err
		}
	}

	// The null record, the previous latest's demotion and the "_latest"
	// pointer removal commit in one atomic batch — no partial state can
	// leave a stale pointer shadowing the suspended write or the old
	// IsLatest record surfacing beside it.
	latestKey := s.latestKeyStorageKey(bucket, key)
	var prevLatest pb.Object
	hasPrevLatest := s.BaseStore.GetProto(latestKey, &prevLatest) == nil

	batchBucket, ok := s.BaseStore.Bucket().(storage.BatchBucket)
	if !ok {
		return fmt.Errorf("s3: storage bucket does not support atomic batches")
	}
	batch := batchBucket.NewBatch()
	defer batch.Close()

	objBytes, err := proto.Marshal(ObjectToProto(obj))
	if err != nil {
		return err
	}
	if err := batch.Put([]byte(storageKey), objBytes); err != nil {
		return err
	}
	if hasPrevLatest {
		if prevLatest.VersionId != "" && prevLatest.VersionId != "null" {
			prevLatest.IsLatest = false
			prevBytes, err := proto.Marshal(&prevLatest)
			if err != nil {
				return err
			}
			if err := batch.Put([]byte(s.versionedStorageKey(bucket, key, prevLatest.VersionId)), prevBytes); err != nil {
				return fmt.Errorf("failed to update previous version: %w", err)
			}
		}
		if err := batch.Delete([]byte(latestKey)); err != nil {
			return err
		}
	}
	return batch.Commit()
}

// versioningStatus reads the bucket's versioning status directly from the
// bucket record. The enabled-bool cache collapses Suspended and never-enabled
// into one value; the delete path must tell them apart. A read failure is
// returned to the caller — degrading it to "" would route the delete into
// the destructive marker-less branch on what may be an enabled bucket whose
// record merely failed to read.
func (s *ObjectStore) versioningStatus(bucket string) (BucketVersioningStatus, error) {
	b, err := s.bucketStore.Get(bucket)
	if err != nil {
		return "", err
	}
	return b.VersioningStatus, nil
}

// DeleteWithVersion deletes a specific version of an object.
// If the version ID is provided, it deletes that specific version.
// If the object is versioned and the deleted version is the latest, it updates the latest pointer.
func (s *ObjectStore) DeleteWithVersion(ctx context.Context, bucket, key, versionId string) (*Object, error) {
	if versionId != "" {
		// A suspended bucket still deletes its pre-suspension versions by
		// id, so the gate is the record's existence, not the bucket's
		// current status. The null version passes the same existence gate:
		// non-enabled buckets keep their records at the null-versioned key
		// (the unversioned persistence tail writes there), and "null" is a
		// legal version address — the delete API documents versionId
		// deletes as permanent removals of the addressed version, with the
		// unversioned object itself being the null version.
		canDeleteVersion := s.isVersioningEnabled(bucket) ||
			s.BaseStore.Exists(s.versionedStorageKey(bucket, key, versionId))
		if canDeleteVersion {
			lockKey := bucket + keySep + key
			s.keyLocker.Lock(lockKey)
			defer s.keyLocker.Unlock(lockKey)

			latestKey := s.latestKeyStorageKey(bucket, key)
			var latestObj pb.Object
			isLatest := false
			if err := s.BaseStore.GetProto(latestKey, &latestObj); err == nil {
				isLatest = (latestObj.VersionId == versionId)
			}

			// The metadata record goes first: a delete failure then leaves
			// only an orphaned blob, while the reverse order could leave a
			// record that lists but cannot be served.
			storageKey := s.versionedStorageKey(bucket, key, versionId)
			// The blob mirrors the record's layout. A real version always
			// lives at its versioned slot; the null version keeps its blob
			// at the plain slot on every non-enabled write — only an
			// encryption update of the null version moves it to the
			// null-versioned slot (the same marker the unversioned
			// persistence tail keys its stale-blob cleanup on).
			var outgoing pb.Object
			recordFound := s.BaseStore.GetProto(storageKey, &outgoing) == nil
			hadSSEBlob := recordFound && outgoing.SseMetadata != nil
			if err := s.BaseStore.Delete(storageKey); err != nil {
				return nil, err
			}
			if versionId == "null" && !hadSSEBlob {
				if err := s.blobStore.Delete(ctx, bucket, key); err != nil {
					return nil, err
				}
			} else if err := s.blobStore.DeleteWithVersion(ctx, bucket, key, versionId); err != nil {
				return nil, err
			}

			if isLatest {
				prefix := bucket + keySep + key + keySep
				var remainingVersions []*Object
				err := s.ScanPrefix(prefix, func(k string, v []byte) error {
					if strings.HasSuffix(k, keySep+"_latest") {
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
					// The new latest is the most recently modified
					// remaining version. Version IDs are UUIDs with no
					// intrinsic order, so identical timestamps (the same
					// write burst) break ties by version ID — an
					// arbitrary but total order, so every re-latesting of
					// the same state picks the same version.
					var newLatest *Object
					for _, v := range remainingVersions {
						if newLatest == nil ||
							v.LastModified.After(newLatest.LastModified) ||
							(v.LastModified.Equal(newLatest.LastModified) && v.VersionID > newLatest.VersionID) {
							newLatest = v
						}
					}
					// Every IsLatest flip and the pointer rewrite commit
					// in one atomic batch: a partial failure can never
					// leave two records flagged latest, or none.
					batchBucket, ok := s.BaseStore.Bucket().(storage.BatchBucket)
					if !ok {
						return nil, fmt.Errorf("s3: storage bucket does not support atomic batches")
					}
					batch := batchBucket.NewBatch()
					defer batch.Close()
					for _, v := range remainingVersions {
						v.IsLatest = v == newLatest
						vBytes, err := proto.Marshal(ObjectToProto(v))
						if err != nil {
							return nil, err
						}
						if err := batch.Put([]byte(s.versionedStorageKey(bucket, key, v.VersionID)), vBytes); err != nil {
							return nil, err
						}
					}
					newLatestBytes, err := proto.Marshal(ObjectToProto(newLatest))
					if err != nil {
						return nil, err
					}
					if err := batch.Put([]byte(latestKey), newLatestBytes); err != nil {
						return nil, err
					}
					if err := batch.Commit(); err != nil {
						return nil, err
					}
				}
			}
			// The removed record is returned so callers report whether the
			// deleted version was a delete marker (DeleteObjects and
			// DeleteObject surface DeleteMarker/VersionId for marker
			// removals); a missing record deletes nothing and reports nil.
			if recordFound {
				return ProtoToObject(&outgoing), nil
			}
			return nil, nil
		} else {
			// VersionId specified but versioning not enabled: a bucket that
			// never had versioning only ever holds the null version, so any
			// other version id references a version that does not exist.
			return nil, ErrVersioningNotEnabled
		}
	}

	status, err := s.versioningStatus(bucket)
	if err != nil {
		return nil, err
	}
	if status == BucketVersioningEnabled {
		deleteMarker, err := s.PutWithVersioning(ctx, bucket, key, nil, "", nil, true, StorageClassStandard, nil)
		if err != nil {
			return nil, err
		}
		return deleteMarker, nil
	}
	if status == BucketVersioningSuspended {
		// A suspended bucket "removes the object that has a null versionId,
		// if there is one, and inserts a delete marker that becomes the
		// current version of the object" (AWS DeleteObject). The marker
		// write lands as a null record and demotes any previous latest.
		// The record goes first, as in every delete branch: a failure
		// after the marker write leaves only the old null data
		// unreachable behind the marker, while the reverse order could
		// leave a null record whose data is already gone.
		nullKey := s.versionedStorageKey(bucket, key, "null")
		var nullObj pb.Object
		hadNullObject := s.BaseStore.GetProto(nullKey, &nullObj) == nil && !nullObj.IsDeleteMarker
		deleteMarker, err := s.PutWithVersioning(ctx, bucket, key, nil, "", nil, true, StorageClassStandard, nil)
		if err != nil {
			return nil, err
		}
		if hadNullObject {
			// The null object's blob may live at either location (see
			// readRecordBlob); remove both.
			if err := s.blobStore.Delete(ctx, bucket, key); err != nil {
				return nil, err
			}
			if err := s.blobStore.DeleteWithVersion(ctx, bucket, key, "null"); err != nil {
				return nil, err
			}
		}
		return deleteMarker, nil
	}

	// The record goes first, as in every delete branch: a failure after
	// the record delete leaves only an orphaned blob, while the reverse
	// order could leave a record that lists but cannot be served.
	storageKey := s.versionedStorageKey(bucket, key, "null")
	if err := s.BaseStore.Delete(storageKey); err != nil {
		return nil, err
	}
	if err := s.blobStore.Delete(ctx, bucket, key); err != nil {
		return nil, err
	}
	return nil, nil
}

// HeadWithVersion retrieves metadata for a specific version of an object.
func (s *ObjectStore) HeadWithVersion(ctx context.Context, bucket, key, versionId string) (*Object, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		// A delete marker reports the marker record alongside not-found so
		// callers can surface the marker's metadata, matching the Get and
		// range read paths.
		if marker := s.findDeleteMarkerRecord(bucket, key, versionId); marker != nil {
			return ProtoToObject(marker), ErrObjectNotFound
		}
		return nil, err
	}

	obj := ProtoToObject(pbObj)

	// The blob tier mirrors the record's layout, not the bucket's current
	// status (see readRecordBlob); the full-read lookup covers both null
	// locations, so the metadata comes from the same resolution.
	reader, blobMeta, err := s.readRecordBlob(ctx, bucket, key, obj.VersionID)
	if err != nil {
		return nil, blobReadErr(err)
	}
	reader.Close()

	applyBlobFacts(obj, blobMeta)

	return obj, nil
}

// getVersionedObjectMeta reads object metadata from Pebble only, without
// touching the blob store. Suitable for callers that need the Object struct
// (e.g. lock/retention metadata) but not the blob data or blob-derived
// metadata (size, ETag from blob).
func (s *ObjectStore) getVersionedObjectMeta(bucket, key, versionId string) (*Object, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		return nil, err
	}
	return ProtoToObject(pbObj), nil
}

// GetRangeWithVersion retrieves a range of bytes from a specific version of an object.
func (s *ObjectStore) GetRangeWithVersion(ctx context.Context, bucket, key, versionId string, offset, length int64) (io.ReadCloser, *Object, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		if marker := s.findDeleteMarkerRecord(bucket, key, versionId); marker != nil {
			return nil, ProtoToObject(marker), ErrObjectNotFound
		}
		return nil, nil, err
	}

	obj := ProtoToObject(pbObj)

	reader, blobMeta, err := s.readRecordBlobRange(ctx, bucket, key, obj.VersionID, offset, length)
	if err != nil {
		return nil, nil, blobReadErr(err)
	}

	applyBlobFacts(obj, blobMeta)

	return reader, obj, nil
}

// SetACLWithVersion sets the access control list for a specific version of
// an object. When no explicit versionId is given the current version is
// updated (the "_latest" pointer with null-record fallback). The versioned
// record and the _latest pointer are kept in sync through the shared
// mutation skeleton.
func (s *ObjectStore) SetACLWithVersion(bucket, key, versionId string, acp *AccessControlPolicy) error {
	return s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		obj.ACL = acp
		return nil
	})
}

// GetACLWithVersion retrieves the access control list for a specific
// version of an object, through the shared layout resolution.
func (s *ObjectStore) GetACLWithVersion(bucket, key, versionId string) (*AccessControlPolicy, error) {
	pbObj, err := s.resolveObjectMetaPB(bucket, key, versionId)
	if err != nil {
		return nil, err
	}
	return ProtoToObject(pbObj).ACL, nil
}
