package s3

import (
	"context"
	"fmt"
)

// SetObjectLegalHold sets the legal hold status for an object version.
func (s *ObjectStore) SetObjectLegalHold(ctx context.Context, bucket, key, versionId string, legalHold *ObjectLockLegalHold) error {
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

	obj.ObjectLockLegalHold = legalHold

	isVersioned := s.isVersioningEnabled(bucket)

	var storageKey string
	if isVersioned {
		vid := versionId
		if vid == "" {
			vid = obj.VersionID
		}
		storageKey = s.versionedStorageKey(bucket, key, vid)
	} else {
		storageKey = s.storageKey(bucket, key)
	}

	if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
		return err
	}

	if isVersioned {
		latestKey := s.latestKeyStorageKey(bucket, key)
		if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
			return err
		}
	}

	return nil
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

	obj.ObjectLockRetention = retention

	isVersioned := s.isVersioningEnabled(bucket)

	var storageKey string
	if isVersioned {
		vid := versionId
		if vid == "" {
			vid = obj.VersionID
		}
		storageKey = s.versionedStorageKey(bucket, key, vid)
	} else {
		storageKey = s.storageKey(bucket, key)
	}

	if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
		return err
	}

	if isVersioned {
		latestKey := s.latestKeyStorageKey(bucket, key)
		if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
			return err
		}
	}

	return nil
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
