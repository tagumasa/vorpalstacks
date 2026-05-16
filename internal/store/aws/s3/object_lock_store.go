package s3

import (
	"context"
	"fmt"
)

// updateObjectLockMetadata is a shared helper for SetObjectLegalHold and
// SetObjectRetention. It locks the key, reads the target version, applies
// the mutation via the callback, writes back the versioned record, and only
// updates the _latest pointer when the modified version is the current latest.
func (s *ObjectStore) updateObjectLockMetadata(ctx context.Context, bucket, key, versionId string, mutate func(obj *Object)) error {
	lockKey := bucket + "#" + key
	s.keyLocker.Lock(lockKey)
	defer func() {
		s.keyLocker.Unlock(lockKey)
		s.keyLocker.Delete(lockKey)
	}()

	_, obj, err := s.GetWithVersion(ctx, bucket, key, versionId)
	if err != nil {
		return err
	}

	mutate(obj)

	isVersioned := s.isVersioningEnabled(bucket)

	var storageKey string
	if isVersioned {
		vid := versionId
		if vid == "" {
			vid = obj.VersionID
		}
		storageKey = s.versionedStorageKey(bucket, key, vid)
	} else {
		storageKey = s.versionedStorageKey(bucket, key, "null")
	}

	if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
		return err
	}

	if isVersioned && obj.IsLatest {
		latestKey := s.latestKeyStorageKey(bucket, key)
		if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
			return err
		}
	}

	return nil
}

// SetObjectLegalHold sets the legal hold status for an object version.
func (s *ObjectStore) SetObjectLegalHold(ctx context.Context, bucket, key, versionId string, legalHold *ObjectLockLegalHold) error {
	return s.updateObjectLockMetadata(ctx, bucket, key, versionId, func(obj *Object) {
		obj.ObjectLockLegalHold = legalHold
	})
}

// GetObjectLegalHold retrieves the legal hold status for an object version.
func (s *ObjectStore) GetObjectLegalHold(ctx context.Context, bucket, key, versionId string) (*ObjectLockLegalHold, error) {
	_, obj, err := s.GetWithVersion(ctx, bucket, key, versionId)
	if err != nil {
		return nil, err
	}

	if obj.ObjectLockLegalHold == nil {
		return &ObjectLockLegalHold{Status: ObjectLockLegalHoldOff}, nil
	}

	return obj.ObjectLockLegalHold, nil
}

// SetObjectRetention sets the retention policy for an object version.
func (s *ObjectStore) SetObjectRetention(ctx context.Context, bucket, key, versionId string, retention *ObjectLockRetention) error {
	return s.updateObjectLockMetadata(ctx, bucket, key, versionId, func(obj *Object) {
		obj.ObjectLockRetention = retention
	})
}

// GetObjectRetention retrieves the retention policy for an object version.
func (s *ObjectStore) GetObjectRetention(ctx context.Context, bucket, key, versionId string) (*ObjectLockRetention, error) {
	_, obj, err := s.GetWithVersion(ctx, bucket, key, versionId)
	if err != nil {
		return nil, err
	}

	if obj.ObjectLockRetention == nil {
		return nil, fmt.Errorf("retention configuration not found")
	}

	return obj.ObjectLockRetention, nil
}
